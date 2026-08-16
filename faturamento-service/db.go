package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=db_faturamento sslmode=disable"

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
	CREATE TABLE IF NOT EXISTS notas (
		numero SERIAL PRIMARY KEY,
		status VARCHAR(20) NOT NULL DEFAULT 'Aberta'
	);

	CREATE TABLE IF NOT EXISTS itens_nota (
		id SERIAL PRIMARY KEY,
		nota_numero INTEGER NOT NULL REFERENCES notas(numero),
		produto_id INTEGER NOT NULL,
		quantidade INTEGER NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("erro ao criar tabelas: %v", err)
	}

	log.Println("tabelas notas e itens_nota prontas")
}