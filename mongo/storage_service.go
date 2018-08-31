package mongo

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/tobscore/aity/model"
	"log"
	"strings"
)

type StorageService struct {
	trackCol    *mgo.Collection
	progressCol *mgo.Collection
	userCol     *mgo.Collection
}

const (
	trackCollectionName    = "tracks"
	progressCollectionName = "progress"
	userCollectionName     = "users"
)

func NewStorageService(session *Session, dbName string) *StorageService {
	trackCol := session.GetCollection(dbName, trackCollectionName)
	trackCol.EnsureIndex(trackIndex())

	progressCol := session.GetCollection(dbName, progressCollectionName)

	userCol := session.GetCollection(dbName, progressCollectionName)
	userCol.EnsureIndex(userIndex())

	return &StorageService{trackCol: trackCol, progressCol: progressCol, userCol: userCol}
}

func toLower(username string) string {
	return strings.ToLower(username)
}

func (s *StorageService) CreateTrack(username string, t *model.Track) (*mgo.ChangeInfo, error) {
	username = toLower(username)
	track := newTrackModel(username, t)
	return s.trackCol.Upsert(bson.M{"username": username}, &track)
}

func (s *StorageService) GetTrackByUsername(username string) (*model.Track, error) {
	username = toLower(username)
	trackMod := trackModel{}
	err := s.trackCol.Find(bson.M{"username": username}).One(&trackMod)
	return trackMod.toTrack(), err
}

func (s *StorageService) GetProgressByUsername(username string) ([]model.Progress, error) {
	username = toLower(username)
	progress := make([]model.Progress, 20)
	err := s.progressCol.Find(bson.M{"username": username}).All(&progress)
	return progress, err
}

func (s *StorageService) AddProgress(username string, p *model.Progress) error {
	username = toLower(username)
	progress := newProgressModel(username, p)
	return s.progressCol.Insert(&progress)
}

func (s *StorageService) UserExists(email string) bool {
	count, err := s.userCol.Find(bson.M{"email": email}).Count()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return count > 0
}

func (s *StorageService) GetUserByEmail(email string) (model.User, error) {
	user := model.User{}
	err := s.userCol.Find(bson.M{"email": email}).One(&user)
	return user, err
}

func (s *StorageService) CreateUser(user *model.User) (model.User, error) {
	err := s.userCol.Insert(user)
	return *user, err
}
