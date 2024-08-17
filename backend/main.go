package main

import (
    "github.com/gin-gonic/gin"
    "log"
    "net/http"
    "os"
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
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

// func calculateSteps(password string) int {
//     // Implement password strength logic here
//     return 0 // Placeholder
// }
func calculateSteps(password string) int {
    n := len(password)
    hasLower := false
    hasUpper := false
    hasDigit := false

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

    changes := 0
    for i := 2; i < n; {
        if password[i] == password[i-1] && password[i] == password[i-2] {
            changes++
            i += 3
        } else {
            i++
        }
    }

    if n < 6 {
        return max(6-n, missingTypes)
    } else if n <= 20 {
        return max(changes, missingTypes)
    } else {
        return n - 20 + max(changes, missingTypes)
    }
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
