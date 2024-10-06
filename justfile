
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
