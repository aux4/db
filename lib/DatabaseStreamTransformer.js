const { Transform } = require("stream");
const fs = require("fs");

class DatabaseStreamTransformer extends Transform {
  constructor(dbConfig, query, params, stream = false) {
    super({
      writableObjectMode: true,
      readableObjectMode: true
    });

    const Database = require(dbConfig.module || `@aux4/db-${dbConfig.type}`);
    this.database = new Database(dbConfig);
    this.query = query;
    this.params = params;
    this.stream = stream;

    this.on("end", async () => {
      await this.database.close();
    });
  }

  async open() {
    await this.database.open();
  }

  _transform(item, encoding, callback) {
    if (this.stream) {
      this.database.stream(this.query, { ...this.params, ...item }).then(stream => {
        stream.on("data", row => console.log(JSON.stringify(row, null, 2)));
        stream.on("error", err => {
          if (params.outputError) {
            fs.appendFileSync(params.outputError, err.message, { encoding: "utf8" });
          } else {
            console.error(err.message.red);
          }
        });
        stream.on("end", () => callback());
      });
      return;
    }

    this.database
      .execute(this.query, { ...this.params, ...item })
      .then(result => {
        this.push(JSON.stringify(result.data, null, 2));
        callback();
      })
      .catch(err => {
        callback(err);
      });
  }
}

module.exports = DatabaseStreamTransformer;
