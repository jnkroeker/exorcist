package telemetry

import (
	"time"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geo"
)

// Represents one second of telemetry data
type Telem struct {
	Accl        []Accelerometer
	Gps         []GPS
	Gyro        []Gyrometer
	GpsFix      GPSFix
	GpsAccuracy GPSPrecision
	Time        GPST
	Temp        Temperature
}

// the json we want
// GPS data might have generated a timestamp and derived track
type TelemOut struct {
	*GPS

	GpsAccuracy uint16  `json:"gps_accuracy,omitempty"`
	GpsFix      uint32  `json:"gps_fix,omitempty"`
	Temp        float32 `json:"temp,omitempty"`
	Track       float64 `json:"track,omitempty"`
}

var lastGoodTrack float64

var pp orb.Point

func (t *Telem) Clear() {
	t.Accl = t.Accl[:0]
	t.Gps = t.Gps[:0]
	t.Gyro = t.Gyro[:0]
	t.Time.Time = time.Time{}
}

func (t *Telem) IsZero() bool {
	return t.Time.Time.IsZero()
}

// try to populate a timestamp for every GPS row. probably bogus.
func (t *Telem) FillTimes(until time.Time) error {
	len := len(t.Gps)
	diff := until.Sub(t.Time.Time)

	offset := diff.Seconds() / float64(len)

	for i, _ := range t.Gps {
		dur := time.Duration(float64(i)*offset*1000) * time.Millisecond
		ts := t.Time.Time.Add(dur)
		t.Gps[i].Timestamp = ts.UnixNano() / 1000
	}

	return nil
}

func (t *Telem) OutputJson() []TelemOut {
	var out []TelemOut

	for i, _ := range t.Gps {
		jobj := TelemOut{&t.Gps[i], 0, 0, 0, 0}
		if 0 == i {
			jobj.GpsAccuracy = t.GpsAccuracy.Accuracy
			jobj.GpsFix = t.GpsFix.F
			jobj.Temp = t.Temp.Temp
		}

		p := orb.Point{jobj.GPS.Longitude, jobj.GPS.Latitute}
		jobj.Track = geo.Bearing(pp, p)
		pp = p

		if jobj.Track < 0 {
			jobj.Track = 360 + jobj.Track
		}

		// only set the track if speed is over 1 m/s
		// if it's slower (eg, stopped) it will drift all over with the location
		if jobj.GPS.Speed > 1 {
			lastGoodTrack = jobj.Track
		} else {
			jobj.Track = lastGoodTrack
		}

		out = append(out, jobj)
	}

	return out
}
