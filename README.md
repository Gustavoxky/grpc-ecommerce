# 📆 grpc-ecommerce

Sistema de pedidos com autenticação de usuários, baseado em **gRPC + Go + PostgreSQL**, pronto para deploy em **Docker**, **Kubernetes** e **GitHub Actions** com **Skaffold**.

> Autenticação via JWT • CLI interativo • Interceptores gRPC • PostgreSQL via Docker • Kubernetes com Skaffold • CI/CD com GitHub Actions

---

## ✅ Funcionalidades

- Registro e login de usuários
- Autenticação com JWT
- CRUD completo de pedidos
- Cliente CLI interativo
- Validação de token em todas as chamadas
- Logs de requisições com método, tempo, status e usuário
- Banco de dados PostgreSQL via Docker
- Deploy automatizado com Kubernetes + Skaffold
- Pipeline CI/CD via GitHub Actions

---

## 🧐 Tecnologias

- Go 1.22+
- gRPC
- Protocol Buffers
- PostgreSQL (via Docker Compose)
- JWT (`golang-jwt/jwt/v5`)
- Bcrypt (`x/crypto/bcrypt`)
- Docker & Docker Compose
- Kubernetes
- Skaffold
- GitHub Actions

---

## 🚀 Executando o projeto

### 1. Clonar o repositório

```bash
git clone https://github.com/seu-usuario/grpc-ecommerce.git
cd grpc-ecommerce
```

### 2. Subir o PostgreSQL com Docker

```bash
docker-compose up -d
```

Banco:
- Host: `localhost`
- Porta: `5432`
- User: `grpcuser`
- Pass: `grpcpass`
- DB: `grpcdb`

### 3. Instalar as dependências do Go

```bash
go mod tidy
```

---

## 💠 Gerar os arquivos gRPC (se alterar o .proto)

```bash
protoc --go_out=. --go-grpc_out=. proto/order.proto
```

> Isso atualiza os arquivos em `pb/`

---

## ▶️ Rodar o servidor

```bash
make run-server
```

## 💻 Rodar o cliente CLI

```bash
make run-client
```

---

## ☸️ Deploy com Kubernetes + Skaffold

1. Instale o [Skaffold](https://skaffold.dev/docs/install/):

```bash
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-amd64
chmod +x skaffold && sudo mv skaffold /usr/local/bin/
```

2. Inicie um cluster local com [Minikube](https://minikube.sigs.k8s.io):

```bash
minikube start
```

3. Aplique o deploy:

```bash
skaffold dev
```

> Isso builda a imagem, aplica os manifests em `k8s/` e faz port-forward para `localhost:50051`

---

## 🧾 Logs de requisições

Logs exibem:

- Método gRPC
- Tempo de resposta
- Status (OK ou erro)
- Usuário autenticado (ou anônimo)

Exemplo:
```
[gRPC] método: /order.OrderService/CreateOrder | usuário: gustavo | duração: 14.1ms | status: OK
```

---

## 📁 Estrutura

```
grpc-ecommerce/
├── client/                # Cliente CLI
├── proto/                 # Arquivo .proto
├── pb/                    # Arquivos gerados pelo protoc
├── server/                # Lógica do servidor (auth, db, handlers)
├── k8s/                   # Manifestos Kubernetes
├── main.go                # Start do servidor
├── docker-compose.yml
├── Dockerfile             # Build do servidor Go
├── Makefile               # Comandos de atalho
├── skaffold.yaml          # Configuração de dev com Skaffold
├── .gitignore
└── README.md
```

---

## 🔁 CI/CD com GitHub Actions

O deploy automático é feito com:

```yaml
.github/workflows/deploy.yaml
```

Você precisa configurar os segredos:

- `GCP_PROJECT_ID`: ID do projeto GCP ou namespace
- `KUBECONFIG`: conteúdo do kubeconfig em base64

---

## 🧺 Usuário teste (opcional)

Você pode registrar um novo usuário via CLI, ou usar o usuário de exemplo:

- **Usuário**: gustavo  
- **Senha**: 12345

---

## 📌 Requisitos

- Go 1.22+
- Docker + Docker Compose
- `protoc` + plugins `protoc-gen-go`, `protoc-gen-go-grpc`
- Kubernetes local (ex: Minikube)
- Skaffold

---

## ✨ Próximos passos

- gRPC-Web + frontend em React/Next.js
- Exportação de logs para arquivo ou banco
- Autorização por nível de usuário (admin, operador)
- Deploy multiambiente (prod, staging, dev)

---

## 🧑‍💻 Autor

Desenvolvido por [Seu Nome]  
Contribuições, PRs e sugestões são bem-vindas!
