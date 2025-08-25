package handlers

import (
    "github.com/Golukpal/task_manager/internal/database"
    "github.com/Golukpal/task_manager/internal/models"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "os"
    "time"
)

func Register(c *fiber.Ctx) error {
    var data models.RegisterInput

    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid input",
        })
    }

    password, _ := bcrypt.GenerateFromPassword([]byte(data.Password), 14)

    user := models.User{
        Username: data.Username,
        Password: string(password),
    }

    result := database.DB.Create(&user)

    if result.Error != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Could not create user",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "User created successfully",
        "user":    user,
    })
}

func Login(c *fiber.Ctx) error {
    var data models.LoginInput

    if err := c.BodyParser(&data); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid input",
        })
    }

    var user models.User
    database.DB.Where("username = ?", data.Username).First(&user)

    if user.ID == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "User not found",
        })
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password)); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Incorrect password",
        })
    }

    claims := jwt.MapClaims{
        "user_id":  user.ID,
        "username": user.Username,
        "exp":      time.Now().Add(time.Hour * 72).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Could not login",
        })
    }

    cookie := fiber.Cookie{
        Name:     "jwt",
        Value:    t,
        Expires:  time.Now().Add(time.Hour * 72),
        HTTPOnly: true,
    }

    c.Cookie(&cookie)

    return c.JSON(fiber.Map{
        "message": "Successfully logged in",
        "user":    user,
    })
}

func CreateTask(c *fiber.Ctx) error {
    var task models.Task

    if err := c.BodyParser(&task); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid input",
        })
    }

    claims := c.Locals("user").(jwt.MapClaims)
    userId := uint(claims["user_id"].(float64))
    task.UserID = userId

    result := database.DB.Create(&task)
    if result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Could not create task",
        })
    }

    return c.Status(fiber.StatusCreated).JSON(task)
}

func GetTasks(c *fiber.Ctx) error {
    claims := c.Locals("user").(jwt.MapClaims)
    userId := uint(claims["user_id"].(float64))

    var tasks []models.Task
    database.DB.Where("user_id = ?", userId).Find(&tasks)

    return c.JSON(tasks)
}

func GetTask(c *fiber.Ctx) error {
    id := c.Params("id")
    claims := c.Locals("user").(jwt.MapClaims)
    userId := uint(claims["user_id"].(float64))

    var task models.Task
    result := database.DB.Where("id = ? AND user_id = ?", id, userId).First(&task)

    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "Task not found",
        })
    }

    return c.JSON(task)
}

func UpdateTask(c *fiber.Ctx) error {
    id := c.Params("id")
    claims := c.Locals("user").(jwt.MapClaims)
    userId := uint(claims["user_id"].(float64))

    var task models.Task
    result := database.DB.Where("id = ? AND user_id = ?", id, userId).First(&task)

    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "Task not found",
        })
    }

    if err := c.BodyParser(&task); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Invalid input",
        })
    }

    database.DB.Save(&task)

    return c.JSON(task)
}

func DeleteTask(c *fiber.Ctx) error {
    id := c.Params("id")
    claims := c.Locals("user").(jwt.MapClaims)
    userId := uint(claims["user_id"].(float64))

    var task models.Task
    result := database.DB.Where("id = ? AND user_id = ?", id, userId).First(&task)

    if result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "message": "Task not found",
        })
    }

    database.DB.Delete(&task)

    return c.JSON(fiber.Map{
        "message": "Task deleted successfully",
    })
}
