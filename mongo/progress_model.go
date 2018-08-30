package mongo

import (
	"github.com/globalsign/mgo/bson"
	"github.com/tobscore/aity/model"
)

type progressModel struct {
	Id       bson.ObjectId `bson:"id,omitempty"`
	Username string
	Date     string
	Distance model.Distance
}

func newProgressModel(username string, p *model.Progress) *progressModel {
	return &progressModel{
		Username: username,
		Date:     p.Date,
		Distance: p.Distance,
	}
}

func (p *progressModel) toProgress() *progressModel {
	return &progressModel{
		Date:     p.Date,
		Distance: p.Distance,
	}
}
