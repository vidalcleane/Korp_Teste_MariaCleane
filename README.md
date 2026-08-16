# Korp ERP — Sistema de Emissão de Notas Fiscais

Projeto técnico desenvolvido para o desafio da Korp ERP: um sistema de emissão de notas fiscais com arquitetura de microsserviços, backend em Go e frontend em Angular.

## Arquitetura

O sistema é composto por três partes independentes:

```
frontend/              → Angular (porta 4200)
estoque-service/        → Go — cadastro de produtos e controle de saldo (porta 8081)
faturamento-service/    → Go — cadastro e impressão de notas fiscais (porta 8082)
```

Cada microsserviço tem seu próprio banco de dados PostgreSQL (`db_estoque` e `db_faturamento`). O Faturamento se comunica com o Estoque via HTTP/REST para dar baixa nos produtos ao imprimir uma nota.

## Funcionalidades

- Cadastro de produtos (código, descrição, saldo)
- Cadastro de notas fiscais com numeração sequencial e múltiplos produtos
- Impressão de nota fiscal, com indicador de processamento, baixa automática de estoque e mudança de status (Aberta → Fechada)
- Tratamento de falha: se o Serviço de Estoque estiver indisponível, o usuário recebe um erro claro e a nota permanece "Aberta" para nova tentativa
- Tratamento de concorrência: baixa de estoque protegida por transação (`SELECT ... FOR UPDATE`), evitando saldo negativo em disputas simultâneas

## Tecnologias

- **Backend:** Go, [chi](https://github.com/go-chi/chi) (roteador HTTP), [lib/pq](https://github.com/lib/pq) (driver PostgreSQL)
- **Frontend:** Angular, Angular Material, RxJS, Reactive Forms
- **Banco de dados:** PostgreSQL

## Como executar

### 1. Banco de dados

Com Docker instalado:

```bash
docker compose up -d
```

Isso sobe um PostgreSQL com os bancos `db_estoque` e `db_faturamento` já criados.

> Alternativa sem Docker: instale o PostgreSQL localmente e crie os dois bancos manualmente (`CREATE DATABASE db_estoque;` e `CREATE DATABASE db_faturamento;`), com usuário `postgres` e senha `postgres`.

### 2. Serviço de Estoque

```bash
cd estoque-service
go mod tidy
go run .
```

Sobe em `http://localhost:8081`.

### 3. Serviço de Faturamento

```bash
cd faturamento-service
go mod tidy
go run .
```

Sobe em `http://localhost:8082`.

### 4. Frontend

```bash
cd frontend
npm install
ng serve
```

Acesse em `http://localhost:4200`.

## Rotas

### Serviço de Estoque (`:8081`)

| Método | Rota | Descrição |
|---|---|---|
| POST | `/produtos` | Cadastra um novo produto |
| GET | `/produtos` | Lista todos os produtos |
| POST | `/produtos/{id}/baixa` | Dá baixa em uma quantidade do saldo |

### Serviço de Faturamento (`:8082`)

| Método | Rota | Descrição |
|---|---|---|
| POST | `/notas` | Cria uma nova nota fiscal (status Aberta) |
| GET | `/notas` | Lista todas as notas fiscais |
| POST | `/notas/{numero}/imprimir` | Imprime a nota: baixa o estoque e fecha a nota |

## Documentação técnica

O detalhamento técnico completo (ciclos de vida do Angular, uso de RxJS, bibliotecas, tratamento de erros, etc.) está no arquivo `Detalhamento_Tecnico_Korp.docx`.
