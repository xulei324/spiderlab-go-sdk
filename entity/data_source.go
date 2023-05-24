package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DataSource struct {
	Id         primitive.ObjectID `json:"_id" bson:"_id"`
	Type       string             `json:"type" bson:"type"`
	Host       string             `json:"host" bson:"host"`
	Port       string             `json:"port" bson:"port"`
	Database   string             `json:"database" bson:"database"`
	Username   string             `json:"username" bson:"username"`
	Password   string             `json:"password" bson:"password"`
	AuthSource string             `json:"auth_source" bson:"auth_source"`
	UserId     primitive.ObjectID `json:"user_id" bson:"user_id"`
	CreateTs   time.Time          `json:"create_ts" bson:"create_ts"`
	UpdateTs   time.Time          `json:"update_ts" bson:"update_ts"`
}
