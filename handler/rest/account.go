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
	"go.mongodb.org/mongo-driver/mongo/options"
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

	accountCollection := database.ConnectCollection("account")

	var existingAccount model.Account
	filter := bson.M{"id_number": accReq.IDNumber}

	err = accountCollection.FindOne(context.TODO(), filter).Decode(&existingAccount)
	if err == nil {
		return BadRequest(c, "This ID Number already used", "existing account")
	} else if err != mongo.ErrNoDocuments {
		return Conflict(c, "Server error! Try again later", "existing account")
	}

	account := model.Account{
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

	_, err = accountCollection.InsertOne(context.TODO(), &account)
	if err != nil {
		return Conflict(c, "Can not Register now! Try again later..", "insert data")
	}

	return OK(c, "Account created successfully!", account)
}

func GetAllAccount(c *fiber.Ctx) error {
	var results []bson.M
	accountCollection := database.ConnectCollection("account")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projection := bson.M{
		"_id":       1,
		"id_number": 1,
		"username":  1,
		"balance":   1,
		"createdAt": 1,
		"updatedAt": 1,
		"deletedAt": 1,
	}

	cursor, err := accountCollection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "can not find account!", err.Error())
		}
		return Conflict(c, "failed to get data! Try again later..", "failed find account")
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var account bson.M
		if err := cursor.Decode(&account); err != nil {
			return Conflict(c, "failed decoded data!", "decodeing data")
		}
		results = append(results, account)
	}
	if err := cursor.Err(); err != nil {
		return Conflict(c, "cursor error! try again", "failed in cursor")
	}

	return OK(c, "Success get all account!", results)
}

func DeleteAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return BadRequest(c, "Account not found!", "find by id")
	}

	accountCollection := database.ConnectCollection("account")
	filter := bson.M{"_id": userId}

	var account bson.M
	err = accountCollection.FindOne(context.TODO(), filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return NotFound(c, "Account not found!", "find in collection")
		}
		return Conflict(c, "Failed to get account details! Please try again...", "conflict find in colecction")
	}

	if deletedAt, ok := account["deletedAt"]; ok && deletedAt != nil {
		return AlreadyDeleted(c, "This account already deleted!", "checking deleted account", deletedAt)
	}

	update := bson.M{"$set": bson.M{"deletedAt": time.Now()}}
	result, err := accountCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return Conflict(c, "Failed to delete this account! Try again...", "update deleted_at line in DB")
	}

	if result.MatchedCount == 0 {
		return NotFound(c, "Can not find this account!", "failed find account data")
	}

	return OK(c, "Account successfully deleted!", userId)

}
