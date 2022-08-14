package main

import (
	"log"

	"gorm/connection"
	"gorm/model"

	"github.com/gofiber/fiber/v2"
)

var StringConnection = "root:12345678@tcp(127.0.0.1:3306)/gorm_db"
var DbConnection = connection.GetConnection(StringConnection)

func main() {

	app := fiber.New()

	/*
		User routes
	*/
	app.Get("/api/users", GetAllUser)
	app.Get("/api/users/:id", GetOneByIdUser)
	app.Post("/api/users", InsertUser)
	app.Delete("/api/users/:id", DeleteByIdUser)
	app.Put("/api/users/:id", UpdateUser)

	/*
		Product routes
	*/
	app.Get("/api/products", GetAllProducts)

	log.Fatal(app.Listen(":8080"))

}

//User implementation

// InsertUser - Insert a new user
func InsertUser(c *fiber.Ctx) error {
	user := new(model.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad request")
	}

	DbConnection.Create(&user)
	return c.JSON(user)
}

func GetAllUser(c *fiber.Ctx) error {
	var users []model.User

	DbConnection.Preload("Products").Find(&users)

	return c.JSON(users)
}

func GetOneByIdUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := model.User{}

	result := DbConnection.Find(&user, id)

	if result.RowsAffected > 0 {
		return c.JSON(user)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"code":    fiber.StatusNotFound,
		"message": "Not found",
	})
}

func DeleteByIdUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user := model.User{}

	DbConnection.Find(&user, id)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "User does not exist in the database",
		})
	}

	DbConnection.Delete(&user)
	return c.SendString("User deleted")
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := model.User{}

	DbConnection.Find(&user, id)
	if user.ID > 0 {
		newUser := model.User{}
		if err := c.BodyParser(&newUser); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": err.Error(),
			})
		}
		user.Name = newUser.Name
		user.Email = newUser.Email
		user.Password = newUser.Password

		DbConnection.Save(&user)
		return c.SendString("User updated")
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"code":    fiber.StatusNotFound,
		"message": "User does not exist in the database",
	})
}

// Product implementation

func GetAllProducts(c *fiber.Ctx) error {
	var products []model.Product

	DbConnection.Find(&products)

	return c.JSON(products)
}
