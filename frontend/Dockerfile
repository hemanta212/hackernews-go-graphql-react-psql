FROM node:16.19-alpine

RUN apk add --no-cache make

RUN npm install --location=global serve

ENV APP_HOME /frontend/app
RUN mkdir -p "$APP_HOME"

WORKDIR "$APP_HOME"

RUN mkdir -p /etc/letsencrypt/live/vps.osac.org.np

COPY package.json .

RUN npm install

EXPOSE 443

CMD make run
