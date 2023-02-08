# FullStack HackerNews Frontend

Running at: https://vps.osac.org.np
Trying to clone: https://news.ycombinator.com

## Installation

- Install docker and docker compose

- clone this repository with git

- cd to the cloned folder, create and populate a .env file.

```shell
PGUSER=user123
PGPASSWORD=pass123
PGHOST=postgres
PGPORT=5432
PGDATABASE=hackernews
PG_DATABASE_URI=postgres://$PGUSER:$PGPASSWORD@$PGHOST/$PGDATABASE

POSTGRES_USER=$PGUSER
POSTGRES_PASSWORD=$PGPASSWORD
POSTGRES_HOST=localhost
POSTGRES_DB=$PGDATABASE

HTTPS_SSL=/folder/to/the/ssl/certs/and/keys
```

- If you don't set the HTTP_SSL then the site will run on http at port 80

- Run the command `docker compose up`

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

enum Sort {
  asc
  desc
}

scalar Time
```
