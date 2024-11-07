init_db:
	psql postgresql://devuser:$(DATABASE_PASSWORD)@localhost:5433/devuser?sslmode=disable -f ./database/db_scripts/init_db.sql

exec_db:
	docker exec -it jolt-db-1 psql -U postgres

migrate:
	go run github.com/jackc/tern/v2@latest migrate -m database/migrations --database jolt_dev ;

migrate_down:
	go run github.com/jackc/tern/v2@latest migrate -m database/migrations --database jolt_dev -d -1 ;

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate ;

templ:
	go run github.com/a-h/templ/cmd/templ@latest generate -watch -proxy="http://localhost:$(PORT)" -open-browser=false ;

templ_no_watch:
	go run github.com/a-h/templ/cmd/templ@latest generate;

tailwind:
	npx tailwindcss -i templates/static/input.css -o templates/static/output.css --watch=always --minify

tailwind_no_watch:
	npx tailwindcss -i templates/static/input.css -o templates/static/output.css --minify

sync_static:
	go run github.com/air-verse/air@latest \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "templates/static" \
	--build.include_ext "css" ;

air:
	go run github.com/air-verse/air@latest ;

dev:
	make -j3 templ tailwind air

deps:
	docker compose up nats -d

test:
	go test ./... ;

lint:
	golangci-lint run

regen: templ_no_watch tailwind_no_watch
