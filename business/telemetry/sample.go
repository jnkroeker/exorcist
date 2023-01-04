package telemetry

import (
	"encoding/binary"
	"errors"
)

// Total number of samples
type TotalSamples struct {
	Samples uint32
}

func (t *TotalSamples) Parse(bytes []byte, scale *Scale) error {
	if 4 != len(bytes) {
		return errors.New("invalid length TSMP packet")
	}

	t.Samples = binary.BigEndian.Uint32(bytes[0:4])

	return nil
}
