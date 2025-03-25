set windows-shell := ["pwsh.exe", "-NoProfile", "-NoLogo", "-Command"]
export AUTH_USERNAME := "test"
export AUTH_PASSWORD := "test"
export CGO_ENABLED := "0"

[working-directory: 'db']
db_code_gen:
	sqlite3 test.db < scheme.sql
	SQLITE_DSN=test.db go run github.com/stephenafamo/bob/gen/bobgen-sqlite@latest
	rm test.db

build:
    # Build the Go application
    go build .

[working-directory: 'ui']
ui-dev:
    bun run dev

[working-directory: 'ui']
ui-build:
    bun run build

ionic-build:
    cd web && ionic build --prod

run: build
     ./vet_note -db backup1.db -port 8001 -cors -dbLog

runall: ionic-build build run

