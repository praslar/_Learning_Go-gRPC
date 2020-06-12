package v1

import "time"

type TodoItem struct {
	ID          string    `json:"to_do_id" bson:"to_do_id"`
	Description string    `json:"description" bson:"description"`
	Title       string    `json:"title" bson:"title"`
	Reminder    time.Time `json:"reminder" bson:"reminder"`
}
