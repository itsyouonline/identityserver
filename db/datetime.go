package db

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	dateTimeFmt       = "2006-01-02T15:04:05.999999999Z07:00"
	dateTimeFmtTicked = `"` + dateTimeFmt + `"`
)

// DateTime is timestamp in "date-time" format defined in RFC3339
type DateTime time.Time

// MarshalJSON override marshalJSON
func (dt *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*dt).Format(dateTimeFmtTicked)), nil
}

// UnmarshalJSON override unmarshalJSON
func (dt *DateTime) UnmarshalJSON(b []byte) error {
	ts, err := time.Parse(dateTimeFmtTicked, string(b))
	if err != nil {
		return err
	}

	*dt = DateTime(ts)
	return nil
}

// String returns it's string representation
func (dt *DateTime) String() string {
	return time.Time(*dt).Format(dateTimeFmt)
}

// GetBSON implements bson.Getter since the bson library does not look at underlying types and matches directly the time.Time type
func (dt DateTime) GetBSON() (interface{}, error) {
	return time.Time(dt), nil
}

// SetBSON implements bson.Setter since the bson library does not look at underlying types and matches directly the time.Time type
func (dt *DateTime) SetBSON(raw bson.Raw) error {
	decoded := time.Time{}

	bsonErr := raw.Unmarshal(&decoded)
	if bsonErr != nil {
		return bsonErr
	}
	*dt = DateTime(decoded)
	return nil

}
