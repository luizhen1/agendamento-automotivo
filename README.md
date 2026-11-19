# Sistema de Agendamento de Serviços Automotivos

Monorepo com API REST em Go, persistência MySQL, dashboard Next.js e infraestrutura Docker Compose.

## Arquitetura do backend

O backend está organizado em camadas dentro de `backend/internal`:

- `domain`: entidades e regras simples do domínio, sem dependência de HTTP ou GORM.
- `application`: casos de uso e interfaces de repositório necessárias à aplicação.
- `httpapi`: rotas, DTOs, validação de entrada e tradução de erros para respostas HTTP.
- `infrastructure/mysql`: modelos GORM, migrações, seed e implementação dos repositórios.
- `bootstrap`: composição das dependências, conexão com o banco e configuração do servidor.

O fluxo principal é `HTTP → application → domain`; a aplicação depende de contratos e o adaptador MySQL os implementa. `backend/main.go` fica responsável apenas por iniciar a aplicação.

## Requisitos

- Docker e Docker Compose
- Node.js 18 ou superior e npm

## Subir banco e API

Na raiz do repositório:

```bash
docker compose up --build -d
```

O MySQL fica disponível em `localhost:3306` e a API em `http://localhost:8080`. O banco aplica as migrações GORM e cria três clientes e três serviços de exemplo quando as respectivas tabelas estão vazias.

Para acompanhar os logs:

```bash
docker compose logs -f backend
```

## Subir o front-end

Em outro terminal:

```bash
cd frontend
npm install && npm run dev
```

Acesse `http://localhost:3000`. Por padrão, o painel consome a API em `http://localhost:8080`. Para outro endereço, defina `NEXT_PUBLIC_API_URL` antes de iniciar o Next.js.

## Endpoints

- `GET, POST /clientes`
- `GET, POST /servicos`
- `DELETE /servicos/{id}`
- `GET, POST /agendamentos`
- `PATCH, DELETE /agendamentos/{id}`

Agendamentos em horários já ocupados retornam HTTP 409. Agendamentos cancelados não bloqueiam novos horários. O CORS permite chamadas do front-end em `http://localhost:3000`.
Serviços vinculados a agendamentos não podem ser excluídos; a API retorna HTTP 409 para preservar a integridade referencial.

## Encerrar

```bash
docker compose down
```

Para remover também os dados persistidos no volume do MySQL, use `docker compose down -v`.
