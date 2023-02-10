# FullStack HackerNews Frontend

Trying to clone: https://news.ycombinator.com

## Installation

- Install docker and docker compose

- clone this repository with git

- cd to the cloned folder, create and populate a `.env` file in the topmost/root folder.

```shell
PGUSER=postgres
PGPASSWORD=pass123
PGHOST=postgres
PGPORT=5432
PGDATABASE=hackernews

POSTGRES_USER=$PGUSER
POSTGRES_PASSWORD=$PGPASSWORD
POSTGRES_HOST=localhost
POSTGRES_DB=$PGDATABASE

# HTTPS_SSL_PATH=/path/to/your/ssl/cert/folder
# SSL_CERT_FILE=fullchain.pem
# SSL_KEY_FILE=privkey.pem

VITE_API_URL=http://localhost:8008/query
VITE_API_WS_URL=ws://localhost:8008/query
```

- OPTIONAL: `HTTPS_SSL_PATH`, refers to the absolute path of folder containing ssl certs and `SSL_CERT_FILE` and `SSL_KEY_FILE` refers to the name of certificate and key file.
 If you have these certs, uncomment above lines in `.env` file, and point `VITE_API_*` vars to https and wss version of your site.

- Run the command `docker compose up` to spin up dev environment, optionally create a dev volume with `docker volume create hackernews-postgres-dev`

- Navigate to: http://localhost:8000 to see the frontend, similarly the graphql api playground will be hosted on http://localhost:8080

- To get a Prod build (without ssl), `export HTTPS_SSL_PATH=/non/existent/path` and run `docker compose -f docker-compose.yml -f docker-compose.prod.yml up`
  This will spin up the production go server at port 8008 and frontend at both ports 443 and 9000.

## Running with SSL mode

- Export the `HTTPS_SSL_PATH`, `SSL_CERT_FILE`, `SSL_KEY_FILE`, `VITE_API_URL` and `VITE_API_WS_URL` env variables in the .env file and source it in your terminal

```sh
docker volume create hackernews-postgres-dev
source .env
docker compose -f docker-compose.yml -f docker-compose.prod.yml up
```

## Frontend
- Uses React.js + Vite with Apollo client library

## Backend
- Uses go, graphql and psql
- Here's the schema

```graphql
type Feed {
     id: ID!
     links: [Link!]!
     count: Int!
}

type Link {
    id: ID!
    description: String!
    postedBy: User!
    url: String!
    createdAt: Time!
    votes: [Vote!]!
}

type User {
     id: ID!
     username: String!
     email: String
}

type Vote {
     id: ID!
     link: Link!
     user: User!
}

type AuthPayload {
     token: String
     user: User
}

type Query {
     feed(filter: String, offset: Int, limit: Int, orderBy: LinkOrderByInput): Feed!
}

type Mutation {
     post(input: NewLink!): Link!
     signup(input: NewUser!): AuthPayload
     login(input: Login!): AuthPayload
     vote(linkID: ID!): Vote
     refreshToken(input: RefreshTokenInput!): String!
}

type Subscription{
     newLink: Link!
     newVote: Vote!
}

input RefreshTokenInput {
      token: String!
}

input NewUser {
      username: String!
      password: String!
      email: String
}

input NewLink {
      description: String!
      url: String!
}

input Login {
      username: String!
      password: String!
}

input LinkOrderByInput {
  description: Sort
  createdAt: Sort
}

enum Sort {
  asc
  desc
}

scalar Time
```
