package apps

import (
	"database/sql"
	"fmt"
	"time"

	db_client "github.com/RandySteven/go-cook/db"
	_ "github.com/go-sql-driver/mysql"
)

type mysqlClient struct {
	db *sql.DB
}

func (m *mysqlClient) Close() {
	_ = m.db.Close()
}

func (m *mysqlClient) Ping() error {
	return m.db.Ping()
}

func (m *mysqlClient) Client() *sql.DB {
	return m.db
}

var _ db_client.DBClient = (*mysqlClient)(nil)

// NewMySQLClient opens a real MySQL connection that still satisfies go-cook's
// DBClient interface (Save/Update/MigrationWorker). go-cook's NewMYSQLClient
// currently builds a Postgres-style DSN and does not register the MySQL driver.
func NewMySQLClient(config *db_client.DBConfig) (db_client.DBClient, error) {
	if config.Db != "mysql" {
		return nil, fmt.Errorf("expected db=mysql, got %s", config.Db)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	maxIdle := config.MaxIdleConns
	if maxIdle == 0 {
		maxIdle = 10
	}
	maxOpen := config.MaxOpenConns
	if maxOpen == 0 {
		maxOpen = 8
	}
	lifeMinutes := config.ConnMaxLifeTime
	if lifeMinutes == 0 {
		lifeMinutes = 10
	}
	idleMinutes := config.ConnMaxIdleTime
	if idleMinutes == 0 {
		idleMinutes = 8
	}

	db.SetMaxIdleConns(maxIdle)
	db.SetMaxOpenConns(maxOpen)
	db.SetConnMaxLifetime(time.Duration(lifeMinutes) * time.Minute)
	db.SetConnMaxIdleTime(time.Duration(idleMinutes) * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return &mysqlClient{db: db}, nil
}
