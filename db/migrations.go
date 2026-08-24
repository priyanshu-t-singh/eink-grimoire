package migrations

import "embed"

// MigrationsFS holds all embedded SQL migration files.
//
//go:embed migrations/*.sql
var MigrationsFS embed.FS
