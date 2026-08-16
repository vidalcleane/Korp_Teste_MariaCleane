package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=db_estoque sslmode=disable"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("erro ao conectar no banco: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("erro ao pingar o banco: %v", err)
	}

	migrate()
}

func migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS produtos (
		id SERIAL PRIMARY KEY,
		codigo VARCHAR(50) UNIQUE NOT NULL,
		descricao VARCHAR(255) NOT NULL,
		saldo INTEGER NOT NULL DEFAULT 0
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("erro ao criar tabela produtos: %v", err)
	}

	log.Println("tabela produtos pronta")
}
