package gograte

func GetSQLDriver(driver string) string {
	switch driver {
	case "mssql":
		return "sqlserver"
	case "sqlite3":
		return "sqlite"
	case "postgres", "redshift":
		return "pgx"
	}
	return ""
}
