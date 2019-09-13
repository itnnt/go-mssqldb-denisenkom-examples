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
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)
	queryStr := `
	select 1 as id, 'abc' as username
	union all 
	select 2 as id, 'abcd' as username
	`
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

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	var final_result map[int]map[string]interface{}
	final_result = make(map[int]map[string]interface{})

	result_id := 0
	for rows.Next() {
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		var tmpStructs map[string]interface{}
		tmpStructs = make(map[string]interface{})

		for i, col := range columns {
			var v interface{}
			v = values[i]
			tmpStructs[col] = v
		}

		final_result[result_id] = tmpStructs
		result_id++
	}

	//fmt.Println(final_result)
	tem := final_result[0]["username"]
	tem2 := final_result[0]["id"]
	fmt.Printf("%T, %v\n", tem, tem)
	fmt.Printf("%T, %v\n", tem2, tem2)

}
