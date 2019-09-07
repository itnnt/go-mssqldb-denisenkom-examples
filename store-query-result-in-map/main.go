package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
)

var (
	debug         = flag.Bool("debug", false, "enable debugging")
	password      = flag.String("password", "123$%^qwe", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "localhost", "the database server")
	user          = flag.String("user", "sa", "the database user")
)

func main() {
	flag.Parse()

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)
	queryStr := `
select 1 as id, 'abc' as username
union all 
select 2 as id, 'abcd' as username
`

	query(queryStr, connString)

}

func query(queryStr string, connString string) map[string]interface{} {
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Error while opening database connection:", err.Error())
	}
	defer conn.Close()

	rows, err := conn.Query(queryStr)
	if err != nil {
		log.Fatal("Query failed:", err.Error())
	}
	defer rows.Close()

	cols, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column,
		// and a second slice to contain pointers to each item in the columns slice.
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			log.Fatal(err)
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			/* val := columnPointers[i].(*interface{})
			m[colName] = *val */
			val := columns[i]
			m[colName] = val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		fmt.Println(m)
	}
	return nil
}
