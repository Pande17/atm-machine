package rest

import (
	"atm-machine/database"
	"atm-machine/model"
	"context"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func RegisterAccount(c *fiber.Ctx) error {
	var accReq struct {
		IDNumber int64  `json:"id_number"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&accReq); err != nil {
		return BadRequest(c, "Form can not be empty!", "body parser")
	}

	if _, err := govalidator.ValidateStruct(&accReq); err != nil {
		return BadRequest(c, "Input not valid", "Register validator")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(accReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return Conflict(c, "Can not hash the password", "hashing password register")
	}

	usersCollection := database.ConnectCollection("users")

	var existingAccount model.UserAccount
	filter := bson.M{"id_number": accReq.IDNumber}

	err = usersCollection.FindOne(context.TODO(), filter).Decode(&existingAccount)
	if err == nil {
		return BadRequest(c, "This ID Number already used", "existing account")
	} else if err != mongo.ErrNoDocuments {
		return Conflict(c, "Server error! Try again later", "existing account")
	}

	users := model.UserAccount{
		IDNumber: accReq.IDNumber,
		Username: accReq.Username,
		Password: string(hash),
		// Balance: ,
		Base: model.Base{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			DeletedAt: nil,
		},
	}

	_, err = usersCollection.InsertOne(context.TODO(), &users)
	if err != nil {
		return Conflict(c, "Can not Register now! Try again later..", "insert data")
	}

	return OK(c, "Account created successfully!", users)
}
