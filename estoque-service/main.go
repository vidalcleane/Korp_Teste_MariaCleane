package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// corsMiddleware libera o navegador (rodando o Angular em localhost:4200)
// a fazer requisições diretamente para este serviço.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	initDB()

	r := chi.NewRouter()
	r.Use(corsMiddleware)

	r.Post("/produtos", criarProduto)
	r.Get("/produtos", listarProdutos)
	r.Post("/produtos/{id}/baixa", baixarEstoque)

	log.Println("serviço de estoque rodando na porta 8081")
	log.Fatal(http.ListenAndServe(":8081", r))
}
