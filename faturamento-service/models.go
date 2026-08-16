package main

type ItemNota struct {
	ProdutoID  int `json:"produto_id"`
	Quantidade int `json:"quantidade"`
}

type Nota struct {
	Numero int        `json:"numero"`
	Status string     `json:"status"`
	Itens  []ItemNota `json:"itens"`
}

type CriarNotaRequest struct {
	Itens []ItemNota `json:"itens"`
}