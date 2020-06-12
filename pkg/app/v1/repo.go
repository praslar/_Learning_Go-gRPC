package v1

import (
	"context"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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

func (r *MongoDBRepository) Create(ctx context.Context, item *TodoItem) error {
	s := r.session.Clone()
	defer s.Close()
	if err := r.collection(s).Insert(item); err != nil {
		return err
	}
	return nil
}

func (r *MongoDBRepository) FindByTitle(ctx context.Context, title string) (*TodoItem, error) {
	selector := bson.M{"title": title}
	s := r.session.Clone()
	defer s.Close()
	var todo *TodoItem
	if err := r.collection(s).Find(selector).One(&todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (r *MongoDBRepository) collection(s *mgo.Session) *mgo.Collection {
	return s.DB("").C("todo")
}
