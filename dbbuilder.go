package main
import (
	"database/sql"
	"fmt"
	"os"
	"github.com/go-sql-driver/mysql"
)



func main() {

	cfg := mysql.Config {
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "food_data",
		AllowNativePasswords: true,
    }



	db,err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	var query string

	query = "CREATE TABLE IF NOT EXISTS Users(id INT AUTO_INCREMENT PRIMARY KEY, Name VARCHAR(255), Password VARCHAR(255), MaxCalorieGoal INT)"

	create,err := db.Exec(query)
	if err != nil {
		panic(err)
	}

	fmt.Println(create)
}