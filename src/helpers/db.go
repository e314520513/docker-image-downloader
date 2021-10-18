
package helpers

import(
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
func DbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbName := "golang"
        dbPassword :="golang"
	protocal := "tcp(localhost:3306)"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPassword+"@"+protocal+"/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}