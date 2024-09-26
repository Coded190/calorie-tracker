package main

import (
	"database/sql"
    "fmt"
    "log"
    "os"
	"html/template"
	"net/http"
	"strconv"
	"github.com/joho/godotenv"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB
var curCalGoal int

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", login)
	http.HandleFunc("/sign_up", sign_up)
	http.HandleFunc("/login_user", login_user)
	http.HandleFunc("/sign_up_user", sign_up_user)
	http.HandleFunc("/results", results)
	http.HandleFunc("/home", home)
	http.HandleFunc("/update_max_cal_goal", update_max_cal_goal)


	http.ListenAndServe(os.Getenv("PORT"), nil)


	/*
	
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

	defer db.Close()

	*/
}


func add_user(username string, password string, maxCalorieGoal int) bool {
	cfg := mysql.Config{
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

	add, err := db.Query("INSERT INTO Users(Name, Password, MaxCalorieGoal) VALUES (?, ?, ?)", (username), (password), (maxCalorieGoal))
	if err != nil {
		panic(err)
	}

	fmt.Println(add)
	defer db.Close()
	return true
}

func check_user(username string, password string) bool {
	cfg := mysql.Config{
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

	var exists bool
	var query string
	query = fmt.Sprintf("SELECT EXISTS(SELECT Name FROM Users WHERE Name='%s' AND Password='%s')", (username),(password))
	row := db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	fmt.Println(row)
	defer db.Close()
	return exists
}


func login(w http.ResponseWriter, r *http.Request){
	var tmplate = template.Must(template.ParseFiles("templates/login.html"))
	tmplate.Execute(w, nil)
}


func home(w http.ResponseWriter, r *http.Request){
	var tmplate = template.Must(template.ParseFiles("templates/index.html"))
	tmplate.Execute(w, nil)
}


func results(w http.ResponseWriter, r *http.Request){

	r.ParseForm()
	var date = r.Form["date"]
	var food_item = r.Form["food_item"]
	var cal = r.Form["calories"]
	calories, err := strconv.Atoi(cal[0])
	if err != nil {
		panic(err)
	}

	if (calculate_total_calories(date[0], food_item[0], calories)) {
		var tmplate = template.Must(template.ParseFiles("templates/results.html"))
		tmplate.Execute(w, nil)
	} else {
		var tmplate = template.Must(template.ParseFiles("templates/error.html"))
		tmplate.Execute(w, nil)
	}
}


func sign_up(w http.ResponseWriter, r *http.Request){
	var tmplate = template.Must(template.ParseFiles("templates/sign_up.html"))
	tmplate.Execute(w, nil)
}

func sign_up_user(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	var username = r.Form["fname"]
	var password = r.Form["Password"]
	var maxCalGoal = r.Form["maxCalGoal"]
	maxCalorieGoal, err := strconv.Atoi(maxCalGoal[0])
	if err != nil {
		panic(err)
	}

	fmt.Println(username, " ", password, " ", maxCalorieGoal)

	if (add_user(username[0], password[0], maxCalorieGoal)) {
		var tmplate = template.Must(template.ParseFiles("templates/index.html"))
		tmplate.Execute(w, nil)
	} else {
		var tmplate = template.Must(template.ParseFiles("templates/error.html"))
		tmplate.Execute(w, nil)
	}

	
}

func login_user(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	var username = r.Form["fname"]
	var password = r.Form["Password"]
	fmt.Println(username, " ", password)

	if (check_user(username[0], password[0])) {
		var tmplate = template.Must(template.ParseFiles("templates/index.html"))
		get_curCalGoal(username[0])
		tmplate.Execute(w, nil)
	} else {
		var tmplate = template.Must(template.ParseFiles("templates/error.html"))
		tmplate.Execute(w, nil)
	}
}

func get_curCalGoal(name string) {
	cfg := mysql.Config{
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

	// Updates curCalGoal with max calorie for user
	err = db.QueryRow("SELECT MaxCalorieGoal FROM Users WHERE Name = ?", name).Scan(&curCalGoal)
	if err != nil {
		panic(err)
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

// Updates the max calorie goal of the user
func update_max_cal_goal(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var fname = r.Form["fname"]
	var maxCalGoal = r.Form["maxCalGoal"]
	maxCalorieGoal, err := strconv.Atoi(maxCalGoal[0])
	if err != nil {
		panic(err)
	}

	if (update_max_calorie_goal(fname[0], maxCalorieGoal)) {
		var tmplate = template.Must(template.ParseFiles("templates/index.html"))
		tmplate.Execute(w, nil)
	} else {
		var tmplate = template.Must(template.ParseFiles("templates/error.html"))
		tmplate.Execute(w, nil)
	}
}

func update_max_calorie_goal(fname string, maxCalGoal int) bool {
	cfg := mysql.Config{
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

	_, err = db.Exec(`UPDATE Users SET MaxCalorieGoal = ? WHERE Name = ?`, maxCalGoal, fname)
	if err != nil {
		return false
	}

	defer db.Close()
	return true

}

func calculate_total_calories(date string, food_item string, calories int) bool {
	cfg := mysql.Config{
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

	// Update calorie_stats table
	_, err = db.Exec("INSERT INTO calorie_stats (day, food_item, calorie_amount) VALUES (?, ?, ?)", date, food_item, calories)
	if err != nil {
		return false
	}

	// Grabs the total calories from a specific day
	var total_calories int
	err = db.QueryRow("SELECT SUM(calorie_amount) FROM calorie_stats WHERE day = ?", date).Scan(&total_calories)
	if err != nil {
		return false
	}
	
	// Updates total_day_calorie for specific day
	_, err = db.Exec("INSERT INTO total_day_calorie (day, total_calorie_amount) VALUES (?, ?) ON DUPLICATE KEY UPDATE total_calorie_amount = ?", date, total_calories, total_calories)
	if err != nil {
		return false
	}

	defer db.Close()

	return true

}


