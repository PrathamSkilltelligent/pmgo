package utils

import (
	"time"

	"github.com/google/uuid"
)

func ToPtr[T any](data T) *T {
	return &data
}

func CombineList[T any](t ...[]T) []T {
	list := make([]T, 0)
	for idx := range t {
		list = append(list, t[idx]...)
	}
	return list
}

func ToByteArrayPtr(id *uuid.UUID) *[]byte {
	if id == nil {
		return nil
	}
	return ToPtr(id[:])
}

func ToByteArray(id uuid.UUID) []byte {
	return id[:]
}

func ToUUID(barray []byte) uuid.UUID {
	return uuid.UUID(barray)
}

func ToUUIDPtr(barray *[]byte) *uuid.UUID {
	if barray == nil {
		return nil
	}
	return ToPtr(uuid.UUID(*barray))
}

func ToDateTime(milliSec int64) time.Time {
	return time.Unix(0, milliSec*int64(time.Millisecond))
}

func ToDateTimePtr(milliSec *int64) *time.Time {
	if milliSec == nil {
		return nil
	}
	return ToPtr(time.Unix(0, *milliSec*int64(time.Millisecond)))
}
