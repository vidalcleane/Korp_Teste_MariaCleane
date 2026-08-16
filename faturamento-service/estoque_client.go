package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// URL do Serviço de Estoque. Em um ambiente real isso viria de variável de ambiente.
const estoqueBaseURL = "http://localhost:8081"

var httpClient = &http.Client{
	// Timeout evita que o Faturamento fique travado para sempre se o
	// Estoque estiver fora do ar — é a base do tratamento de falha.
	Timeout: 5 * time.Second,
}

type baixaEstoquePayload struct {
	Quantidade int `json:"quantidade"`
}

// baixarEstoqueRemoto chama o Serviço de Estoque para dar baixa em um produto.
// Retorna um erro descritivo se o serviço estiver fora do ar, com timeout,
// ou se o próprio Estoque recusar a operação (ex: saldo insuficiente).
func baixarEstoqueRemoto(produtoID int, quantidade int) error {
	payload := baixaEstoquePayload{Quantidade: quantidade}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/produtos/%d/baixa", estoqueBaseURL, produtoID)

	resp, err := httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("serviço de estoque indisponível: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var erroResp map[string]string
		json.NewDecoder(resp.Body).Decode(&erroResp)
		return fmt.Errorf("estoque recusou a baixa do produto %d: %s", produtoID, erroResp["erro"])
	}

	return nil
}