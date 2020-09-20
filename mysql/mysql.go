package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"vdicalc/functions"

	"github.com/doug-martin/goqu"
	_ "github.com/doug-martin/goqu/dialect/mysql" // import the dialect
)

// Insert exported
// This function executes a SQL Insert
func Insert(db *sql.DB, sqlInsert string) error {

	// [START cloud_sql_mysql_databasesql_connection]

	if _, err := db.Exec(sqlInsert); err != nil {
		return fmt.Errorf("DB.Exec: %v", err)
	}

	return nil
	// [END cloud_sql_mysql_databasesql_connection]
}

// SQLBuilderInsert export
// This functions uses goqu packages to create a mySQL compatible SQL
// statement and require input as map[string]interface{}
// github.com/doug-martin/goqu
func SQLBuilderInsert(s ...interface{}) (string, []interface{}) {

	dialect := goqu.Dialect("mysql")
	ds := dialect.Insert(functions.MustGetenv("DB_NAME")).Rows(s)

	sql, args, err := ds.ToSQL()
	if err != nil {
		fmt.Println("An error occurred while generating the SQL", err.Error())
	}

	return sql, args
}

// InitSocketConnectionPool initializes a Unix socket connection pool for
// a Cloud SQL instance of SQL Server.
func InitSocketConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_socket]
	var (
		dbUser                 = functions.MustGetenv("DB_USER")
		dbPwd                  = functions.MustGetenv("DB_PASS")
		instanceConnectionName = functions.MustGetenv("INSTANCE_CONNECTION_NAME")
		dbName                 = functions.MustGetenv("DB_NAME")
	)

	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
	if !isSet {
		socketDir = "/cloudsql"
	}

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@unix(/%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_socket]
}

// configureConnectionPool sets database connection pool properties.
// For more information, see https://golang.org/pkg/database/sql
func configureConnectionPool(dbPool *sql.DB) {
	// [START cloud_sql_mysql_databasesql_limit]

	// Set maximum number of connections in idle connection pool.
	dbPool.SetMaxIdleConns(5)

	// Set maximum number of open connections to the database.
	dbPool.SetMaxOpenConns(7)

	// [END cloud_sql_mysql_databasesql_limit]

	// [START cloud_sql_mysql_databasesql_lifetime]

	// Set Maximum time (in seconds) that a connection can remain open.
	dbPool.SetConnMaxLifetime(1800)

	// [END cloud_sql_mysql_databasesql_lifetime]
}

// InitTCPConnectionPool initializes a TCP connection pool for a Cloud SQL
// instance of SQL Server.
func InitTCPConnectionPool() (*sql.DB, error) {
	// [START cloud_sql_mysql_databasesql_create_tcp]
	var (
		dbUser    = functions.MustGetenv("DB_USER")
		dbPwd     = functions.MustGetenv("DB_PASS")
		dbTcpHost = functions.MustGetenv("DB_TCP_HOST")
		dbPort    = functions.MustGetenv("DB_PORT")
		dbName    = functions.MustGetenv("DB_NAME")
	)

	var dbURI string
	dbURI = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPwd, dbTcpHost, dbPort, dbName)

	// dbPool is the pool of database connections.
	dbPool, err := sql.Open("mysql", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	// [START_EXCLUDE]
	configureConnectionPool(dbPool)
	// [END_EXCLUDE]

	return dbPool, nil
	// [END cloud_sql_mysql_databasesql_create_tcp]
}
