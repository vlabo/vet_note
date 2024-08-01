#!/bin/sh

cd web
ionic build
cd ..
go build .
