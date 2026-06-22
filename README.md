# Mapa das Sementes do Brasil

> Plataforma colaborativa para catalogar, mapear e preservar sementes nativas brasileiras e o conhecimento tradicional a elas associado.

---

## Sobre o Projeto

O **Mapa das Sementes do Brasil** é uma API REST desenvolvida em Go para documentar a biodiversidade de sementes crioulas e nativas do Brasil, conectando guardiões de sementes, especialistas e pesquisadores em uma base de dados aberta e colaborativa.

### Módulos

| Módulo | Descrição |
|---|---|
| Mapeamento Geográfico | Registros geolocalizados de ocorrência de sementes |
| Catálogo de Sementes | Fichas completas com características botânicas |
| Conhecimento Tradicional | Saberes de comunidades e guardiões de sementes |
| Registro Fotográfico | Upload e galeria de imagens das sementes |
| Validação Colaborativa | Avaliação por especialistas e comunidade |
| Dashboard de Dados | Visualização e estatísticas da coleção |

## Stack Tecnológica - Ferramentas utilizadas para o desenvolvimento

- **Linguagem:** Go 1.25
- **Framework:** Gin
- **Banco de Dados:** PostgreSQL + GORM
- **Autenticação:** JWT (Bearer Token)
- **Deploy:** Railway
- **Containerização:** Docker
- **VS Code 1.125.0:** IDE principal, com extensões Go (gopls, delve)
- **pgAdmin 4:** Administração e modelagem do banco PostgreSQL
- **GitHub:** Versionamento, histórico de etapas e colaboração
- **Postman:** Testes dos endpoints REST com autenticação JWT
- **Railway:** Deploy contínuo com variáveis de ambiente protegidas

---

## Estrutura do Projeto

```
mapa-sementes-brasil/
├── config/         # Configurações de ambiente
├── database/       # Conexão e migrações do banco
├── docs/           # Documentação da API e do catálogo
├── handlers/       # Controllers das rotas HTTP
├── middleware/     # Auth JWT, CORS, logging
├── models/         # Modelos de dados (GORM)
├── scripts/        # Scripts utilitários
├── services/       # Lógica de negócio e integrações
├── static/         # Frontend estático (HTML/CSS/JS)
├── uploads/        # Imagens enviadas pelos usuários
└── utils/          # Funções auxiliares
```

---

## Como Rodar Localmente

### Pré-requisitos

- Go 1.21+
- PostgreSQL 14+
- Git

### Instalação

```bash
# 1. Clone o repositório
git clone https://github.com/seu-usuario/mapa-sementes-brasil.git
cd mapa-sementes-brasil

# 2. Configure as variáveis de ambiente
cp .env.example .env
# Edite o .env com suas credenciais

# 3. Instale as dependências
go mod tidy

# 4. Execute
go run main.go
```

### Com Docker

```bash
docker-compose up --build
```

A API estará disponível em `http://localhost:8080`

---

## Autenticação

A API usa **JWT Bearer Token**. Para acessar rotas protegidas:

```bash
# 1. Registre um usuário
POST /api/auth/registro

# 2. Faça login
POST /api/auth/login
# Retorna: { "token": "eyJ..." }

# 3. Use o token no header
Authorization: Bearer eyJ...
```
---

## Principais Endpoints

| Método | Rota | Descrição |
|---|---|---|
| `POST` | `/api/auth/registro` | Registrar usuário |
| `POST` | `/api/auth/login` | Login |
| `GET` | `/api/especies` | Listar espécies |
| `POST` | `/api/especies` | Cadastrar espécie |
| `GET` | `/api/sementes` | Listar sementes |
| `POST` | `/api/sementes` | Cadastrar semente |
| `POST` | `/api/sementes/:id/imagem` | Upload de imagem |
| `GET` | `/api/registros` | Registros geolocalizados |
| `POST` | `/api/registros` | Novo registro com coordenadas |
| `GET` | `/api/conhecimentos` | Conhecimentos tradicionais |
| `POST` | `/api/avaliacoes` | Avaliar registro |

---


## Contribuindo
Em breve!!!
Consulte o [`docs/CONTRIBUTING.md`](docs/CONTRIBUTING.md) para diretrizes de contribuição.

---

## Licença

Este projeto é open source sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

---

> Desenvolvido para preservar a biodiversidade brasileira.

---

