import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { setContext } from "@apollo/client/link/context";
import { split } from "@apollo/client";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { createClient } from "graphql-ws";
import { getMainDefinition } from "@apollo/client/utilities";

import {
  ApolloProvider,
  ApolloClient,
  createHttpLink,
  InMemoryCache,
} from "@apollo/client";

import { AUTH_TOKEN } from "./constants";
import "./styles/index.css";
import App from "./components/App";

const httpLink = createHttpLink({
  uri:
    import.meta.env.MODE == "production"
      ? import.meta.env.VITE_API_URL
      : import.meta.env["VITE_DEV_API_URL"] || "http://localhost:8080/query",
});

const authLink = setContext((_, { headers }) => {
  const token = localStorage.getItem(AUTH_TOKEN);
  return {
    headers: {
      ...headers,
      authorization: token ? `${token}` : "",
    },
  };
});

const wsLink = new GraphQLWsLink(
  createClient({
    url:
      import.meta.env.MODE == "production"
        ? import.meta.env.VITE_API_WS_URL
        : import.meta.env.VITE_DEV_API_WSL_URL || "ws://localhost:8080/query",
    options: {
      connectionParams: {
        authToken: localStorage.getItem(AUTH_TOKEN),
      },
    },
  })
);

const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === "OperationDefinition" &&
      definition.operation === "subscription"
    );
  },
  wsLink,
  authLink.concat(httpLink)
);

const client = new ApolloClient({
  link: splitLink,
  cache: new InMemoryCache(),
});

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <ApolloProvider client={client}>
        <App />
      </ApolloProvider>
    </BrowserRouter>
  </React.StrictMode>
);
