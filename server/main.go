package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// Database connection string (replace with your details)
const dbURL = "postgres://postgres:admin123@localhost:5437/try_db"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from Go backend!")
}

func main() {

	// Create a connection pool
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()
	fmt.Println("✅ Connected to PostgreSQL successfully!")

	// Test: Get all accounts
	GetAllAccounts(dbpool)

	InsertAccount(dbpool, "player3", "player3@example.com")

	GetAllCharacters(dbpool)

	GetAllScores(dbpool)

	GetRankings(dbpool)

	// Generate and insert fake data
	GenerateFakeAccounts(dbpool, 100000)
	GenerateFakeCharacters(dbpool, 100000)
	GenerateFakeScores(dbpool, 100000)

	fmt.Println("🚀 Data generation complete!")

	godotenv.Load()
	http.HandleFunc("/", handler)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

// Function to get all accounts from the database
func GetAllAccounts(dbpool *pgxpool.Pool) {
	rows, err := dbpool.Query(context.Background(), "SELECT acc_id, username, email FROM Account")
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	fmt.Println("📝 List of Accounts:")
	for rows.Next() {
		var accID int
		var username, email string
		err := rows.Scan(&accID, &username, &email)
		if err != nil {
			log.Fatalf("Row scan failed: %v\n", err)
		}
		fmt.Printf("🔹 ID: %d | Username: %s | Email: %s\n", accID, username, email)
	}
}
func InsertAccount(dbpool *pgxpool.Pool, username, email string) {
	// Check if the username already exists
	var exists bool
	err := dbpool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM Account WHERE username=$1)", username).Scan(&exists)

	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}

	if exists {
		fmt.Println("⚠️ Username already exists! Skipping insertion.")
		return
	}

	// Insert the new account
	_, err = dbpool.Exec(context.Background(),
		"INSERT INTO Account (username, email) VALUES ($1, $2)", username, email)
	if err != nil {
		log.Fatalf("Insert failed: %v\n", err)
	}
	fmt.Println("✅ New account added successfully!")
}
func GetAllCharacters(dbpool *pgxpool.Pool) {
	rows, err := dbpool.Query(context.Background(), "SELECT char_id, acc_id, class_id FROM Character")
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	fmt.Println("🎮 List of Characters:")
	for rows.Next() {
		var charID, accID, classID int
		err := rows.Scan(&charID, &accID, &classID)
		if err != nil {
			log.Fatalf("Row scan failed: %v\n", err)
		}
		fmt.Printf("🔹 Character ID: %d | Account ID: %d | Class ID: %d\n", charID, accID, classID)
	}
}
func GetAllScores(dbpool *pgxpool.Pool) {
	rows, err := dbpool.Query(context.Background(), "SELECT score_id, char_id, reward_score FROM Scores")
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	fmt.Println("🏆 List of Scores:")
	for rows.Next() {
		var scoreID, charID, rewardScore int
		err := rows.Scan(&scoreID, &charID, &rewardScore)
		if err != nil {
			log.Fatalf("Row scan failed: %v\n", err)
		}
		fmt.Printf("🔹 Score ID: %d | Character ID: %d | Score: %d\n", scoreID, charID, rewardScore)
	}
}
func InsertCharacter(dbpool *pgxpool.Pool, accID, classID int) {
	// Check if the account already has 8 classes
	var count int
	err := dbpool.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM Character WHERE acc_id=$1", accID).Scan(&count)

	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}

	if count >= 8 {
		fmt.Println("⚠️ Cannot add more classes! An account can have only 8 classes.")
		return
	}

	// Insert the new character
	_, err = dbpool.Exec(context.Background(),
		"INSERT INTO Character (acc_id, class_id) VALUES ($1, $2) ON CONFLICT (acc_id, class_id) DO NOTHING",
		accID, classID)

	if err != nil {
		log.Fatalf("Insert failed: %v\n", err)
	} else {
		fmt.Println("✅ New character added successfully!")
	}
}
func GetScoresForAccount(dbpool *pgxpool.Pool, username string) {
	rows, err := dbpool.Query(context.Background(),
		`SELECT c.class_id, s.reward_score 
		FROM Scores s 
		JOIN Character c ON s.char_id = c.char_id 
		JOIN Account a ON c.acc_id = a.acc_id 
		WHERE a.username = $1`, username)

	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	fmt.Printf("🏆 Scores for %s:\n", username)
	for rows.Next() {
		var classID, rewardScore int
		err := rows.Scan(&classID, &rewardScore)
		if err != nil {
			log.Fatalf("Row scan failed: %v\n", err)
		}
		fmt.Printf("🔹 Class ID: %d | Score: %d\n", classID, rewardScore)
	}
}
func GetRankings(dbpool *pgxpool.Pool) {
	rows, err := dbpool.Query(context.Background(),
		`SELECT a.username, c.class_id, s.reward_score,
		        RANK() OVER (PARTITION BY c.class_id ORDER BY s.reward_score DESC) AS rank
		 FROM Scores s
		 JOIN Character c ON s.char_id = c.char_id
		 JOIN Account a ON c.acc_id = a.acc_id
		 ORDER BY c.class_id, rank;`)

	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
	}
	defer rows.Close()

	fmt.Println("🏆 WIRA RANKING DASHBOARD 🏆")
	for rows.Next() {
		var username string
		var classID, score, rank int
		err := rows.Scan(&username, &classID, &score, &rank)
		if err != nil {
			log.Fatalf("Row scan failed: %v\n", err)
		}
		fmt.Printf("🔹 Rank: %d | Class ID: %d | Player: %s | Score: %d\n", rank, classID, username, score)
	}
}
func GenerateFakeAccounts(dbpool *pgxpool.Pool, count int) {
	fmt.Println("📝 Generating fake accounts...")

	for i := 0; i < count; i++ {
		username := faker.Username()
		email := faker.Email()
		_, err := dbpool.Exec(context.Background(),
			"INSERT INTO Account (username, email) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING", username, email)
		if err != nil {
			log.Fatalf("❌ Insert failed: %v\n", err)
		}
	}
	fmt.Println("✅ Fake accounts inserted successfully!")
}

