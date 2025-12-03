package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// OrderedMap preserves key insertion order
type OrderedMap struct {
	keys   []string
	values map[string]interface{}
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		keys:   []string{},
		values: make(map[string]interface{}),
	}
}

func (om *OrderedMap) Set(key string, value interface{}) {
	if _, exists := om.values[key]; !exists {
		om.keys = append(om.keys, key)
	}
	om.values[key] = value
}

func (om *OrderedMap) Get(key string) (interface{}, bool) {
	value, exists := om.values[key]
	return value, exists
}

func (om *OrderedMap) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	
	for i, key := range om.keys {
		if i > 0 {
			buf.WriteByte(',')
		}
		
		// Marshal key
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')
		
		// Marshal value
		valueBytes, err := json.Marshal(om.values[key])
		if err != nil {
			return nil, err
		}
		buf.Write(valueBytes)
	}
	
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func (om *OrderedMap) UnmarshalJSON(data []byte) error {
	// Parse JSON maintaining field order
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	
	// Parse opening brace
	token, err := decoder.Token()
	if err != nil {
		return err
	}
	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected opening brace")
	}
	
	for decoder.More() {
		// Parse key
		keyToken, err := decoder.Token()
		if err != nil {
			return err
		}
		key := keyToken.(string)
		
		// Parse value
		var value interface{}
		if err := decoder.Decode(&value); err != nil {
			return err
		}
		
		om.Set(key, value)
	}
	
	// Parse closing brace
	_, err = decoder.Token()
	return err
}

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: aux4-db <command> <args...>\n")
		fmt.Fprintf(os.Stderr, "Commands: prepare|sql-bind\n")
		os.Exit(1)
	}

	command := args[0]

	// Handle sql-bind command
	if command == "sql-bind" {
		handleSqlBinder(args)
		return
	}

	// Handle prepare command
	if command == "prepare" {
		handlePrepare(args)
		return
	}

	fmt.Fprintf(os.Stderr, "Invalid command: %s. Supported commands: prepare, sql-bind\n", command)
	os.Exit(1)
}

func handlePrepare(args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: aux4-db prepare <action> <query> [file] [inputStream] [tx] [ignore] [params] [key value...]\n")
		os.Exit(1)
	}

	action := args[1]

	// Validate action
	if action != "execute" && action != "stream" {
		fmt.Fprintf(os.Stderr, "Invalid action: %s. Supported actions: execute, stream\n", action)
		os.Exit(1)
	}

	query := args[2]
	file := ""
	inputStream := false
	tx := false
	ignore := false
	params := NewOrderedMap()

	if len(args) > 3 {
		file = args[3]
	}
	if len(args) > 4 {
		inputStream = args[4] == "true"
	}
	if len(args) > 5 {
		tx = args[5] == "true"
	}
	if len(args) > 6 {
		ignore = args[6] == "true"
	}
	if len(args) > 7 {
		json.Unmarshal([]byte(args[7]), params)
	}

	// Process additional key-value arguments
	for i := 8; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break // Ensure we have pairs
		}
		key, value := args[i], args[i+1]
		// Parse value based on context and action
		parsedValue := parseParameterValue(value, action)
		params.Set(key, parsedValue)
	}

	// Get SQL
	sql := query
	if sql == "" && file != "" {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading SQL file: %v\n", err)
			os.Exit(1)
		}
		sql = string(content)
	}
	if sql == "" {
		fmt.Fprintf(os.Stderr, "Either query or file parameter must be provided\n")
		os.Exit(1)
	}

	if inputStream {
		processInputStream(action, sql, params, tx, ignore)
	} else {
		output(action, sql, params, ignore)
	}
}

