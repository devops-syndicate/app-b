const {format, createLogger, transports} = require("winston");
const {json, timestamp, combine} = format

const logger = createLogger({
  level: "debug",
  format: combine(timestamp(), json()),
  transports: [new transports.Console()],
});

module.exports = logger;