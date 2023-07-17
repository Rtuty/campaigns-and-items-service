package migrations

import "embed"

//go:embed postgres/*.sql
var FS embed.FS
