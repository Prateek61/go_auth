# Go Auth

User authentication implemened using graphql with GO gqlgen.

## Run Locally

### Setup environment variables

```
DB_URL= database url
JWT_SECRET= secret key for jwt
```

### Install dependencies

```bash
go mod download
```

### Run migrations

```bash
migrate -path ./postgres/migrations -database "$DB_URL" up
```

### Run server

```bash
go run server.go
```