package migrations

import "embed"

//go:embed postgres/*.sql
var FS embed.FS

//go:embed clickhouse/*.sql
var MFS embed.FS
