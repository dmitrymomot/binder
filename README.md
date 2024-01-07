# binder

HTTP request data binder.

[![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/dmitrymomot/binder)](https://github.com/dmitrymomot/binder)
[![Tests](https://github.com/dmitrymomot/binder/actions/workflows/tests.yml/badge.svg)](https://github.com/dmitrymomot/binder/actions/workflows/tests.yml)
[![CodeQL Analysis](https://github.com/dmitrymomot/binder/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/dmitrymomot/binder/actions/workflows/codeql-analysis.yml)
[![GolangCI Lint](https://github.com/dmitrymomot/binder/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/dmitrymomot/binder/actions/workflows/golangci-lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/dmitrymomot/binder)](https://goreportcard.com/report/github.com/dmitrymomot/binder)
[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrymomot/binder.svg)](https://pkg.go.dev/github.com/dmitrymomot/binder)
[![License](https://img.shields.io/github/license/dmitrymomot/binder)](https://github.com/dmitrymomot/binder/blob/main/LICENSE)

## Features

- [x] Bind query string parameters to struct fields
- [x] Bind form values to struct fields
- [x] Bind JSON body to struct fields
- [x] Get file from multipart form
- [x] Bind multipart form values to struct fields (limited support, see [supported types](#supported-types))
- [x] Binder interface implementation

### Supported types

Supported types for binding from multipart form:

- [x] `string`
- [x] `int`, `int8`, `int16`, `int32`, `int64`
- [x] `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- [x] `float32`, `float64`
- [x] `bool`
- [x] `map[string]interface{}`
- [x] `*binder.File` & `binder.File`
- [ ] `time.Time`
- [ ] `[]string`
- [ ] `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`
- [ ] `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`
- [ ] `[]float32`, `[]float64`
- [ ] `[]bool`
- [ ] `[]*binder.File` & `[]binder.File`
- [ ] `[]time.Time`

## Installation

```bash
go get -u github.com/dmitrymomot/binder
```

## Usage

Comming soon...