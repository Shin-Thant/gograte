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

// Help
gograte --help

// Actions
init - Create a new migration directory
create - Create a new migration file
status - Show the status of all migrations
migrate - Apply or rollback migrations

// Migrate Actions
up - Apply all available migrations
down - Rollback all migrations
up-one - Apply the next available migration
down-one - Rollback the last applied migration
up-to - Apply all migrations up to a specific version
down-to - Rollback all migrations down to a specific version

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
