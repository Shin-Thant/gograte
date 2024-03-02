package gograte

var UsageMessage = `Usage:
// Init
gograte init

// Create
gograte create [migration_name]

// Migrate
gograte migrate [db_driver] [db_url] [migrate_action]

// Status
gograte status [db_driver] [db_url]

// Examples
gograte status postgres "postgres://localhost:5432/db"

gograte migrate postgres "postgres://localhost:5432/db" up
gograte migrate postgres "postgres://localhost:5432/db" down
gograte migrate postgres "postgres://localhost:5432/db" up-one
gograte migrate postgres "postgres://localhost:5432/db" up-to 20240228042140
gograte migrate postgres "postgres://localhost:5432/db" down-one
gograte migrate postgres "postgres://localhost:5432/db" down-to 20240228042140`

var initUsage = `No migration directory found.
Please run "gograte init" to create a migrations directory.`

var initContent = `
-- +gograte Up
-- SQL in section 'Up' is executed when this migration is applied

-- +gograte Down
-- SQL section 'Down' is executed when this migration is rolled back
`
