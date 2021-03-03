package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"gqlexample/v2/graph"
	"gqlexample/v2/graph/generated"
	"gqlexample/v2/internal/auth"
	database "gqlexample/v2/internal/pkg/db/mongodb"
	"gqlexample/v2/internal/pkg/utils"
	"net/http"
)


func main() {
	utils.InitLog()

	router := chi.NewRouter()
	router.Use(auth.Middleware())

	database.InitMongoDBClient("gql_db","graphql","graphql")

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", server)

	logrus.Printf("connect to http://localhost:8080/ for GraphQL playground")
	logrus.Fatal(http.ListenAndServe(":8080", router))
}
