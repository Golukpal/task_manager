package main

import (
    "log"
    "os"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/joho/godotenv"
    "github.com/Golukpal/task_manager/internal/database"
    "github.com/Golukpal/task_manager/internal/handlers"
    "github.com/Golukpal/task_manager/internal/middleware"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    database.ConnectDB()

    app := fiber.New()

    app.Use(cors.New(cors.Config{
        AllowCredentials: true,
    }))

    // Public routes
    app.Post("/api/register", handlers.Register)
    app.Post("/api/login", handlers.Login)

    // Protected routes
    api := app.Group("/api", middleware.Protected())
    api.Post("/tasks", handlers.CreateTask)
    api.Get("/tasks", handlers.GetTasks)
    api.Get("/tasks/:id", handlers.GetTask)
    api.Put("/tasks/:id", handlers.UpdateTask)
    api.Delete("/tasks/:id", handlers.DeleteTask)

    port := os.Getenv("APP_PORT")
    log.Fatal(app.Listen(":" + port))
}
