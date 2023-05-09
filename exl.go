package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/tealeg/xlsx"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"Name"`
	Email string `json:"email"`
	City  string `json:"city"`
	State string `json:"state"`
}

func main() {
	// Connect to PostgreSQL database
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=123456 dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize Gin router
	router := gin.Default()

	// Define route handlers
	router.GET("/excl", func(c *gin.Context) {
		// Query the database for all users
		rows, err := db.Query("SELECT * FROM excl")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query database"})
			return
		}
		defer rows.Close()

		// Create slice to store users
		var users []User

		// Iterate through rows and populate users slice
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.City, &user.State)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row"})
				return
			}
			users = append(users, user)
		}

		// Check for errors after iterating through rows
		err = rows.Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to iterate rows"})
			return
		}

		// Create new Excel file and add sheet
		file := xlsx.NewFile()
		sheet, err := file.AddSheet("excl")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sheet"})
			return
		}

		// Add header row to sheet
		row := sheet.AddRow()
		row.AddCell().SetValue("ID")
		row.AddCell().SetValue("Name")
		row.AddCell().SetValue("Email")
		row.AddCell().SetValue("City")
		row.AddCell().SetValue("State")

		// Add user data to sheet
		for _, user := range users {
			row := sheet.AddRow()
			row.AddCell().SetValue(fmt.Sprintf("%d", user.ID))
			row.AddCell().SetValue(user.Name)
			row.AddCell().SetValue(user.Email)
			row.AddCell().SetValue(user.City)
			row.AddCell().SetValue(user.State)
		}

		// Save Excel file to disk
		err = file.Save("excl.xlsx")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Return success message to client
		c.JSON(http.StatusOK, gin.H{"message": "Excel file created successfully"})
	})

	router.Run(":8080")
}
