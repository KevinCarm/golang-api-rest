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

	app.Get("/", GetAll)

	log.Fatal(app.Listen(":8080"))

}

func GetAll(c *fiber.Ctx) error {
	users := []model.User{}

	DbConnection.Find(&users)

	return c.JSON(users)
}
