module github.com/developerc/reductorUrl

go 1.22.12

replace (
	golang.org/x/tools => golang.org/x/tools v0.31.0
)

require go.uber.org/zap v1.27.0

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

require (
	github.com/go-chi/chi/v5 v5.2.1
	github.com/gorilla/securecookie v1.1.2
	github.com/jackc/pgerrcode v0.0.0-20240316143900-6e2875d9b438
	github.com/jackc/pgx/v5 v5.7.2
	go.uber.org/multierr v1.10.0 // indirect
)
