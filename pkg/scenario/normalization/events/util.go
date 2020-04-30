package events

import (
	"fmt"
	"time"

	"github.com/docker/go-units"
)

// DesugarTimestamp takes in a formatted duration string and returns the duration specified by this string in milliseconds.
func DesugarTimestamp(t string) (int, error) {
	duration, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}

	return int(duration.Milliseconds()), nil
}

// GetInBytes takes in a formatted size string and returns the size specified by this string in bytes.
// It errors if the size is below zero.
func GetInBytes(unit string, unitName string) (int64, error) {
	unitInt, err := units.RAMInBytes(unit)
	if err != nil {
		return 0, err
	} else if unitInt < 0 {
		return 0, fmt.Errorf("%s usage should be at least 0", unitName)
	}
	return unitInt, nil
}
