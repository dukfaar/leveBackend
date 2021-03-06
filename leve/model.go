package leve

import (
	"github.com/dukfaar/goUtils/relay"
	"github.com/globalsign/mgo/bson"
)

type Model struct {
	ID      bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	XivdbID string        `json:"xivdbid,omitempty" bson:"xivdbid,omitempty"`
	Name    string        `json:"name,omitempty" bson:"name,omitempty"`
	Level   int32         `json:"level,omitempty" bson:"level,omitempty"`
	Xp      int32         `json:"xp,omitempty" bson:"xp,omitempty"`
	Class   string        `json:"class,omitempty" bson:"class,omitempty"`
	Gil     int32         `json:"gil,omitempty" bson:"gil,omitempty"`
}

var GraphQLType = `
type Leve {
	_id: ID
	xivdbid: String
	name: String
	level: Int
	xp: Int
	class: String
	gil: Int
}
` +
	relay.GenerateConnectionTypes("Leve")
