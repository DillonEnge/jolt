migrate:
	go run github.com/jackc/tern/v2@latest migrate -m database/migrations --database market-wise-api-dev ;

migrate_down:
	go run github.com/jackc/tern/v2@latest migrate -m database/migrations --database market-wise-api-dev -d -1 ;

sqlc:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate ;

templ:
	go run github.com/a-h/templ/cmd/templ@latest generate -watch -proxy="http://localhost:$(PORT)" -open-browser=false ;

tailwind:
	npx tailwindcss -i templates/static/input.css -o templates/static/output.css --watch=always

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

test:
	go test ./... ;

lint:
	golangci-lint run
