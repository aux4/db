package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
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
	if len(args) < 2 || (args[0] != "execute" && args[0] != "stream") {
		fmt.Fprintf(os.Stderr, "Usage: preprocessor.js <execute|stream> <query> <file> <inputStream> <tx> <ignore> <params>\n")
		os.Exit(1)
	}
	
	action := args[0]
	query := args[1]
	file := ""
	inputStream := false
	tx := false
	ignore := false
	params := NewOrderedMap()
	
	if len(args) > 2 {
		file = args[2]
	}
	if len(args) > 3 {
		inputStream = args[3] == "true"
	}
	if len(args) > 4 {
		tx = args[4] == "true"
	}
	if len(args) > 5 {
		ignore = args[5] == "true"
	}
	if len(args) > 6 {
		json.Unmarshal([]byte(args[6]), params)
	}
	
	// Process additional key-value arguments
	for i := 7; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break // Ensure we have pairs
		}
		key, value := args[i], args[i+1]
		if value == "true" {
			params.Set(key, true)
		} else if value == "false" {
			params.Set(key, false)
		} else if num, err := strconv.ParseFloat(value, 64); err == nil {
			params.Set(key, num)
		} else {
			params.Set(key, value)
		}
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