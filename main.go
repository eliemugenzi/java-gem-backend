package main

import (
	"java-gem/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// GraphQL Handler
	graphQlHandler := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	r.POST("/graphql", func(c *gin.Context) {
		graphQlHandler.ServeHTTP(c.Writer, c.Request)
	})

	// GraphQL Playground
	playgroundHandler := playground.Handler("GraphQL", "/graphql")
	r.GET("/", func(c *gin.Context) {
		playgroundHandler.ServeHTTP(c.Writer, c.Request)
	})

	r.Run(":8080")

}