func processInputStream(action, sql string, baseParams *OrderedMap, tx bool, ignore bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var items []*OrderedMap
	var input strings.Builder
	
	for scanner.Scan() {
		input.WriteString(scanner.Text())
		input.WriteByte('\n')
	}
	
	inputStr := input.String()
	if len(inputStr) == 0 {
		return
	}
	
	// Parse JSON input maintaining order
	if inputStr[0] == '[' {
		// Parse JSON array maintaining field order
		decoder := json.NewDecoder(strings.NewReader(inputStr))
		
		// Parse opening bracket
		token, err := decoder.Token()
		if err == nil && token == json.Delim('[') {
			for decoder.More() {
				item := NewOrderedMap()
				if err := decoder.Decode(item); err == nil {
					items = append(items, item)
				}
			}
		}
	} else if inputStr[0] == '{' {
		// Check if it's line-delimited JSON (contains newlines between objects)
		if strings.Contains(inputStr, "}\n{") || strings.Contains(inputStr, "}\r\n{") {
			// Line-delimited JSON
			scanner := bufio.NewScanner(strings.NewReader(inputStr))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line != "" {
					item := NewOrderedMap()
					if json.Unmarshal([]byte(line), item) == nil {
						items = append(items, item)
					}
				}
			}
		} else {
			// Single JSON object
			item := NewOrderedMap()
			if json.Unmarshal([]byte(inputStr), item) == nil {
				items = append(items, item)
			}
		}
	} else {
		// Other line-delimited formats
		scanner := bufio.NewScanner(strings.NewReader(inputStr))
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				item := NewOrderedMap()
				if json.Unmarshal([]byte(line), item) == nil {
					items = append(items, item)
				}
			}
		}
	}
	
	// Merge base params into each item (params override input stream data)
	mergedItems := make([]*OrderedMap, len(items))
	for i, item := range items {
		mergedItem := NewOrderedMap()
		// Add item params first to maintain their order
		for _, key := range item.keys {
			if value, exists := item.Get(key); exists {
				mergedItem.Set(key, value)
			}
		}
		// Add base params, which override item params
		for _, key := range baseParams.keys {
			if value, exists := baseParams.Get(key); exists {
				mergedItem.Set(key, value)
			}
		}
		mergedItems[i] = mergedItem
	}

	// When ignore=true, always process items individually to allow partial success
	if ignore && !tx {
		for _, mergedItem := range mergedItems {
			output(action, sql, mergedItem, ignore)
		}
	} else if tx || (action == "execute" && len(items) > 1) {
		// For streaming with transaction, we still use batch but the db-sqlite will stream results
		if action == "stream" && tx {
			outputStreamBatch(sql, baseParams, mergedItems, tx, ignore)
		} else {
			outputBatch(sql, baseParams, mergedItems, tx || (action == "execute" && len(items) > 1), ignore)
		}
	} else {
		for _, mergedItem := range mergedItems {
			output(action, sql, mergedItem, ignore)
		}
	}
}

// filterAux4Params removes aux4-specific parameters, keeping only database parameters
func filterAux4Params(params *OrderedMap) *OrderedMap {
	// List of aux4-specific parameter keys to exclude
	excludedKeys := map[string]bool{
		"action":        true,
		"aux4HomeDir":   true,
		"configDir":     true,
		"packageDir":    true,
		"file":          true,
		"inputStream":   true,
		"tx":            true,
		"ignore":        true,
		"query":         true,
	}
	
	filtered := NewOrderedMap()
	for _, key := range params.keys {
		if !excludedKeys[key] {
			if value, exists := params.Get(key); exists {
				filtered.Set(key, value)
			}
		}
	}
	return filtered
}

func output(action, sql string, params *OrderedMap, ignore bool) {
	filtered := filterAux4Params(params)
	result := NewOrderedMap()
	result.Set("action", action)
	result.Set("sql", sql)
	result.Set("params", filtered)
	result.Set("ignore", ignore)
	json.NewEncoder(os.Stdout).Encode(result)
}

func outputBatch(sql string, params *OrderedMap, items []*OrderedMap, useTx bool, ignore bool) {
	result := NewOrderedMap()
	result.Set("action", "executeBatch")
	result.Set("sql", sql)
	result.Set("params", filterAux4Params(params))
	result.Set("items", items)
	result.Set("tx", useTx)
	result.Set("ignore", ignore)
	json.NewEncoder(os.Stdout).Encode(result)
}

