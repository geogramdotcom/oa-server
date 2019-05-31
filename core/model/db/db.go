package db

import (
	"database/sql"
	"fmt"
	"runtime/debug"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

type Datastore interface {
	Escape(string) string
	UserInterface
	OrgInterface
	AccountInterface
	TransactionInterface
	PriceInterface
	SessionInterface
	ApiKeyInterface
	SystemHealthInteface
}

func NewDB(dataSourceName string) (*DB, error) {
	fmt.Println("starting NewDB Func in db.go")

	var err error
	db, err := sql.Open("mysql", dataSourceName)

	if err != nil {
		debug.PrintStack()
		return nil, err
	}
	// fmt.Println("No errors; trying SetConnMaxLifetime")
	// db.SetConnMaxLifetime(time.Second * 5)
	// fmt.Println("No errors; trying SetMaxIdleConns")
	// db.SetMaxIdleConns(0)
	// fmt.Println("No errors; trying SetMaxOpenConns")
	// db.SetMaxOpenConns(151)

	if err = db.Ping(); err != nil {
		fmt.Println("try ping")
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) Escape(sql string) string {
	dest := make([]byte, 0, 2*len(sql))
	var escape byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]

		escape = 0

		switch c {
		case 0: /* Must be escaped for 'mysql' */
			escape = '0'
			break
		case '\n': /* Must be escaped for logs */
			escape = 'n'
			break
		case '\r':
			escape = 'r'
			break
		case '\\':
			escape = '\\'
			break
		case '\'':
			escape = '\''
			break
		case '"': /* Better safe than sorry */
			escape = '"'
			break
		case '\032': /* This gives problems on Win32 */
			escape = 'Z'
			break
		case '%':
			escape = '%'
		}

		if escape != 0 {
			dest = append(dest, '\\', escape)
		} else {
			dest = append(dest, c)
		}
	}

	return string(dest)
}
