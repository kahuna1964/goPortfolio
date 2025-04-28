package migrations

import "embed"

// This file tells the compiler for build to include any file with .sql in this directory into the app

//go:embed *.sql
var FS embed.FS
