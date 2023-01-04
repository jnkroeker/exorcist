package telemetry

import (
	"errors"
	"time"
)

type GPST struct {
	Time time.Time
}

func (gpsu *GPST) Parse(bytes []byte) error {
	if 16 != len(bytes) {
		return errors.New("invalid length GPSU packet")
	}

	t, err := time.Parse("060102150405", string(bytes))
	if err != nil {
		return err
	}

	gpsu.Time = t

	return nil
}
