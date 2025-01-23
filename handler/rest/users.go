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

func DeleteAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Account not found!", "find by id")
	}

	usersCollection := database.ConnectCollection("users")
	filter := bson.M{"_id": userId}

	var userAccount bson.M
	err = usersCollection.FindOne(context.TODO(), filter).Decode(&userAccount)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Account not found!", "find in collection")
		}
		return Conflict(c, "Failed to get account details! Please try again...", "conflict find in colecction")
	}

	if deletedAt, ok := userAccount["deleted_at"]; ok && deletedAt != nil {
		return AlreadyDeleted(c, "This account already deleted!", "checking deleted account", deletedAt)
	}

	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	result, err := usersCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Conflict(c, "Failed to delete this account! Try again...", "update deleted_at line in DB")
	}

	if result.MatchedCount == 0 {
		return NotFound(c, "Can not find this account!", "failed find account data")
	}

	return OK(c, "Account successfully deleted!", userId)

}
