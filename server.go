package main

import (
	"log"
	"net/http"
	"server/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const port = "8080"

func InitGraphQLServer() error {

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

	return nil
}

// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/99designs/gqlgen/graphql/handler"
// 	"github.com/99designs/gqlgen/graphql/handler/extension"
// 	"github.com/99designs/gqlgen/graphql/handler/lru"
// 	"github.com/99designs/gqlgen/graphql/handler/transport"
// 	"github.com/99designs/gqlgen/graphql/playground"
// 	"github.com/vektah/gqlparser/v2/ast"

// 	"server/graph"
// )

// const (
// 	graphqlPort = "8080" // Port for GraphQL server
// 	restPort    = "8000" // Port for REST server
// )

// // InitGraphQLServer initializes and starts the GraphQL server.
// func InitGraphQLServer() {
// 	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

// 	// Configure transports and caching
// 	srv.AddTransport(transport.Options{})
// 	srv.AddTransport(transport.GET{})
// 	srv.AddTransport(transport.POST{})
// 	srv.SetQueryCache(lru.Newospection{})
// 	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New })

// 	// Routes for Graphle("/", playground.Handler("GraphQL Playground", "/query"))
// 	http.Handle("/query", srv)

// 	log.Printf("GraphQL server running at http://localhost:%s/", graphqlPort)
// 	log.Fatal(http.ListenAndServe(":"+graphqlPort, nil))
// }

// // InitRESTServer initializes and starts the RESTful server using Gin.
// func InitRESTServer() {
// 	router := gin.Default()

// 	// Define REST endpoints
// 	router.GET("/ping", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "pong"})
// 	})

// 	// Additional REST routes can go here...

// 	log.Printf("RESTful server running at http://localhost:%s/", restPort)
// 	log.Fatal(router.Run(":" + restPort))
// }

// func main() {
// 	// Run GraphQL and REST servers concurrently
// 	go InitGraphQLServer() // Run GraphQL server in a separate goroutine
// 	InitRESTServer()       // Run REST server in the main goroutine
// }
