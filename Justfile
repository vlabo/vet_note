set windows-shell := ["pwsh.exe", "-NoProfile", "-NoLogo", "-Command"]
export AUTH_USERNAME := "test"
export AUTH_PASSWORD := "test"
export CGO_ENABLED := "0"

[working-directory: 'db']
db-code-gen $SQLITE_DSN="test.db":
	cat scheme.sql | sqlite3 test.db
	go run github.com/stephenafamo/bob/gen/bobgen-sqlite@latest
	rm test.db

[working-directory: 'ui']
ui-dev:
    bun run dev

[working-directory: 'ui']
ui-build:
    bun run build

tygo: 
    tygo generate

build:
    # Build the Go application
    go build .

build-linux $GOOS="linux":
    # Build the Go application
    go build .

run: build
     ./vet_note -db backup1.db -port 8001 -cors -dbLog


