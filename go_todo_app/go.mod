module github.com/koralle/go-web-application-development/go_todo_app

go 1.18

require golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4

require (
	github.com/caarlos0/env/v6 v6.9.3 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
)

replace github.com/koralle/go-web-application-development/go_todo_app/server => ../server

replace github.com/koralle/go-web-application-development/go_todo_app/config => ../config

replace github.com/kolalle/go-web-application-development/go_todo_app/store => ../store

replace github.com/kolalle/go-web-application-development/go_todo_app/entity => ../entity
