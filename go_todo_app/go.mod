module github.com/koralle/go-web-application-development/go_todo_app

go 1.18

require (
	github.com/caarlos0/env/v6 v6.9.3
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-playground/validator/v10 v10.11.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/go-cmp v0.5.8
	github.com/jmoiron/sqlx v1.3.5
	github.com/matryer/moq v0.2.7
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4
)

require (
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	golang.org/x/crypto v0.0.0-20220817201139-bc19a97f63c8 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sys v0.0.0-20220825204002-c680a09ffe64 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.12 // indirect
)

replace github.com/koralle/go-web-application-development/go_todo_app/server => ../server

replace github.com/koralle/go-web-application-development/go_todo_app/config => ../config

replace github.com/koralle/go-web-application-development/go_todo_app/store => ../store

replace github.com/koralle/go-web-application-development/go_todo_app/entity => ../entity
