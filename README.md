# FullStack HackerNews Frontend

Trying to clone: https://news.ycombinator.com

## Installation

- Install docker and docker compose

- clone this repository with git

- cd to the cloned folder, create and populate a `.env` file in the topmost/root folder.

```shell
PGUSER=user123
PGPASSWORD=pass123
PGHOST=postgres
PGPORT=5432
PGDATABASE=hackernews

POSTGRES_USER=$PGUSER
POSTGRES_PASSWORD=$PGPASSWORD
POSTGRES_HOST=localhost
POSTGRES_DB=$PGDATABASE

# HTTPS_SSL=/folder/to/the/ssl/certs/and/keys

VITE_API_URL=http://localhost:8008/query
VITE_API_WS_URL=ws://localhost:8008/query
```

- OPTIONAL: `HTTPS_SSL` refers to the folder containing, `fullchain.pem` and `privkey.pem` file, If you have these certs, uncomment the line in above `.env` file, and point `VITE_API_*` vars to https and wss version of your site, otherwise let it be commented, and it will run on port 80 with http.

- Run the command `docker compose up`

- Navigate to: http://localhost or http://localhost:9000 to see the react site, similarly the graphql api playground will be hosted on http://localhost:8080

## Running with ssl mode

- Export the `HTTPS_SSL`, `VITE_API_URL` and `VITE_API_WS_URL` env variable in the .env file and source it in your terminal

```sh
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
