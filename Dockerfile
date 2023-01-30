FROM node:18 as builder

WORKDIR /usr/src/app
COPY package*.json ./

RUN ["npm", "ci", "--only=production"]

FROM gcr.io/distroless/nodejs:18

ENV NODE_ENV = "production"

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/ ./

COPY *.js .

CMD ["server.js"]
