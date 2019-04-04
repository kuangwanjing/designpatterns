/*
introduce the first technique of customize json encoding and decoding:
we need to transform some fields of the struct into some other fields.
this refers to wonderful post: https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
*/
package json

import (
	"encoding/json"
	"fmt"
	"time"
)

// in this example, we need to transform field "born_at" from string into unix timestamp.
type Dog struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	BornAt Time   `json:"born_at"`
}

// here define a Time type as a bridge between two data types.
type Time struct {
	time.Time
}

// define MarshalJSON method
func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Time.Unix())
}

// and UnmarshalJSON
// note: the caller of this method should be the pointer of Time instance otherwise assignment is invalid.
func (t *Time) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)
	return nil
}

// clean codes to settle this problem: no field copying
func TestTech1() {
	dog := Dog{1, "bowser", "husky", Time{time.Now()}}
	b, err := json.Marshal(dog)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	b = []byte(`{
    "id":1,
    "name":"bowser",
    "breed":"husky",
    "born_at":1480979203}`)
	dog = Dog{}
	json.Unmarshal(b, &dog)
	fmt.Println(dog)
}
