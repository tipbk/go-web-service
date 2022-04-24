package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/tipbk/go-web-service/model"
	"github.com/tipbk/go-web-service/openapi"
)

func main() {
	app := fiber.New()

	SetupRoutes(app)

	app.Listen(":8080")
}

var tempMemory []model.Person

func SetupRoutes(app *fiber.App) {
	app.Get("/", GetRootRoute)
	app.Post("/", PostRootRoute)
	app.Post("/person", PostCreatePerson)
	app.Get("/person", GetAllPeople)
	app.Post("/user", AddUser)
}

func GetRootRoute(c *fiber.Ctx) error {
	return c.SendString("Hello, World!")
}

func PostRootRoute(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(model.Person{Friends: nil, Firstname: "Chanathip", Lastname: "Nateprapai"})
}

func PostCreatePerson(c *fiber.Ctx) error {
	person := model.Person{}
	err := c.BodyParser(&person)
	var myerror []openapi.OpenApiError
	if err != nil {
		myerror = append(myerror, openapi.OpenApiError{ErrorCode: "ERR01", ErrorReason: "Request cannot be processed"})
		return c.Status(fiber.ErrBadRequest.Code).JSON(myerror)
	}
	if person.Firstname == "" || person.Lastname == "" {
		myerror = append(myerror, openapi.OpenApiError{ErrorCode: "ERR02", ErrorReason: "Missing Required Parameters"})
		return c.Status(fiber.ErrBadRequest.Code).JSON(myerror)
	}

	tempMemory = append(tempMemory, person)
	return c.JSON(person)
}

func GetAllPeople(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(tempMemory)
}

type Job struct {
	Type   string `validate:"required,min=3,max=32"`
	Salary int    `validate:"required"`
}

type User struct {
	Name string `validate:"required,min=3,max=32"`
	// use `*bool` here otherwise the validation will fail for `false` values
	// Ref: https://github.com/go-playground/validator/issues/319#issuecomment-339222389
	IsActive *bool  `validate:"required"`
	Email    string `validate:"required,email,min=6,max=32"`
	Job      Job    `validate:"dive"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

var validate = validator.New()

func ValidateStruct(user User) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func AddUser(c *fiber.Ctx) error {
	//Connect to database

	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := ValidateStruct(*user)
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)

	}

	//Do something else here

	//Return user
	return c.JSON(user)
}