// ✅ Generate Fake Characters
func GenerateFakeCharacters(dbpool *pgxpool.Pool, count int) {
	fmt.Println("🎮 Generating fake characters...")

	// Fetch all existing account IDs from the database
	var accIDs []int
	rows, err := dbpool.Query(context.Background(), "SELECT acc_id FROM Account")
	if err != nil {
		log.Fatalf("❌ Failed to get account IDs: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var accID int
		if err := rows.Scan(&accID); err != nil {
			log.Fatalf("❌ Failed to scan account ID: %v\n", err)
		}
		accIDs = append(accIDs, accID)
	}

	// Check if there are valid accounts
	if len(accIDs) == 0 {
		log.Fatalf("❌ No valid accounts found! Cannot insert characters.")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// Generate characters with existing account IDs
	for i := 0; i < count; i++ {
		accID := accIDs[rand.Intn(len(accIDs))] // Pick a random valid acc_id
		classID := rand.Intn(8) + 1             // Class ID between 1 and 8

		_, err := dbpool.Exec(context.Background(),
			"INSERT INTO Character (acc_id, class_id) VALUES ($1, $2) ON CONFLICT (acc_id, class_id) DO NOTHING",
			accID, classID)

		if err != nil {
			log.Fatalf("❌ Insert failed: %v\n", err)
		}
	}

	fmt.Println("✅ Fake characters inserted successfully!")
}

// ✅ Generate Fake Scores
func GenerateFakeScores(dbpool *pgxpool.Pool, count int) {
	fmt.Println("🏆 Generating fake scores...")

	// Fetch all existing character IDs and their class IDs
	type CharacterInfo struct {
		CharID  int
		ClassID int
	}
	var characters []CharacterInfo

	rows, err := dbpool.Query(context.Background(), "SELECT char_id, class_id FROM Character")
	if err != nil {
		log.Fatalf("❌ Failed to get character data: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var charInfo CharacterInfo
		if err := rows.Scan(&charInfo.CharID, &charInfo.ClassID); err != nil {
			log.Fatalf("❌ Failed to scan character data: %v\n", err)
		}
		characters = append(characters, charInfo)
	}

	// Check if there are valid characters
	if len(characters) == 0 {
		log.Fatalf("❌ No valid characters found! Cannot insert scores.")
	}

	// Seed random generator
	rand.Seed(time.Now().UnixNano())

	// Generate scores with existing character and class IDs
	for i := 0; i < count; i++ {
		randomChar := characters[rand.Intn(len(characters))] // Pick a random character
		rewardScore := rand.Intn(1000)                       // Random score (0-999)

		// ✅ Use `ON CONFLICT DO NOTHING` to avoid duplicate (`char_id`, `class_id`) errors
		_, err := dbpool.Exec(context.Background(),
			`INSERT INTO Scores (char_id, class_id, reward_score) 
			 VALUES ($1, $2, $3) 
			 ON CONFLICT (char_id, class_id) DO NOTHING`,
			randomChar.CharID, randomChar.ClassID, rewardScore)

		if err != nil {
			log.Fatalf("❌ Insert failed: %v\n", err)
		}
	}

	fmt.Println("✅ Fake scores inserted successfully!")
}
