package main

import (
	"database/sql"
    "fmt"
    "log"
    "os"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
    // Start DB logic and Capture connection properties.
    cfg := mysql.Config{
        User:   os.Getenv("DBUSER"),
        Passwd: os.Getenv("DBPASS"),
        Net:    "tcp",
        Addr:   "127.0.0.1:3306",
        DBName: "food_data",
		AllowNativePasswords: true,
    }

    // Get a database handle.
    var err error
    db, err = sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
        log.Fatal(err)
    }

    pingErr := db.Ping()
    if pingErr != nil {
        log.Fatal(pingErr)
    }
    fmt.Println("Connected!")






	firstName, lastName, maxCalorieGoal := greetUser()
	curCalorieGoal := maxCalorieGoal
	fmt.Printf("Hello, %s %s, let's get you to your calorie goal of %d calories!\n", firstName, lastName, curCalorieGoal)

	for {
		food, calories, day := getUserFood()
		curCalorieGoal -= calories

		totalCalories, err := calculateTotalCaloriesForDay(day)
		if err != nil {
			log.Fatal("Error calculating total calories: ", err)
		}

		err = updateTotalDayCalories(day, totalCalories)
		if err != nil {
			log.Fatal("Error updating total_day_calorie table: ", err)
		}

		fmt.Printf("Total calories for %s: %d\n", day, totalCalories)

		if curCalorieGoal > 0 {
			fmt.Printf("The %s you ate left you with %d calories left to reach your goal\n", food, curCalorieGoal)
		} else {
			fmt.Printf("Congrats you have reached your goal of %d calories\n", maxCalorieGoal)
			fmt.Printf("The %s you ate put you over your goal by %d calories\n", food, curCalorieGoal * -1)
			break
		}
	}

}


func greetUser() (string, string, int) {
	var firstName, lastName string
	var calorieGoal int

	fmt.Println("Enter your first name: ")
	fmt.Scan(&firstName)

	fmt.Println("Enter your last name: ")
	fmt.Scan(&lastName)

	fmt.Println("Enter your calorie goal: ")
	fmt.Scan(&calorieGoal)

	return firstName, lastName, calorieGoal
}

func getUserFood() (string, int, string) {
	var day string
	var food string
	var calories int

	fmt.Println("Enter the entry date in this format (YYYY-MM-DD): ")
	fmt.Scan(&day)
	fmt.Println("Enter the food you ate: ")
	fmt.Scan(&food)
	fmt.Println("Enter the number of calories in the food: ")
	fmt.Scan(&calories)

	_, err := db.Exec("INSERT INTO calorie_stats (day, food_item, calorie_amount) VALUES (?, ?, ?)", day, food, calories)
	if err != nil {
		log.Fatal("Error inserting into calorie_stats table: ", err)
	}
	return food, calories, day
}

// Sums the total calories for a given day
func calculateTotalCaloriesForDay(day string) (int, error) {
	var totalCalories int
	err := db.QueryRow("SELECT SUM(calorie_amount) FROM calorie_stats WHERE day = ?", day).Scan(&totalCalories)
	if err != nil {
		return 0, err
	}
	return totalCalories, nil
}

// Inserts or updates the total_day_calorie table with the summed calories for the day
func updateTotalDayCalories(day string, totalCalories int) error {
	_, err := db.Exec(`
		INSERT INTO total_day_calorie (day, total_calorie_amount) 
		VALUES (?, ?) 
		ON DUPLICATE KEY UPDATE total_calorie_amount = ?`, day, totalCalories, totalCalories)
	return err
}