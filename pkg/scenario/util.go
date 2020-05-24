package scenario

import (
	"github.com/pkg/errors"

	"github.com/docker/go-units"
)

// GetInBytes takes in a formatted size string and returns the size specified by this string in bytes.
// It errors if the size is below zero.
func GetInBytes(unit string, unitName string) (int64, error) {
	unitInt, err := units.RAMInBytes(unit)
	if err != nil {
		return 0, errors.Wrap(err, "converting storage string to bytes failed")
	} else if unitInt < 0 {
		return 0, errors.Errorf("%s usage should be at least 0", unitName)
	}
	return unitInt, nil
}
