package main

import (
	"gorm/connection"
	"gorm/model"
	"log"

	"github.com/gofiber/fiber/v2"
)

var StringConnection = "root:12345678@tcp(127.0.0.1:3306)/gorm_db"
var DbConnection = connection.GetConnection(StringConnection)

func main() {

	app := fiber.New()

	/*
		User routes
	*/
	go app.Get("/api/users", GetAllUser)
	go app.Get("/api/users/:id", GetOneByIdUser)
	go app.Post("/api/users", InsertUser)
	go app.Delete("/api/users/:id", DeleteByIdUser)
	go app.Put("/api/users/:id", UpdateUser)

	/*
		Product routes
	*/
	go app.Get("/api/products", GetAllProducts)
	go app.Get("/api/products/:id", GetOneByIdProduct)
	go app.Post("/api/products/:id", InsertProduct)
	go app.Delete("/api/products", DeleteByIdProduct)
	go app.Put("/api/products/:id", UpdateProduct)

	log.Fatal(app.Listen(":8080"))

}

//User implementation

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

	return c.Status(fiber.StatusOK).JSON(products)
}

func GetOneByIdProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product := model.Product{}

	result := DbConnection.Find(&product, id)

	if result.RowsAffected > 0 {
		return c.Status(fiber.StatusOK).JSON(product)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"code":    fiber.StatusNotFound,
		"message": "Not found",
	})
}

func InsertProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	user := model.User{}
	DbConnection.Find(&user, id)

	if user.ID > 0 {
		product := model.Product{}
		if err := c.BodyParser(&product); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Bad request",
			})
		}
		DbConnection.Create(&product)
		return c.Status(fiber.StatusOK).JSON(product)
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"code":    fiber.StatusNotFound,
		"message": "User not found",
	})
}

func DeleteByIdProduct(c *fiber.Ctx) error {
	userId := c.Query("user")
	productId := c.Query("product")
	user := model.User{}

	DbConnection.Find(&user, userId)
	if user.ID > 0 {
		product := model.Product{}
		DbConnection.Find(&product, productId)
		if product.Id > 0 {
			DbConnection.Delete(&product, productId)
			return c.Status(fiber.StatusOK).SendString("Product deleted successfully")
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"code":    fiber.StatusNotFound,
			"message": "Product not found",
		})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"code":    fiber.StatusNotFound,
		"message": "User not found",
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	findProduct := model.Product{}
	DbConnection.Find(&findProduct, id)

	if findProduct.Id > 0 {
		newProduct := model.Product{}
		if err := c.BodyParser(&newProduct); err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    fiber.StatusBadRequest,
				"message": "Bad request",
				"err":     err.Error(),
			})
		}
		findProduct.Name = newProduct.Name
		findProduct.Price = newProduct.Price
		findProduct.UserId = newProduct.UserId
		DbConnection.Save(&findProduct)

		return c.Status(fiber.StatusOK).JSON("Product updated successfully")
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"code":    fiber.StatusNotFound,
		"message": "Product not found",
	})
}
