package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {

	/*stringConnection := "root:12345678@tcp(127.0.0.1:3306)/gorm_db"
	dbConnection := connection.GetConnection(stringConnection)*/

	app := fiber.New()

	/*user := model.User{
		Name:     "Miguel",
		Email:    "miguel@correo.com",
		Password: "12345678",
	}

	dbConnection.Create(&user)

	users := []model.User{}

	dbConnection.Find(&users)

	fmt.Println(users)*/

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello")
	})

	log.Fatal(app.Listen(":8080"))

}
