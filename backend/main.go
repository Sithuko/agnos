package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	router := gin.Default()

	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	router.POST("/api/strong_password_steps", getStrongPasswordSteps)
	router.Run(":8080")
}

type PasswordRequest struct {
	InitPassword string `json:"init_password"`
}

type PasswordResponse struct {
	NumOfSteps int `json:"num_of_steps"`
}

func getStrongPasswordSteps(c *gin.Context) {
	var request PasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	steps := calculateSteps(request.InitPassword)

	response := PasswordResponse{NumOfSteps: steps}

	// Log request and response to the database
	logRequestResponse(request.InitPassword, steps)

	c.JSON(http.StatusOK, response)
}

//	func calculateSteps(password string) int {
//	    // Implement password strength logic here
//	    return 0 // Placeholder
//	}
func calculateSteps(password string) int {
	n := len(password)
	hasLower := false
	hasUpper := false
	hasDigit := false

	// Check character types
	for i := 0; i < n; i++ {
		if 'a' <= password[i] && password[i] <= 'z' {
			hasLower = true
		} else if 'A' <= password[i] && password[i] <= 'Z' {
			hasUpper = true
		} else if '0' <= password[i] && password[i] <= '9' {
			hasDigit = true
		}
	}

	missingTypes := 0
	if !hasLower {
		missingTypes++
	}
	if !hasUpper {
		missingTypes++
	}
	if !hasDigit {
		missingTypes++
	}

	replace := 0
	one := 0
	two := 0
	for i := 0; i < n; {
		j := i
		for i < n && password[i] == password[j] {
			i++
		}
		length := i - j
		if length >= 3 {
			replace += length / 3
			if length%3 == 0 {
				one++
			} else if length%3 == 1 {
				two++
			}
		}
	}

	if n < 6 {
		return max(missingTypes, 6-n)
	} else if n <= 20 {
		return max(missingTypes, replace)
	} else {
		delete := n - 20
		replace -= min(delete, one*1+two*2) / 3
		return delete + max(missingTypes, replace)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func logRequestResponse(password string, steps int) {
	sqlStatement := `
        INSERT INTO logs (password, steps)
        VALUES ($1, $2)`
	_, err := db.Exec(sqlStatement, password, steps)
	if err != nil {
		log.Fatalf("Unable to log request and response: %v", err)
	}
}
