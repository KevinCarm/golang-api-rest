package main

import (
	"log"

	"gorm/connection"
	"gorm/model"

	"github.com/gofiber/fiber/v2"
)

var StringConnection string = "root:12345678@tcp(127.0.0.1:3306)/gorm_db"
var DbConnection = connection.GetConnection(StringConnection)

func main() {

	app := fiber.New()

	app.Get("/api/users", GetAll)
	app.Get("/api/users/:id", GetOneById)
	app.Post("/api/users", Insert)
	app.Delete("/api/users/:id", DeleteById)
	app.Put("/api/users/:id", Update)

	log.Fatal(app.Listen(":8080"))

}

func Insert(c *fiber.Ctx) error {
	user := new(model.User)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad request")
	}

	DbConnection.Create(&user)
	return c.JSON(user)
}

func GetAll(c *fiber.Ctx) error {
	users := []model.User{}

	DbConnection.Preload("Products").Find(&users)

	return c.JSON(users)
}

func GetOneById(c *fiber.Ctx) error {
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

func DeleteById(c *fiber.Ctx) error {
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

func Update(c *fiber.Ctx) error {
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
