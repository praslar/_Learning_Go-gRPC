package db

import (
	"errors"
	"sync"

	"github.com/globalsign/mgo"
	"github.com/praslar/to-do-list-micro/pkg/database/mongodb"
)

var (
	session     *mgo.Session
	sessionOnce sync.Once
)

// IsErrNotFound return true if the given error is a not found error
func IsErrNotFound(err error) bool {
	return errors.Is(err, mgo.ErrNotFound)
}

// Connect to mongoDB
func DialDefaultMongoDB() (*mgo.Session, error) {
	repoConf := mongodb.LoadConfigFromEnv()
	var err error
	sessionOnce.Do(func() {
		session, err = mongodb.Dial(repoConf)
	})
	if err != nil {
		return nil, err
	}
	s := session.Clone()
	return s, nil
}
