package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func criarNota(w http.ResponseWriter, r *http.Request) {
	var req CriarNotaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErro(w, http.StatusBadRequest, "corpo da requisição inválido")
		return
	}

	if len(req.Itens) == 0 {
		respondErro(w, http.StatusBadRequest, "a nota precisa ter ao menos um produto")
		return
	}

	tx, err := db.Begin()
	if err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao iniciar transação")
		return
	}
	defer tx.Rollback()

	var numero int
	if err := tx.QueryRow(`INSERT INTO notas (status) VALUES ('Aberta') RETURNING numero`).Scan(&numero); err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao criar nota")
		return
	}

	for _, item := range req.Itens {
		_, err := tx.Exec(
			`INSERT INTO itens_nota (nota_numero, produto_id, quantidade) VALUES ($1, $2, $3)`,
			numero, item.ProdutoID, item.Quantidade,
		)
		if err != nil {
			respondErro(w, http.StatusInternalServerError, "erro ao adicionar item na nota")
			return
		}
	}

	if err := tx.Commit(); err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao confirmar criação da nota")
		return
	}

	respondJSON(w, http.StatusCreated, Nota{Numero: numero, Status: "Aberta", Itens: req.Itens})
}

func listarNotas(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT numero, status FROM notas ORDER BY numero`)
	if err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao consultar notas")
		return
	}
	defer rows.Close()

	notas := []Nota{}
	for rows.Next() {
		var n Nota
		if err := rows.Scan(&n.Numero, &n.Status); err != nil {
			respondErro(w, http.StatusInternalServerError, "erro ao ler nota")
			return
		}
		n.Itens = buscarItens(n.Numero)
		notas = append(notas, n)
	}

	respondJSON(w, http.StatusOK, notas)
}

func buscarItens(numeroNota int) []ItemNota {
	itens := []ItemNota{}
	rows, err := db.Query(`SELECT produto_id, quantidade FROM itens_nota WHERE nota_numero = $1`, numeroNota)
	if err != nil {
		return itens
	}
	defer rows.Close()

	for rows.Next() {
		var it ItemNota
		rows.Scan(&it.ProdutoID, &it.Quantidade)
		itens = append(itens, it)
	}
	return itens
}

// imprimirNota é o coração do requisito de "tratamento de falhas":
// 1. Confere que a nota está Aberta.
// 2. Para cada item, chama o Serviço de Estoque para dar baixa.
// 3. Se o Estoque falhar (fora do ar, timeout, ou saldo insuficiente),
//    a nota continua Aberta e o usuário recebe um erro claro.
// 4. Só se TODAS as baixas derem certo, a nota vira Fechada.
func imprimirNota(w http.ResponseWriter, r *http.Request) {
	numero := chi.URLParam(r, "numero")

	var status string
	err := db.QueryRow(`SELECT status FROM notas WHERE numero = $1`, numero).Scan(&status)
	if errors.Is(err, sql.ErrNoRows) {
		respondErro(w, http.StatusNotFound, "nota não encontrada")
		return
	}
	if err != nil {
		respondErro(w, http.StatusInternalServerError, "erro ao consultar nota")
		return
	}

	if status != "Aberta" {
		respondErro(w, http.StatusConflict, "só é possível imprimir notas com status Aberta")
		return
	}

	var numeroInt int
	db.QueryRow(`SELECT numero FROM notas WHERE numero = $1`, numero).Scan(&numeroInt)
	itens := buscarItens(numeroInt)

	for _, item := range itens {
		if err := baixarEstoqueRemoto(item.ProdutoID, item.Quantidade); err != nil {
			respondErro(w, http.StatusBadGateway,
				"não foi possível imprimir a nota: "+err.Error())
			return
		}
	}

	if _, err := db.Exec(`UPDATE notas SET status = 'Fechada' WHERE numero = $1`, numero); err != nil {
		respondErro(w, http.StatusInternalServerError, "baixa feita no estoque, mas erro ao fechar a nota")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "Fechada"})
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondErro(w http.ResponseWriter, status int, mensagem string) {
	respondJSON(w, status, map[string]string{"erro": mensagem})
}