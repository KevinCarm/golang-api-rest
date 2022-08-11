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
	app.Post("/api/users", Insert)

	log.Fatal(app.Listen(":8080"))

}

func Insert(c *fiber.Ctx) error {
	user := new(model.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad request")
	}

	DbConnection.Create(&user)
	return c.JSON(user)
}

func GetAll(c *fiber.Ctx) error {
	users := []model.User{}

	DbConnection.Find(&users)

	return c.JSON(users)
}
