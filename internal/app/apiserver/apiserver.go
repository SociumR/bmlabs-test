package apiserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SociumR/bmlabs-test/store"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type apiserver struct {
	router *mux.Router
	store  *store.Store
}

// APIServer interface
type APIServer interface {
	Start()
}

// New ...
func new() *apiserver {
	mDB := MongoDB("mongodb://127.0.0.1:27018", "bmlabs")

	router := mux.NewRouter()

	s := &apiserver{
		router: router,
		store:  store.New(mDB),
	}

	s.configureRouter()

	return s
}

// Start ...
func Start() error {
	s := new()

	return http.ListenAndServe(":9008", s)
}

func (s *apiserver) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func (s *apiserver) configureRouter() {

	s.router.HandleFunc("/users", s.CreateUser).Methods(http.MethodPost)
	s.router.HandleFunc("/users", s.GetUsers).Methods(http.MethodGet)

	s.router.HandleFunc("/games", s.GetGames).Methods(http.MethodGet)
	s.router.HandleFunc("/games", s.CreateGame).Methods(http.MethodPost)

	s.router.HandleFunc("/games/aggregate", s.Stat).Methods(http.MethodGet)

	s.router.Use(s.middleware)
}

func (s *apiserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// MongoDB ...
func MongoDB(dns, db string) *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(dns))
	if err != nil {

		log.Fatalf(fmt.Sprintf("Error occured while establishing connection to mongoDB: %s", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	return client.Database(db)
}
