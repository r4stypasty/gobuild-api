package main

import (
	"fmt"
	"net/http"

	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/r4stypasty/gobuild-api/internal/handlers"
	"github.com/r4stypasty/gobuild-api/internal/repository"
)

func main() {

	db, err := repository.InitDB("materials.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()
	fmt.Println("Database connected successfully")

	materialHandler := &handlers.MaterialHandler{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("gobuild-api is running!"))
	})

	r.Get("/materials", materialHandler.GetMaterials)
	r.Post("/materials", materialHandler.CreateMaterial)
	r.Post("/materials/{id}/write-off", materialHandler.WriteOffMaterial)
	r.Get("/alerts", materialHandler.GetAlerts)

	fmt.Println("Server is starting on http://localhost:8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
