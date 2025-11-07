# Books API (Go + Chi + Postgres)

API para estudos de Go construindo um CRUD de livros com Chi, Postgres, migrações com Goose, Makefile, Docker e live reload com Air.

## Stack

- Go (Chi Router, CORS, net/http)
- PostgreSQL (pgx stdlib)
- Goose (migrações)
- Docker Compose (banco)
- Testcontainers (integração DB)
- Air (live reload)
- Makefile (atalhos de desenvolvimento)

Arquivos principais:
- [cmd/api/main.go](cmd\api\main.go)
- Server e rotas: 
  - [internal/server/server.go](internal\server\server.go)
  - [internal/server/routes.go](internal\server\routes.go)
  - Teste: [internal/server/routes_test.go](internal\server\routes_test.go)
- Camada de livros (arquitetura Repository → Service → Handler):
  - Handler: [internal/book/handler.go](internal\book\handler.go)
  - Service: [internal/book/service.go](internal\book\service.go)
  - Repository: [internal/book/repository.go](internal\book\repository.go)
  - Modelos/DTOs: [internal/book/model.go](internal\book\model.go)
- Banco de dados:
  - Serviço DB: [internal/database/database.go](internal\database\database.go)
  - Migrações: [internal/database/migrations/](internal\database\migrations)
  - Testes de integração: [internal/database/database_test.go](internal\database\database_test.go)
- Utilitários de resposta: [util/response.go](util\response.go)
- Infra:
  - Compose: [docker-compose.yml](docker-compose.yml)
  - Makefile: [Makefile](Makefile)
  - Air: [.air.toml](.air.toml)
  - Módulos: [go.mod](go.mod)

## Pré-requisitos

- Go instalado (recomendado Go 1.22+)
- Docker Desktop
- Make (opcional — você pode usar os comandos “go run” e “docker compose” diretamente)
- Ferramentas de CLI:
  - Goose: `go install github.com/pressly/goose/v3/cmd/goose@latest`
  - Air (opcional — o `make watch` instala se faltar)

## Configuração

Crie um arquivo `.env` na raiz com as variáveis esperadas em [internal/database/database.go](internal\database\database.go) e [docker-compose.yml](docker-compose.yml):

```env
# Porta HTTP da API
PORT=8080

# Postgres (Docker)
DB_HOST=localhost
DB_PORT=5432
DB_USERNAME=postgres
DB_PASSWORD=postgres
DB_DATABASE=books
DB_SCHEMA=public
```

Dica: O docker-compose mapeia a porta como "${DB_PORT}:5432". Se mudar DB_PORT no .env, ajuste os clientes.

## Subindo o banco (Docker)

- Via Make:
```bash
make docker-run
```

- Ou diretamente:
```bash
docker compose up --build
```

## Migrações (Goose)

- Instalar (uma vez): 
```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

- Aplicar migrações:
```bash
make migrate-up
```

- Reverter última:
```bash
make migrate-down
```

As migrações estão em [internal/database/migrations](internal\database\migrations), ex.: [20251105011307_users.sql](internal\database\migrations\20251105011307_users.sql) (tabela books).

## Rodando a API

- Compilar:
```bash
make build
```

- Rodar:
```bash
make run
```

- Live reload (Windows, via Air):
```bash
make watch
```

Endpoints servidos em `http://localhost:${PORT}`.

## Endpoints

- GET `/` — Hello World JSON
- GET `/health` — Status da conexão com o banco
- Base de livros: `/api/books`
  - POST `/` — cria livro
  - PUT `/` — atualiza livro
  - GET `/{id}` — busca por ID
  - DELETE `/{id}` — remove por ID

Exemplos (PowerShell):

```powershell
# Hello
curl http://localhost:8080/

# Health
curl http://localhost:8080/health

# Criar livro
$body = @{
  title="Clean Architecture"
  author="Robert C. Martin"
  published="2017-09-20"   # FlexibleDate aceita vários formatos
  image="https://example.com/img.png"
  description="A book about software architecture"
} | ConvertTo-Json
curl -Method POST -Uri http://localhost:8080/api/books/ -ContentType "application/json" -Body $body

# Atualizar livro
$update = @{
  id="PUT-UUID-AQUI"
  title="Clean Architecture (2nd)"
  author="Robert C. Martin"
  published="20/09/2017"
  image="https://example.com/img2.png"
  description="Updated description"
} | ConvertTo-Json
curl -Method PUT -Uri http://localhost:8080/api/books/ -ContentType "application/json" -Body $update

# Buscar por ID
curl http://localhost:8080/api/books/PUT-UUID-AQUI

# Remover por ID
curl -Method DELETE http://localhost:8080/api/books/PUT-UUID-AQUI
```

Observações:
- A API usa CORS liberado para http://* e https://* em [internal/server/routes.go](internal\server\routes.go).
- O tipo `FlexibleDate` aceita formatos comuns (YYYY-MM-DD, DD/MM/YYYY, MM/DD/YYYY, ISO8601, RFC3339) — veja [internal/book/model.go](internal\book\model.go).

## Testes

- Testes gerais:
```bash
make test
```

- Testes de integração (Postgres via Testcontainers):
```bash
make itest
```

Pré-requisito: Docker em execução. Os testes sobem um Postgres temporário — veja [internal/database/database_test.go](internal\database\database_test.go).

## Arquitetura

- Roteamento e middlewares: [internal/server/routes.go](internal\server\routes.go)
- Handler (HTTP) → Service (regras/timeout) → Repository (SQL):
  - Handler: [internal/book/handler.go](internal\book\handler.go)
  - Service (timeouts, validações): [internal/book/service.go](internal\book\service.go)
  - Repository (SQL puro): [internal/book/repository.go](internal\book\repository.go)
- Conexão DB e health: [internal/database/database.go](internal\database\database.go)

## Troubleshooting

- goose: “command not found” — instale com `go install ...` e garanta que `$GOPATH/bin` está no PATH.
- Falha ao conectar no DB — confira variáveis do `.env` e se o container está no ar (`make docker-run`).
- Air no Windows — `make watch` instala automaticamente se não existir.

## Comandos Make úteis

```bash
make all         # build + test
make build       # compila binário
make run         # roda a API
make watch       # live reload com Air
make test        # testes
make itest       # testes de integração DB
make docker-run  # sobe Postgres
make docker-down # derruba Postgres
make migrate-up  # aplica migrações
make migrate-down# reverte migrações
make clean       # limpa binários
```

---
Projeto para estudo. Sinta-se livre para adaptar e evoluir.