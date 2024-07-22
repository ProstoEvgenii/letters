package functions

import (
	"letters/db"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var AuthUsers = map[string]int64{}

func CheckAuthUsers() {
	now := time.Now()
	timestampNow := now.Unix()
	for item := range AuthUsers {
		if timestampNow > AuthUsers[item]+(60*60*24) {
			delete(AuthUsers, item)
		}
	}
}
func CheckInDB(tmp interface{}, filter bson.M, collName string) bool {
	find := db.FindOne(filter, collName)
	if err := find.Decode(&tmp); err != nil {
		return false
	}
	return true
}
