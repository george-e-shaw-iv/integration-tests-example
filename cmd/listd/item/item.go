package item

import "time"

// Record is a type that contains the proper struct tags for both
// a JSON and Postgres representation of an item.
type Record struct {
	ID       int       `json:"id" db:"item_id"`
	ListID   int       `json:"listID" db:"list_id"`
	Name     string    `json:"name" db:"name"`
	Quantity int       `json:"quantity" db:"quantity"`
	Created  time.Time `json:"created" db:"created"`
	Modified time.Time `json:"modified" db:"modified"`
}
