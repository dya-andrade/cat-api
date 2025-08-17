```markdown
# 🐱 Cat API - Go + PostgreSQL

Uma API moderna escrita em **Go** para cadastro de gatos, utilizando **PostgreSQL** como banco de dados.  
Projeto estruturado com boas práticas, paralelismo e pronto para escalar.

---

## 🚀 Tecnologias

- **Go 1.25+**
- **PostgreSQL 15+**
- **pgx** (driver performático para Postgres)
- **chi** (roteador HTTP leve e rápido)
- **golang-migrate** (migrações do banco)
- **Docker + docker-compose** (ambiente pronto)
- **errgroup / worker pool** (paralelismo)

---

## 📂 Estrutura

```

.
├── cmd/
│   └── api/            # entrypoint da aplicação
├── internal/
│   ├── cats/           # módulo principal (CRUD de gatos)
│   ├── db/             # conexão e queries
│   └── workers/        # exemplo de paralelismo
├── migrations/         # scripts SQL versionados
├── .env                # variáveis de ambiente
└── README.md

````

---

## ⚙️ Variáveis de Ambiente

Definidas no arquivo **`.env`** na raiz do projeto:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=catsdb
DB_SSLMODE=disable
````

---

## 🐘 Banco de Dados e Schemas

O projeto utiliza **dois schemas** no PostgreSQL:

* `public` → usado pelo sistema e bibliotecas internas.
* `cats` → schema **dedicado** para tabelas da aplicação (ex.: `cats.gatos`).

👉 Isso garante **organização e isolamento** das tabelas da aplicação sem poluir o schema `public`.

Exemplo de criação na migração inicial:

```sql
CREATE SCHEMA IF NOT EXISTS cats;

CREATE TABLE IF NOT EXISTS cats.gatos (
    id SERIAL PRIMARY KEY,
    nome TEXT NOT NULL,
    idade INT NOT NULL,
    raca TEXT NOT NULL,
    criado_em TIMESTAMP DEFAULT now()
);
```

---

## ▶️ Como rodar

### 1. Clone o repositório

```bash
git clone https://github.com/dya-andrade/cat-api.git
cd cat-api
```

### 2. Suba o Postgres com Docker

```bash
docker-compose up -d
```

### 3. Execute as migrações

```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/catsdb?sslmode=disable" up
```

### 4. Rode a aplicação

```bash
go run ./cmd/api
```

A API estará rodando em [http://localhost:8080](http://localhost:8080)

---

## 🛠 Endpoints

* `GET /cats` → lista todos os gatos
* `POST /cats` → cria um novo gato
* `GET /cats/{id}` → busca gato por ID
* `PUT /cats/{id}` → atualiza gato
* `DELETE /cats/{id}` → remove gato

Exemplo de `POST /cats`:

```json
{
  "nome": "Mingau",
  "idade": 2,
  "raca": "SRD"
}
```

---

## ⚡ Paralelismo

A aplicação utiliza:

* **errgroup** para rodar o servidor HTTP + workers em paralelo.
* **worker pool** para processar tarefas em background (exemplo: logs, notificações).

---

## ✅ Healthcheck

Disponível em:

```
GET /health
```

Retorna `200 OK` se a API e o banco estiverem funcionando.

---

## 📜 Licença

MIT

```
