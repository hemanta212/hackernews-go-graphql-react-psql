package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"

	"github.com/hemanta212/hackernews-go-graphql/graph"
	"github.com/hemanta212/hackernews-go-graphql/graph/model"
	"github.com/hemanta212/hackernews-go-graphql/internal/auth"
	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/postgresql"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://vps.osac.org.np",
			"http://vps.osac.org.np",
			"http://vps.osac.org.np:8000",
			"https://vps.hemantasharma.com.np",
			"http://localhost:8000",
			"http://localhost", "http://localhost:80",
			"http://localhost:9000",
			"http://localhost:8080", "http://localhost:8008",
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
		Debug:            true,
	}))
	router.Use(auth.Middleware())
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	database.InitDB()
	defer database.CloseDB()

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		LinkObservers: map[string]chan *model.Link{},
		VoteObservers: map[string]chan *model.Vote{},
	}}))

	srv.AddTransport(transport.POST{})
	// websockets for subcriptions
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	SSLPath := os.Getenv("HTTPS_SSL")
	if SSLPath != "" {
		log.Printf("HTTPS_SSL env var is exported to %s, initiating https protocol", SSLPath)
		log.Printf("connect to https://localhost:%s/ for GraphQL playground", port)
		log.Fatal(http.ListenAndServeTLS("0.0.0.0:"+port, SSLPath+"/fullchain.pem", SSLPath+"/privkey.pem", router))
	} else {
		log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
		log.Fatal(http.ListenAndServe("0.0.0.0:"+port, router))
	}
}
