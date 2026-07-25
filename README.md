# finances-go

Comandos utiles para levantar base, correr migraciones y seeders sin instalar dependencias locales.

## Requisitos

- Docker
- Docker Compose
- Make

## Servicios

Levantar todo:

```bash
docker compose up -d --build
```

## Migraciones

Migrar con una sola orden (usa Docker, no requiere `migrate` local):

```bash
make migrate
```

Migrar con el target anterior (solo migraciones en `root`):

```bash
make migrate-up-docker
```

## Seeder

Seeder con una sola orden (usa Docker, no requiere `go` local):

```bash
make seed-docker
```

Migracion + seeder:

```bash
make migrate-seed
```

## Flujo completo

Bootstrap no destructivo (levanta DB + migra + seed):

```bash
make db-bootstrap
```

Reset total destructivo (borra volumenes y vuelve a crear todo):

```bash
make db-reset-seed
```

## Conexion en Adminer

- URL: `http://localhost:8081`
- Sistema: `PostgreSQL`
- Servidor: `db`
- Usuario: `root`
- Contrasena: `root`
- DB sugerida: `postgres` o `root` (segun donde ejecutes migraciones/seed)

## Comandos directos sin Make

Migrar `postgres`:

```bash
docker run --rm --network finances-go_default -v "$PWD/migrations:/migrations" migrate/migrate -path=/migrations -database "postgres://root:root@db:5432/postgres?sslmode=disable" up
```

Migrar `root`:

```bash
docker run --rm --network finances-go_default -v "$PWD/migrations:/migrations" migrate/migrate -path=/migrations -database "postgres://root:root@db:5432/root?sslmode=disable" up
```

Seeder en `postgres`:

```bash
docker run --rm --network finances-go_default -v "$PWD:/app" -w /app -e DATABASE_URL="postgres://root:root@db:5432/postgres?sslmode=disable" -e DATABASE_HOST=db -e DATABASE_PORT=5432 -e DATABASE_USER=root -e DATABASE_PASSWORD=root -e DATABASE_NAME=postgres golang:1.26-alpine go run ./cmd/seed
```

Seeder en `root`:

```bash
docker run --rm --network finances-go_default -v "$PWD:/app" -w /app -e DATABASE_URL="postgres://root:root@db:5432/root?sslmode=disable" -e DATABASE_HOST=db -e DATABASE_PORT=5432 -e DATABASE_USER=root -e DATABASE_PASSWORD=root -e DATABASE_NAME=root golang:1.26-alpine go run ./cmd/seed
```
