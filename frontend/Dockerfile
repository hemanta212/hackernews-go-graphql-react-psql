FROM node:16.19-alpine as dev
RUN apk add --no-cache make
WORKDIR /node/src/frontend
COPY package.json package-lock.json ./
RUN npm install
EXPOSE 8000

FROM dev as prod
RUN apk add --no-cache make\
    && npm install --location=global serve
EXPOSE 443
CMD ["make", "run"]
