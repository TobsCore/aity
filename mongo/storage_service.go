package mongo

import (
	"errors"
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
	activeTrack *mgo.Collection
}

const (
	trackCollectionName        = "tracks"
	progressCollectionName     = "progress"
	userCollectionName         = "users"
	currentTrackCollectionName = "active_tracks"
)

func NewStorageService(session *Session, dbName string) *StorageService {
	trackCol := session.GetCollection(dbName, trackCollectionName)
	trackCol.EnsureIndex(trackIndex())

	progressCol := session.GetCollection(dbName, progressCollectionName)

	userCol := session.GetCollection(dbName, userCollectionName)
	userCol.EnsureIndex(userIndex())

	activeTrackCol := session.GetCollection(dbName, currentTrackCollectionName)
	activeTrackCol.EnsureIndex(currentTrackIndex())

	return &StorageService{trackCol: trackCol, progressCol: progressCol, userCol: userCol, activeTrack: activeTrackCol}
}

func toLower(username string) string {
	return strings.ToLower(username)
}

func (s *StorageService) CreateTrack(username string, t *model.Track) (string, error) {
	username = toLower(username)
	track := newTrackModel(username, t)
	_, err := s.trackCol.UpsertId(track.Id, &track)
	return track.Id.Hex(), err
}

func (s *StorageService) GetTrackById(id, username string) (*model.Track, error) {
	t := trackModel{}
	objId := bson.ObjectIdHex(id)
	err := s.trackCol.Find(bson.M{"_id": objId}).One(&t)
	if strings.Compare(t.Username, username) != 0 {
		return &model.Track{}, errors.New("user " + username + " does not own the track with id " + id)
	}
	track := t.toTrack()
	s.checkActive(track)
	return track, err
}

func (s *StorageService) GetCurrentTrackByUsername(username string) (*model.Track, error) {
	username = toLower(username)
	currTrack := currentTrackModel{}
	trackMod := trackModel{}

	// Find the current track for the given username
	err := s.activeTrack.Find(bson.M{"username": username}).One(&currTrack)
	if err != nil {
		log.Fatal("Cannot retrieve current track info for user " + username)
		return trackMod.toTrack(), err
	}

	// Retreive the track document from the tracks collection
	err = s.trackCol.Find(bson.M{"_id": currTrack.Track}).One(&trackMod)

	// Set the current track as active
	track := trackMod.toTrack()
	track.Active = true
	return track, err
}

func (s *StorageService) GetProgressByUsername(trackid, username string) ([]progressModel, error) {
	username = toLower(username)
	progress := make([]progressModel, 20)
	err := s.progressCol.Find(bson.M{"username": username, "track": bson.ObjectIdHex(trackid)}).All(&progress)
	return progress, err
}

func (s *StorageService) AddProgressToTrack(trackid, username string, p *model.Progress) error {
	username = toLower(username)
	progress, err := newProgressModel(trackid, username, p)
	if err != nil {
		return err
	}
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

func (s *StorageService) CreateUser(user *model.User) error {
	userMod := NewUserModel(user)
	err := s.userCol.Insert(userMod)
	return err
}
func (s *StorageService) UpdateCurrentTrack(user string, t model.Track) error {
	currentTrack := NewCurrentTrackModel(user, t.Id)
	_, err := s.activeTrack.Upsert(bson.M{"username": user}, &currentTrack)
	return err
}

// checkActive checks, if the given track is active and sets the active value accordingly.
func (s *StorageService) checkActive(track *model.Track) {
	c, err := s.activeTrack.Find(bson.M{"track_id": bson.ObjectIdHex(track.Id)}).Count()
	if err != nil {
		log.Fatal(err)
	}

	active := c == 1
	track.Active = active
}
