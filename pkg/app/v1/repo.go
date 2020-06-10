package v1

import (
	"github.com/globalsign/mgo"
)

type (
	MongoDBRepository struct {
		session *mgo.Session
	}
)

func NewMongoRepository(session *mgo.Session) *MongoDBRepository {
	return &MongoDBRepository{
		session: session,
	}
}
