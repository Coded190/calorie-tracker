package main

import (
	"database/sql"
	//"context"
    "fmt"
    "log"
    "os"
	"html/template"
	"net/http"
	"strconv"
	//"github.com/joho/godotenv"

	//"github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB
var curCalGoal int

func main() {
    connString := "postgres://chris13:@localhost:5432/food_data"

    var err error
    db, err = sql.Open("pgx", connString)
    if err != nil {
        log.Fatal("Error connecting to database: ", err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal("Error pinging database: ", err)
    }
    fmt.Println("Connected to PostgreSQL")

    // File server for static assets
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Route handlers
    http.HandleFunc("/", login)
    http.HandleFunc("/sign_up", sign_up)
    http.HandleFunc("/login_user", login_user)
    http.HandleFunc("/sign_up_user", sign_up_user)
    http.HandleFunc("/results", results)
    http.HandleFunc("/home", home)
    http.HandleFunc("/update_max_cal_goal", update_max_cal_goal)

    // Starting the server
    port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
    log.Println("Starting server on port", port)
    err = http.ListenAndServe(":"+port, nil)
    if err != nil {
        log.Fatal("ListenAndServe failed:", err)
    }
	defer db.Close()
}


func add_user(username string, password string, maxCalorieGoal int) bool {
	add, err := db.Query("INSERT INTO Users(Name, Password, MaxCalorieGoal) VALUES ($1, $2, $3)", (username), (password), (maxCalorieGoal))
	if err != nil {
		panic(err)
	}

	fmt.Println(add)
	return true
}

func check_user(username string, password string) bool {
	rows, err := db.Query("SELECT username, password FROM users WHERE username = $1 AND password = $2", username, password)
	if err != nil {
		log.Fatal("Error querying users table: ", err)
	}

	for rows.Next() {
		var dbUsername string
		var dbPassword string

		err := rows.Scan(&dbUsername, &dbPassword)
		if err != nil {
			fmt.Println("User not found: ", err)
			return false
		}
	}
	defer rows.Close()
	return true
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
	var username = r.Form["username"]
	var password = r.Form["Password"]
	fmt.Println(username[0], " ", password[0])

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
	rows, err := db.Query("SELECT max_calorie_goal FROM users WHERE username = $1", name)
	if err != nil {
		log.Fatal("Error querying users table: ", err)
	}

	for rows.Next() {
		var dbUsername string

		err := rows.Scan(&dbUsername)
		if err != nil {
			log.Fatal("Error scanning rows: ", err)
		}
	}
	defer rows.Close()
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

	_, err := db.Exec("INSERT INTO calorie_stats (day, food_item, calorie_amount) VALUES ($1, $2, $3)", day, food, calories)
	if err != nil {
		log.Fatal("Error inserting into calorie_stats table: ", err)
	}
	return food, calories, day
}

// Sums the total calories for a given day
func calculateTotalCaloriesForDay(day string) (int, error) {
	var totalCalories int
	err := db.QueryRow("SELECT SUM(calorie_amount) FROM calorie_stats WHERE day = $1", day).Scan(&totalCalories)
	if err != nil {
		return 0, err
	}
	return totalCalories, nil
}

// Inserts or updates the total_day_calorie table with the summed calories for the day
func updateTotalDayCalories(day string, totalCalories int) error {
	_, err := db.Exec(`
		INSERT INTO total_day_calorie (day, total_calorie_amount) 
		VALUES ($1, $2) 
		ON DUPLICATE KEY UPDATE total_calorie_amount = $3`, day, totalCalories, totalCalories)
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
	_, err := db.Exec(`UPDATE Users SET MaxCalorieGoal = $1 WHERE Name = $2`, maxCalGoal, fname)
	if err != nil {
		return false
	}

	return true

}

func calculate_total_calories(date string, food_item string, calories int) bool {
	// Update calorie_stats table
	_, err := db.Exec("INSERT INTO calorie_stats (day, food_item, calorie_amount) VALUES ($1, $2, $3)", date, food_item, calories)
	if err != nil {
		return false
	}

	// Grabs the total calories from a specific day
	var total_calories int
	err = db.QueryRow("SELECT SUM(calorie_amount) FROM calorie_stats WHERE day = $1", date).Scan(&total_calories)
	if err != nil {
		return false
	}
	
	// Updates total_day_calorie for specific day
	_, err = db.Exec("INSERT INTO total_day_calorie (day, total_calorie_amount) VALUES ($1, $2) ON DUPLICATE KEY UPDATE total_calorie_amount = $3", date, total_calories, total_calories)
	if err != nil {
		return false
	}

	return true

}


