const http = require('http')
const express = require('express')
const logger = require('./logger')
const app = express()
const promBundle = require("express-prom-bundle");
const { createTerminus } = require('@godaddy/terminus')
const PORT = 8080

const metricsMiddleware = promBundle({
  includeMethod: true, 
  includePath: true, 
  includeStatusCode: true, 
  includeUp: true,
  customLabels: {app: 'app-b'},
  promClient: {collectDefaultMetrics: {}}
});

app.use(metricsMiddleware)

app.get('/', (_, res) => {
  logger.info('Call hello endpoint')
  res.send('Hello "app-b"!')
})

const server = http.createServer(app)

async function onShutdown() {
  logger.info('server shuts down')
}

createTerminus(server, { 
  onShutdown,
  signals: ['SIGINT', 'SIGTERM'],
  useExit0: true 
})

server.listen(PORT, () => {logger.info(`app-b listening on port ${PORT}`)})