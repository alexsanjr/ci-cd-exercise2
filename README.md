# 🚀 API em Go

## 📌 Sobre o Projeto
API simples desenvolvida em GO como forma de treinamento para criação do dockerfile, compose e ci/cd

Pré-requisitos
- Docker Desktop (recomendado)
- (Opcional) Go instalado localmente

Passo a passo
- Sem Go instalado (via Docker): veja as seções de Docker/Compose abaixo.
- Com Go instalado:
  - `go mod tidy`
  - `go run main.go`

A API estará disponível em:
- http://localhost:3000

## 🐳 (Atividade 1) Rodando com Docker

Você deverá criar um Dockerfile para essa aplicação

### Requisitos:

- Utilizar imagem oficial do Go
- Build da aplicação
- Usar multi-stage build
- Expor porta 3000
- Executar a API

### Como rodar

- Build:
  - `docker build -t api-go:local .`
- Run:
  - `docker run --rm -p 3000:3000 api-go:local`

Acesse: `http://localhost:3000/users`


## 🐳 (Atividade 2) Rodando com Compose

Você deverá modificar a aplicação para fazer acesso ao banco de dados. Crie um docker compose para executar o PostreSQL, o PGAdmin e a aplicação em GO atraves do dockerfile que você criou

### Requisitos:

- Utilizar o dockerfile criado na atividade 1
- Criar um docker compose com:
  - A aplicação em go
  - O banco PostgreSQL
  - O PGAdmin

### Como rodar

- Subir tudo:
  - `docker compose up -d --build`

Serviços:
- API: `http://localhost:3000`
- PostgreSQL: `localhost:5432` (usuário `app`, senha `app`, database `app`)
- PGAdmin: `http://localhost:5050`
  - Email: `admin@admin.com`
  - Senha: `admin`

Para conectar no PGAdmin:
- Create Server
  - Host: `db`
  - Port: `5432`
  - Username: `app`
  - Password: `app`

Parar:
- `docker compose down`




## ⚙️ (DESAFIO) CI/CD

Crie um CI/CD no github actions com as seguintes etapas

- CI (Integração Contínua)
  - Build da aplicação
  - Testes unitários
  - Testes de integração
  - Lint
  - Análise de qualidade de código (SonarQube)
  - SAST (Semgrep ou Checkmarx ou Fortify, etc)

- Container
  - Docker Lint
  - Build da imagem
  - Scan de vulnerabilidades (Trivy)
  - Push da imagem no dockerhub

- CD (Entrega Contínua)
   - Deploy em homologação com Render
   - DAST (OWASP ZAP)
   - Criação da aprovação manual
   - Deploy em produção

### Workflow implementado

O workflow está em `.github/workflows/ci-cd.yml` e roda em `push`/`pull_request`.

Notas:
- Testes de integração rodam com a tag `integration` e precisam de Postgres (no CI é um service container).
- Passos de Sonar/Push/Deploy/DAST são **opcionais** e só rodam se você configurar os `secrets` no repositório.

Secrets suportados (opcionais):
- `DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`
- `SONAR_TOKEN`, `SONAR_HOST_URL`, `SONAR_PROJECT_KEY`
- `RENDER_HOMOLOG_DEPLOY_HOOK`, `HOMOLOG_BASE_URL`
- `RENDER_PROD_DEPLOY_HOOK`, `PROD_BASE_URL`

Para aprovação manual no deploy de produção, configure o Environment `production` no GitHub (Settings → Environments) com required reviewers.