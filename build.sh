#!/bin/sh

cd web
ionic build --prod
cd ..
go build .
