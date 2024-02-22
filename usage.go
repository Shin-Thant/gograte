package gograte

var UsageMessage = `Usage:
// Init
gograte init

// Create
gograte create [migration_name]

// Migrate
gograte migrate [db_driver] [db_url] [migrate_action]`

var initUsage = `No migration directory found.
Please run "gograte init" to create a migrations directory.`

var initContent = `
-- +gograte Up
-- SQL in section 'Up' is executed when this migration is applied

-- +gograte Down
-- SQL section 'Down' is executed when this migration is rolled back
`
