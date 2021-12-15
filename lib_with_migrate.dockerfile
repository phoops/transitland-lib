FROM alpine:3.13

# Copy the migrations 

WORKDIR /migrations

COPY ./internal/schema/postgres/migrations /migrations

RUN apk add curl && curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz -C /usr/bin

CMD ["sh", "-c", "migrate -path /migrations -database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE} up"]