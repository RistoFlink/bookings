#!/bin/bash

#go build -o bookings cmd/web/*.go && ./bookings
go run cmd/web/main.go cmd/web/middleware.go cmd/web/routes.go
