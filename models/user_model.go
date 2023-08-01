package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id          primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Username    string             `json:"username,omitempty" validate:"required"`
	Expiry_date int64              `json:"expiry_date,omitempty" validate:"required"`
	Outputs     []string           `json:"outputs,omitempty" validate:"required"`
	Password    string             `json:"password,omitempty" validate:"required"`
}

type GetUserResponse struct {
	Username    string   `json:"username,omitempty" validate:"required"`
	Expiry_date int64    `json:"expiry_date,omitempty" validate:"required"`
	Outputs     []string `json:"outputs,omitempty" validate:"required"`
}

type EditUser struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username    string             `bson:"username,omitempty" json:"username,omitempty"`
	Expiry_date int64              `bson:"expiry_date,omitempty" json:"expiry_date,omitempty"`
	Outputs     []string           `bson:"outputs,omitempty" json:"outputs,omitempty"`
	Password    string             `bson:"password,omitempty" json:"password,omitempty"`
}
