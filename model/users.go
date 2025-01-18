package model

type UserAccount struct {
	IDNumber int64   `json:"id_number" bson:"id_number"`
	Username string  `json:"username" bson:"username"`
	Password string  `json:"password" bson:"password"`
	Balance  float32 `json:"balance" bson:"balance"`
}
