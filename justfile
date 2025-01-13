[working-directory: 'db']
db_code_gen:
	sqlite3 test.db < scheme.sql
	SQLITE_DSN=test.db go run github.com/stephenafamo/bob/gen/bobgen-sqlite@latest
	rm test.db

build:
    #!/bin/bash
    # Set environment variables to use Zig as the C compiler and linker
    export CC="zig cc -target x86_64-linux-musl"
    export CXX="zig c++ -target x86_64-linux-musl"
    export CGO_ENABLED=1
    export CGO_CFLAGS="-static"
    export CGO_LDFLAGS="-static"
    
    # Build the Go application
    go build .

ionic-build:
    cd web && ionic build --prod

run: build
    AUTH_USERNAME="test" AUTH_PASSWORD="test" ./vet_note -db backup1.db -port 8001 -cors -dbLog

runall: ionic-build build run