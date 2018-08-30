package mongo

import (
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

func newTrackModel(username string, t *model.Track) *trackModel {
	return &trackModel{
		Username: username,
		Start:    t.Start,
		End:      t.End,
		Distance: t.Distance,
	}
}

func (t *trackModel) toTrack() *model.Track {
	return &model.Track{
		Start:    t.Start,
		End:      t.End,
		Distance: t.Distance,
	}
}
