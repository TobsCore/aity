package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tobscore/aity/model"
)

type trackModel struct {
	Id       bson.ObjectId `bson:"_id,omitempty"`
	Username string
	Start    model.Location
	End      model.Location
	Distance model.Distance
}

func trackIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"_id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func newTrackModel(username string, t *model.Track) *trackModel {
	id := bson.NewObjectId()
	return &trackModel{
		Id:       id,
		Username: username,
		Start:    t.Start,
		End:      t.End,
		Distance: t.Distance,
	}
}

func (t *trackModel) toTrack() *model.Track {
	return &model.Track{
		Id:       t.Id.Hex(),
		Start:    t.Start,
		End:      t.End,
		Distance: t.Distance,
	}
}
