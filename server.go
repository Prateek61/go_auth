package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Prateek61/go_auth/graph"
	"github.com/Prateek61/go_auth/postgres"
	"github.com/go-pg/pg/v10"
)

const defaultPort = "8080"

func main() {
	// Setup DB
	DB := postgres.New(&pg.Options{
		User: "postgres",
		Password: "postgres",
		Database: "go_auth_dev",
	})

	defer DB.Close()

	DB.AddQueryHook(postgres.DBLogger{})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	config := graph.Config{Resolvers: &graph.Resolver{
		TodosRepo: postgres.TodosRepo{DB: DB},
		UsersRepo: postgres.UsersRepo{DB: DB},
	}}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(config))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