func outputStreamBatch(sql string, params *OrderedMap, items []*OrderedMap, useTx bool, ignore bool) {
	result := NewOrderedMap()
	result.Set("action", "streamBatch")
	result.Set("sql", sql)
	result.Set("params", filterAux4Params(params))
	result.Set("items", items)
	result.Set("tx", useTx)
	result.Set("ignore", ignore)
	json.NewEncoder(os.Stdout).Encode(result)
}

func handleSqlBinder(args []string) {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: aux4-db sql-bind <query> <params_json> [key value...]\\n")
		os.Exit(1)
	}

	query := args[1]
	paramsJSON := args[2]

	// Parse base parameters from JSON
	var baseParams map[string]interface{}
	if paramsJSON != "" {
		if err := json.Unmarshal([]byte(paramsJSON), &baseParams); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing params JSON: %v\\n", err)
			os.Exit(1)
		}
	} else {
		baseParams = make(map[string]interface{})
	}

	// Process additional key-value arguments starting from index 3
	for i := 3; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break // Ensure we have pairs
		}
		key, value := args[i], args[i+1]
		// Parse value for sql-bind
		parsedValue := parseParameterValue(value, "sql-bind")
		baseParams[key] = parsedValue
	}

	// Parse and replace parameters in the query
	parsedQuery, err := parseQuery(query, baseParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing query: %v\\n", err)
		os.Exit(1)
	}

	fmt.Print(parsedQuery)
}

func parseParameterValue(value string, action string) interface{} {
	// For sql-bind action, be more aggressive about type conversion
	if action == "sql-bind" {
		// Handle null
		if value == "null" {
			return nil
		}
		// Handle booleans
		if value == "true" {
			return true
		}
		if value == "false" {
			return false
		}
		// Handle JSON arrays
		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
			var arr []interface{}
			if err := json.Unmarshal([]byte(value), &arr); err == nil {
				return arr
			}
		}
		// Handle numbers
		if num, err := strconv.ParseFloat(value, 64); err == nil {
			// Check if it's actually an integer
			if num == float64(int64(num)) {
				return int64(num)
			}
			return num
		}
	} else {
		// For other actions, use simpler logic
		if value == "true" {
			return true
		} else if value == "false" {
			return false
		} else if num, err := strconv.ParseFloat(value, 64); err == nil {
			return num
		}
	}
	return value
}

func parseQuery(query string, params map[string]interface{}) (string, error) {
	// Regular expression to find named parameters like :name, :age, etc.
	paramRegex := regexp.MustCompile(`:(\w+)`)

	// Find all parameter names in the query
	matches := paramRegex.FindAllStringSubmatch(query, -1)
	paramNames := make(map[string]bool)
	for _, match := range matches {
		paramNames[match[1]] = true
	}

	// Check for missing parameters
	for paramName := range paramNames {
		if _, exists := params[paramName]; !exists {
			return "", fmt.Errorf("missing parameter: %s", paramName)
		}
	}

	// Replace each parameter with its safely quoted value
	result := query
	for paramName, value := range params {
		placeholder := ":" + paramName
		if strings.Contains(result, placeholder) {
			quotedValue := quoteValue(value)
			result = strings.ReplaceAll(result, placeholder, quotedValue)
		}
	}

	return result, nil
}

func quoteValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return "NULL"
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case string:
		// Escape single quotes by doubling them (SQL standard)
		escaped := strings.ReplaceAll(v, "'", "''")
		return "'" + escaped + "'"
	case []interface{}:
		// Handle arrays for IN clauses
		if len(v) == 0 {
			return "(NULL)" // Empty array becomes (NULL) for safety
		}
		var items []string
		for _, item := range v {
			items = append(items, quoteValue(item))
		}
		return "(" + strings.Join(items, ", ") + ")"
	default:
		// Check if it's a JSON array that wasn't properly parsed
		if str := fmt.Sprintf("%v", v); strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]") {
			// Try to parse as JSON array
			var arr []interface{}
			if err := json.Unmarshal([]byte(str), &arr); err == nil {
				return quoteValue(arr)
			}
		}
		// For any other type, convert to string and treat as string
		str := fmt.Sprintf("%v", v)
		escaped := strings.ReplaceAll(str, "'", "''")
		return "'" + escaped + "'"
	}
}