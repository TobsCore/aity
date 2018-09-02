package mongo

import (
	"github.com/globalsign/mgo/bson"
	"github.com/tobscore/aity/model"
	"time"
)

type progressModel struct {
	Id       bson.ObjectId `bson:"id,omitempty"`
	Username string
	Date     time.Time
	Distance model.Distance
}

func newProgressModel(username string, p *model.Progress) (*progressModel, error) {
	date, err := p.GetDate()
	if err != nil {
		return &progressModel{}, err
	}
	return &progressModel{
		Username: username,
		Date:     date,
		Distance: p.Distance,
	}, nil
}

func (p *progressModel) ToProgress() *model.Progress {
	date := p.Date.Format(model.TimeFormat)
	return &model.Progress{
		Date:     date,
		Distance: p.Distance,
	}
}
