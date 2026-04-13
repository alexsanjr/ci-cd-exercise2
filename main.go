package main

import (
	"api-go/user"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	var repo user.UserRepository
	if os.Getenv("DB_HOST") != "" || os.Getenv("DATABASE_URL") != "" {
		pgRepo, err := user.NewPostgresUserRepositoryFromEnv()
		if err != nil {
			log.Fatal(err)
		}
		repo = pgRepo
	} else {
		repo = user.NewUserRepository()
	}

	service := user.NewUserService(repo)
	controller := user.NewUserController(service)

	mux := http.NewServeMux()

	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListUsers(w, r)
		case http.MethodPost:
			controller.CreateUser(w, r)
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUser(w, r)
		case http.MethodPut:
			controller.UpdateUser(w, r)
		case http.MethodDelete:
			controller.DeleteUser(w, r)
		}
	})

	log.Printf("Server running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}
