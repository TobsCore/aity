package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tobscore/aity/model"
	"time"
)

type userModel struct {
	Id             bson.ObjectId `bson:"id,omitempty"`
	Name           string
	Email          string
	RegisteredDate time.Time
}

func userIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func NewUserModel(u *model.User) *userModel {
	return &userModel{
		Name:           u.Name,
		Email:          u.Email,
		RegisteredDate: u.RegisteredDate,
	}
}
func (u *userModel) toUser() *model.User {
	return &model.User{
		Email:          u.Email,
		Name:           u.Name,
		RegisteredDate: u.RegisteredDate,
	}
}
