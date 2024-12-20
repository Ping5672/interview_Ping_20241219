package main

import (
    "log"
    "interview_Ping_20241219/internal/api"
    "interview_Ping_20241219/internal/database"
)

func main() {
    // Initialize database
    database.InitDB()

    // Create and setup server
    server := api.NewServer()

    // Start server
    log.Fatal(server.Router().Run(":8080"))
}