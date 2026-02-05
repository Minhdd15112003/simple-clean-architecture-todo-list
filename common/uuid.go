package common

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcutil/base58"
)

type UID struct {
	localID    uint32
	objectType int
	shardID    uint32
}

func NewUID(localID uint32, objectType int, shardID uint32) UID {
	return UID{
		localID:    localID,
		objectType: objectType,
		shardID:    shardID,
	}
}

func (uid UID) String() string {
	val := uint64(uid.localID)<<28 | uint64(uid.objectType)<<18 | uint64(uid.shardID)<<0
	return base58.Encode([]byte(fmt.Sprintf("%d", val)))
}

func (uid UID) GetLocalID() uint32 {
	return uid.localID
}

func (uid UID) GetShardID() uint32 {
	return uid.shardID
}

func (uid UID) GetObjectType() int {
	return uid.objectType
}

func DecomposeUID(s string) (UID, error) {
	// Decode base58 string về bytes
	decoded := base58.Decode(s)
	if len(decoded) == 0 {
		return UID{}, errors.New("invalid base58 string")
	}

	// Parse bytes thành uint64
	uid, err := strconv.ParseUint(string(decoded), 10, 64)
	if err != nil {
		return UID{}, err
	}

	if (1 << 18) > uid {
		return UID{}, errors.New("wrong uid")
	}

	u := UID{
		localID:    uint32(uid >> 28),
		objectType: int(uid >> 18 & 0x3FF),
		shardID:    uint32(uid >> 0 & 0x3FFFF),
	}
	return u, nil
}

func (uid UID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", uid.String())), nil
}

func (uid *UID) UnmarshalJSON(data []byte) error {
	// Remove quotes
	s := string(data[1 : len(data)-1])

	u, err := DecomposeUID(s)
	if err != nil {
		return err
	}

	*uid = u
	return nil
}

// FromBase58 decodes a base58 encoded UID string to real local ID
func FromBase58(s string) (uint32, error) {
	uid, err := DecomposeUID(s)
	if err != nil {
		return 0, err
	}
	return uid.GetLocalID(), nil
}
