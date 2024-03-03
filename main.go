package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nikitahahua/RssAgg/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	log.SetPrefix("georgian greetings: ")
	log.SetFlags(0)

	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	connection, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("cant connect to the data base by url : \n" + dbUrl)
	}

	dbQueries := database.New(connection)

	apiCfg := apiConfig{
		DB: dbQueries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get(
		"/healthz",
		handleRequest,
	)
	v1Router.Get(
		"/err",
		handlerErr,
	)
	v1Router.Get(
		"/users",
		apiCfg.middlewareAuth(apiCfg.handleGetUser),
	)
	v1Router.Post(
		"/users",
		apiCfg.handlerCreateUser,
	)

	v1Router.Post(
		"/feeds",
		apiCfg.middlewareAuth(apiCfg.handlerCreateFeed),
	)
	v1Router.Get(
		"/feeds",
		apiCfg.handlerGetFeeds,
	)
	v1Router.Post(
		"/feed_follows",
		apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows),
	)
	v1Router.Get(
		"/feed_follows",
		apiCfg.middlewareAuth(apiCfg.handlerFeedFollowsGet),
	)
	v1Router.Delete(
		"/feed_follows/{feed_follow_id}",
		apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow),
	)

	v1Router.Get(
		"/user_posts",
		apiCfg.middlewareAuth(apiCfg.handleGetPostsForUser),
	)

	router.Mount("/v1", v1Router)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	const collectionConcurrency = 10
	const collectionInterval = time.Minute
	go startScraping(dbQueries, collectionConcurrency, collectionInterval)

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
