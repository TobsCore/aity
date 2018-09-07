package mongo

import (
	"fmt"
	"github.com/globalsign/mgo"
)

type Session struct {
	session *mgo.Session
}

type Conn struct {
	Host string
	Port int
	User string
	Pwd  string
}

func NewSession(c Conn) (*Session, error) {
	url := fmt.Sprintf("%s:%d", c.Host, c.Port)
	cred := mgo.Credential{
		Username:    c.User,
		Password:    c.Pwd,
	}

	// Connect to the mongo instance
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	err = session.Login(&cred)
	if err != nil {
		return nil, err
	}
	return &Session{session:session}, err
}

func (s *Session) Copy() *Session {
	return &Session{s.session.Copy()}
}

func (s *Session) GetCollection(db , col string) *mgo.Collection {
	return s.session.DB(db).C(col)
}

func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}