#!/bin/sh

cd web
ionic build --prod
cd ..
CGO_ENABLED=1 go build .
