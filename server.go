package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Prateek61/go_auth/graph"
	"github.com/Prateek61/go_auth/postgres"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-pg/pg/v10"
	"github.com/rs/cors"
	customMiddleware "github.com/Prateek61/go_auth/middleware"
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

	router := chi.NewRouter()

	userRepo := postgres.UsersRepo{DB: DB}

	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
		Debug: true,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(customMiddleware.AuthMiddleware(userRepo))

    
	config := graph.Config{Resolvers: &graph.Resolver{
		TodosRepo: postgres.TodosRepo{DB: DB},
		UsersRepo: userRepo,
	}}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(config))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
