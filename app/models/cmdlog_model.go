package models

import (
	`context`
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	_ "os"
	"strings"
)
//goland:noinspection Annotator
/*------------------------- GLOBAL VARs ------------------------------------------------------*/
var dbConnection *sql.DB
var ctx context.Context
/*------------------------- STRUCT & INTERFACE ------------------------------------------------------*/
type CmdLogger struct {
	id, time, command, result string
}
// Connect
///*------------------------- FUNCTIONS ------------------------------------------------------*/

func Connect() *sql.DB {
	fmt.Println(122)
	// Capture connection props
	// Get a database handle
	// os.Getenv("DASHPASS"),
	cfg := mysql.Config{
		User:   "dashhub",
		Passwd: os.Getenv("DASHPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:6603",
		DBName: "dash_main",
	}
	a := cfg.FormatDSN()
	fmt.Sprintf(a)
	// env details.
	pass := os.Getenv("DBPASS")
	user := os.Getenv("DBUSER")
	var err error
	//db_connection, err = sql.Open("mysql", user + ":" + pass +"@/dash_main")
	dbConnection, err = sql.Open("mysql", strings.Join([]string{user, ":", pass, "@tcp(localhost:6033)/dash_main"}, ""))
	//dbConnection, err = sql.Open("mysql", strings.Join([]string{"dashhub:phprest@tcp(0.0.0.0:6033)/dash_main"}, ""))
	if err != nil {
		log.Fatal(err)
	}
	pingErr := dbConnection.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
	return dbConnection
}

func CmdCreate(
	command string,
	result string) (sql.Result, bool) {
	db := Connect()
	stmt, err := db.Prepare("INSERT INTO dash_cmd_log(time, command, result) VALUES(NOW(),?,?)")
	checkErr(err)
	res, err := stmt.Exec(command, result)
	if err != nil {
		panic(err.Error())
		return res, false
	} else {
		log.Println(res)
		return res, true
	}
}

func CmdRead(
	all bool,
	query string,
	section string) (*sql.Rows, bool) {
	// Connect to Database.
	db := Connect()
	// Structure to define, what commands should be ran.
	if all == true && section == "nil" {
		// SQL QUERY
		sqlQuery := strings.Join([]string{"SELECT * FROM dash_cmd_log"}, "")
		rows, err := db.Query(sqlQuery)
		// err handl
		checkErr(err)

		return rows, true
	} else {
		if section != "nil" && all != true {
			// Execute command for specific query.
			sqlString := strings.Join([]string{"SELECT * FROM dash_cmd_log WHERE ",section, "=?"}, "")
			//fmt.Println(sqlString
			// Query
			rows, err := db.Query(sqlString, query)
			checkErr(err)
			return rows, true
		}
		db.Close()
	}
	if section != "nil" && all == true {
		// Execute all in that section, like cmd or timestamp etc.
		sqlString := strings.Join([]string{"SELECT ", section, " FROM dash_cmd_log"}, "")
		// Query
		rows, err := db.Query(sqlString)
		checkErr(err)

		return rows, true
	} else {
		fmt.Println("line 130")
		var x2 *sql.Rows
		return x2, false
	}
}

func CmdUpdate(
	queryID string,
	fieldID string,
	update string) (*sql.Rows, bool) {
	var sqlQuery string
	// Connect to Database.
	db := Connect()
	// Switch statement to define what fieldID and create the query string.
	switch fieldID {
	case "command":
		sqlQuery = strings.Join([]string{"UPDATE dash_cmd_log SET command='", update, "' WHERE id=?"}, "")
		fmt.Println("3")
	case "result":
		sqlQuery = strings.Join([]string{"UPDATE dash_cmd_log SET result='", update, "' WHERE id=?"}, "")
		fmt.Println("4")
	}
	// If there are no errors execute the command.
	res, err := db.Query(sqlQuery, queryID)
	checkErr(err)
	if err := res.Err(); err != nil {
		log.Fatal(err)
		return res, false
	} else {
		return res, true
	}
}

func CmdRemove(queryID string) (sql.Result, bool) {
	// Create string and prepared statement.
	sqlQuery := strings.Join([]string{"DELETE FROM dash_cmd_log WHERE dash_cmd_log.id=?"}, "")
	db := Connect()
	stmt, err := db.Prepare(sqlQuery)
	if err != nil {
		panic(err.Error())
	}
	// Execute command.
	res, err := stmt.Exec(queryID)
	if err != nil {
		panic(err.Error())
		return res, false
	} else {
		log.Println(res)
		return res, true
	}
}

func Query(password string, SQLInput string) (sql.Result, bool) {
	// SQL.Result Var
	var result sql.Result

	// dbConn *sql.DB, userInput string,
	db := Connect()
	// Const for verification hash... it's a hard word I promise x)
	const verifyHash = "68182ec52d1c385f2f847ac214ebbb675da00ab318274904c08eda394b648a05"
	// Hashing
	hash := sha256.New()
	hash.Write([]byte(password))
	hashSum := hash.Sum(nil)
	// Saving hash
	userHash := fmt.Sprintf("%x", hashSum)
	// (redutent)fmt.Println(userHash)
	// Verify Hash.
	if verifyHash == userHash {
		fmt.Println("Hello HASHSUM!")
		// Execution and Err Handling
		result, err := db.Exec(SQLInput)
		if err != nil {
			panic(err.Error())
			return result, false
		}
		// SQL Command Executed.
		return result, true
	} else {
		//fmt.Println(string(hashSum))
		fmt.Println("Failed to verify hash.")
		return result, false
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
