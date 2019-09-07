package main

import (
	"database/sql"
	"encoding/json"
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

	result := executeSQL(queryStr, connString)
	myString := string(result[:])
	fmt.Println(myString)

	result = executesqlV2(queryStr, connString)
	myString = string(result[:])
	fmt.Println(myString)
}

func executeSQL(queryStr string, connString string) []byte {
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

	var v struct {
		Data []interface{} // `json:"data"`
	}

	for rows.Next() {
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal(err)
		}
		v.Data = append(v.Data, values)
	}
	jsonMsg, err := json.Marshal(v)
	return jsonMsg
}

func executesqlV2(queryStr string, connString string) []byte {
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

	var v struct {
		Data []interface{} // `json:"data"`
	}

	for rows.Next() {
		values := make([]interface{}, count)
		valuePtrs := make([]interface{}, count)
		for i, _ := range columns {
			valuePtrs[i] = &values[i]
		}
		if err := rows.Scan(valuePtrs...); err != nil {
			log.Fatal(err)
		}

		//Created a map to handle the issue
		var m map[string]interface{}
		m = make(map[string]interface{})
		for i := range columns {
			m[columns[i]] = values[i]
		}
		v.Data = append(v.Data, m)
	}
	jsonMsg, err := json.Marshal(v)
	return jsonMsg
}
