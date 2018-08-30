package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tobscore/aity/model"
)

type StorageService struct {
	trackCol *mgo.Collection
	userCol  *mgo.Collection // TODO: Add user support to storage service
}

const (
	trackCollectionName = "tracks"
	userCollectionName  = "users"
)

func NewStorageService(session *Session, dbName string) *StorageService {
	trackCol := session.GetCollection(dbName, trackCollectionName)
	trackCol.EnsureIndex(defaultIndex())

	userCol := session.GetCollection(dbName, userCollectionName)
	userCol.EnsureIndex(defaultIndex())

	return &StorageService{trackCol: trackCol, userCol: userCol}
}

func defaultIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func (p *StorageService) CreateTrack(username string, t *model.Track) (*mgo.ChangeInfo, error) {
	track := newTrackModel(username, t)
	return p.trackCol.Upsert(bson.M{"username": username}, &track)
}

func (p *StorageService) GetTrackByUsername(username string) (*model.Track, error) {
	trackMod := trackModel{}
	err := p.trackCol.Find(bson.M{"username": username}).One(&trackMod)
	return trackMod.toTrack(), err
}

