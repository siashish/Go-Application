package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"siashish/application/configs"
	"siashish/application/models"
	"siashish/application/responses"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	var data models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	_ = userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&data)
	if data.Username == user.Username {
		return c.Status(http.StatusConflict).JSON(responses.UserResponse{Status: http.StatusConflict, Message: "error", Data: &fiber.Map{"data": "username already exist"}})
	}

	password := HashPassword(user.Password)
	user.Password = password

	newUser := models.User{
		Id:          primitive.NewObjectID(),
		Username:    user.Username,
		Expiry_date: user.Expiry_date,
		Outputs:     user.Outputs,
		Password:    user.Password,
	}

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	err = userCollection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": user}})
}

func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	username := c.Params("username")
	var user models.User

	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	res := models.GetUserResponse{
		Username:    user.Username,
		Expiry_date: user.Expiry_date,
		Outputs:     user.Outputs,
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": res}})
}

func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	username := c.Params("username")
	var user models.EditUser
	var data models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	if user.Username != "" {
		_ = userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&data)
		if data.Username == user.Username {
			return c.Status(http.StatusConflict).JSON(responses.UserResponse{Status: http.StatusConflict, Message: "error", Data: &fiber.Map{"data": "username already exist"}})
		}
		data.Username = user.Username
	} else if user.Password != "" {
		password := HashPassword(user.Password)
		data.Password = password
	} else if user.Outputs != nil {
		data.Outputs = user.Outputs
	} else if user.Expiry_date != 0 {
		data.Expiry_date = user.Expiry_date
	}

	update := bson.M{"username": data.Username, "expiry_date": data.Expiry_date, "outputs": data.Outputs, "password": data.Password}

	_, err = userCollection.UpdateOne(ctx, bson.M{"username": username}, bson.M{"$set": update})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "The user " + username + " has been successfully updated!"})
}

func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	username := c.Params("username")
	defer cancel()

	result, err := userCollection.DeleteOne(ctx, bson.M{"username": username})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified username not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

func GetAllUsers(c *fiber.Ctx) error {
	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var users []models.User
	//defer cancel()

	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	username := c.Query("username", "")
	expiryDate, _ := strconv.Atoi(c.Query("expiry_date", ""))
	outputs := c.Query("outputs", "")
	operator := c.Query("operator", "")
	maxConnections, _ := strconv.Atoi(c.Query("max_connections", ""))
	sortBy := c.Query("sortBy", "")
	order := c.Query("order", "asc")

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64((page - 1) * limit))

	sortField := sortBy
	if order == "desc" {
		sortField = "-" + sortField
	}

	findOptions.SetSort(bson.M{sortField: 1})

	filter := bson.M{}
	if username != "" {
		filter["username"] = username
	}
	if expiryDate != 0 {
		filter["expiry_date"] = bson.M{"$gte": expiryDate}
	}
	if outputs != "" {
		filter["outputs"] = outputs
	}
	if operator != "" && maxConnections != 0 {
		switch operator {
		case "eq":
			filter["max_connections"] = maxConnections
		case "neq":
			filter["max_connections"] = bson.M{"$ne": maxConnections}
		case "gt":
			filter["max_connections"] = bson.M{"$gt": maxConnections}
		case "gte":
			filter["max_connections"] = bson.M{"$gte": maxConnections}
		case "lt":
			filter["max_connections"] = bson.M{"$lt": maxConnections}
		case "lte":
			filter["max_connections"] = bson.M{"$lte": maxConnections}
		}
	}

	cursor, err := userCollection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	if err = cursor.All(context.Background(), &users); err != nil {
		return c.Status(500).SendString(err.Error())
	}

	//reading from the db in an optimal way
	// defer cursor.Close(ctx)
	// for cursor.Next(ctx) {
	// 	var singleUser models.User
	// 	if err = cursor.Decode(&singleUser); err != nil {
	// 		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	// 	}

	// 	users = append(users, singleUser)
	// }

	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("login or passowrd is incorrect")
		check = false
	}

	return check, msg
}
