package main

//go:generate go tool ogen --target internal/adapters/oas --package oas --clean ../../docs/api/openapi.yml
//
//go:generate cp ../../docs/api/openapi.yml ./internal/adapters/api/openapi_copy.yml
