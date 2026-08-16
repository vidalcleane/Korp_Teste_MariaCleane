package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func criarProduto(w http.ResponseWriter, r *http.Request) {
	var p Produto
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondErro(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if p.Codigo == "" || p.Descricao == "" {
		respondErro(w, http.StatusBadRequest, "codigo e descricao são obrigatórios")
		return
	}

	query := `INSERT INTO produtos (codigo, descricao, saldo) VALUES ($1, $2, $3) RETURNING id`
	if err := db.QueryRow(query, p.Codigo, p.Descricao, p.Saldo).Scan(&p.ID); err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao salvar produto: "+err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, p)
}

func listarProdutos(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, codigo, descricao, saldo FROM produtos ORDER BY id`)
	if err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao consultar produtos")
		return
	}
	defer rows.Close()

	produtos := []Produto{}
	for rows.Next() {
		var p Produto
		if err := rows.Scan(&p.ID, &p.Codigo, &p.Descricao, &p.Saldo); err != nil {
			respondErro(w, http.StatusInternalServerError, "erro ao ler produto")
			return
		}
		produtos = append(produtos, p)
	}

	respondJSON(w, http.StatusOK, produtos)
}

// baixarEstoque dá baixa no saldo de um produto de forma segura para concorrência.
// Usa uma transação com "SELECT ... FOR UPDATE", que trava a linha do produto
// no banco até a transação terminar. Assim, se duas notas tentarem baixar o
// mesmo produto ao mesmo tempo, a segunda espera a primeira terminar em vez
// de ler um saldo desatualizado.
func baixarEstoque(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req BaixaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Quantidade <= 0 {
		respondErro(w, http.StatusBadRequest, "quantidade inválida")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao iniciar transação")
		return
	}
	defer tx.Rollback()

	var saldoAtual int
	err = tx.QueryRow(`SELECT saldo FROM produtos WHERE id = $1 FOR UPDATE`, id).Scan(&saldoAtual)
	if errors.Is(err, sql.ErrNoRows) {
		respondErro(w, http.StatusNotFound, "produto não encontrado")
		return
	}
	if err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao consultar saldo")
		return
	}

	if saldoAtual < req.Quantidade {
		respondErro(w, http.StatusConflict, "saldo insuficiente")
		return
	}

	novoSaldo := saldoAtual - req.Quantidade
	if _, err := tx.Exec(`UPDATE produtos SET saldo = $1 WHERE id = $2`, novoSaldo, id); err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao atualizar saldo")
		return
	}

	if err := tx.Commit(); err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao confirmar transação")
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"saldo": novoSaldo})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondErro(w http.ResponseWriter, status int, mensagem string) {
	respondJSON(w, status, map[string]string{"erro": mensagem})
}