package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

type currentTrackModel struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Username string
	Track    bson.ObjectId `bson:"track_id"`
}

func currentTrackIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func NewCurrentTrackModel(username string, trackId string) *currentTrackModel {
	return &currentTrackModel{
		Username: username,
		Track: bson.ObjectIdHex(trackId),
	}
}
