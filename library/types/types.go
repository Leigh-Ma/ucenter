package types

import (
	"fmt"
	"strconv"
)

const (
	InvalidIdString = IdString("0x0")
)

type ObjectID uint64

func (id ObjectID) ToIdString() IdString {
	return IdString(fmt.Sprintf("0x%x", id))
}

type IdString string

func (id IdString) ToObjectID() ObjectID {
	value, err := strconv.ParseUint(string(id), 0, 64)
	if err != nil {
		fmt.Printf("prase idstring %s to object id error: %s\n", id, err.Error())
	}
	return ObjectID(value)
}

func (id IdString) ToString() string {
	return string(id)
}

func (id IdString) IsValid() bool {
	v, err := strconv.ParseUint(string(id), 0, 64)
	return err == nil && v != 0
}
