package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
	"vdicalc/functions"

	"github.com/doug-martin/goqu"
	_ "github.com/doug-martin/goqu/dialect/mysql" // import the dialect
)

// DBInit export
/* This function initializes GCP mysql database connectivity */
func DBInit() *sql.DB {

	var db *sql.DB
	var err error

	// If the optional DB_TCP_HOST environment variable is set, it contains
	// the IP address and port number of a TCP connection pool to be created,
	// such as "127.0.0.1:3306". If DB_TCP_HOST is not set, a Unix socket
	// connection pool will be created instead.
	if os.Getenv("DB_TCP_HOST") != "" {

		db, err = InitTCPConnectionPool()
		if err != nil {
			log.Fatalf("initTCPConnectionPool: unable to connect: %v", err)
		}
	} else {
		db, err = InitSocketConnectionPool()
		if err != nil {
			log.Fatalf("initSocketConnectionPool: unable to connect: %v", err)
		}
	}

	return db

}

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

// QueryUser exported
/* This public function retrieve a single user from vdicalc.users */
func QueryUser(db *sql.DB, UserID string) bool {

	sqlSelect, _ := sqlBuilderSelectWhere("vdicalc.users", map[string]interface{}{
		"guserid": UserID,
	})

	var (
		id       int
		datetime string
		guserid  string
		email    string
	)

	rows, err := db.Query(sqlSelect)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &datetime, &guserid, &email)
		if err != nil {
			log.Fatal(err)
		}
		if guserid == UserID {
			return true
		}
	}

	return false

}

// CreateUser export
/* This function inserts new user into vdicalc.users */
func CreateUser(db *sql.DB, userid, email string) {

	/* Build MySQL statement  */
	sqlInsert, _ := SQLBuilderInsert("users", map[string]interface{}{
		"datetime": time.Now(),
		"guserid":  userid,
		"email":    email,
	})

	/* This function execues the SQL estatement on Google SQL Run database */
	Insert(db, sqlInsert)

}

// sqlBuilderSelectWhere export
// This functions uses goqu packages to create a mySQL compatible SQL statement
// github.com/doug-martin/goqu
func sqlBuilderSelectWhere(table string, s map[string]interface{}) (string, []interface{}) {

	dialect := goqu.Dialect("mysql")
	ds := dialect.From(table).Where(goqu.Ex(s))

	sql, args, err := ds.ToSQL()
	if err != nil {
		fmt.Println("An error occurred while generating the SQL", err.Error())
	}

	return sql, args
}

// SQLBuilderInsert export
// This functions uses goqu packages to create a mySQL compatible SQL
// statement and require input as map[string]interface{}
// github.com/doug-martin/goqu
func SQLBuilderInsert(table string, s ...interface{}) (string, []interface{}) {

	dialect := goqu.Dialect("mysql")
	ds := dialect.Insert(table).Rows(s)

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
