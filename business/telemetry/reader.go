package telemetry

import (
	"fmt"
	"io"
	"io/ioutil"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Read(f io.Reader) (*Telem, error) {
	labels := []string{
		"ACCL",
		"DEVC",
		"DVID",
		"DVNM",
		"EMPT",
		"GPRO",
		"GPS5",
		"GPSF",
		"GPSP",
		"GPSU",
		"GYRO",
		"HD5.",
		"SCAL",
		"SIUN",
		"STRM",
		"TMPC",
		"TSMP",
		"UNIT",
		"TICK",
		"STNM",
		"ISOG",
		"SHUT",
		"TYPE",
		"FACE",
		"FCNM",
		"ISOE",
		"WBAL",
		"WRGB",
		"MAGN",
		"STMP",
		"STPS",
		"SROT",
		"TIMO",
		"UNIF",
		"MTRX",
		"ORIN",
		"ALLD",
		"ORIO",
	}

	label := make([]byte, 4, 4) // 4 byte ascii label of data
	desc := make([]byte, 4, 4)  // 4 byte description of length of data

	// keep a copy of the scale to apply to subsequent sentences
	s := Scale{}

	// the full telemetry for this period
	t := &Telem{}

	for {
		// pick out the label
		read, err := io.ReadFull(f, label)
		if err == io.EOF || read == 0 {
			return nil, err
		}

		label_string := string(label)

		if !stringInSlice(label_string, labels) {
			//skip label if not recognised, instead of stopping
			fmt.Printf("Could not find label in list: %s (%x)\n", label, label)
			//return nil, err
		}

		// pick out the label description
		read, err = io.ReadFull(f, desc)
		if err == io.EOF || read == 0 {
			break
		}

		// first byte is zero, there is no length
		if 0x0 == desc[0] {
			continue
		}

		// skip empty packets
		if "EMPT" == label_string {
			if _, err := io.CopyN(ioutil.Discard, f, 4); err != nil {
				return nil, err
			}
			continue
		}

		// extract the size and length
		val_size := int64(desc[1])
		num_values := (int64(desc[2]) << 8) | int64(desc[3])
		length := val_size * num_values

		// uncomment to see label, type, size and length
		//fmt.Printf("%s (%c) of size %v and len %v\n", label, desc[0], val_size, length)

		if "SCAL" == label_string {
			value := make([]byte, val_size*num_values, val_size*num_values)
			read, err = io.ReadFull(f, value)
			if err == io.EOF || read == 0 {
				return nil, err
			}

			// clear the scales
			s.Values = s.Values[:0]

			err := s.Parse(value, val_size)
			if err != nil {
				return nil, err
			}
		} else {
			value := make([]byte, val_size)

			for i := int64(0); i < num_values; i++ {
				read, err := io.ReadFull(f, value)
				if err == io.EOF || read == 0 {
					return nil, err
				}

				// I think DVID is the payload boundary; this might be a bad assumption
				if "DVID" == label_string {

					// XXX: I think this might skip the first sentence
					return t, nil
				} else if "GPS5" == label_string {
					g := GPS{}
					g.Parse(value, &s)
					t.Gps = append(t.Gps, g)
				} else if "GPSU" == label_string {
					g := GPST{}
					g.Parse(value)
					t.Time = g
				} else if "ACCL" == label_string {
					a := Accelerometer{}
					err := a.Parse(value, &s)
					if err != nil {
						return nil, err
					}
					t.Accl = append(t.Accl, a)
				} else if "TMPC" == label_string {
					tmp := Temperature{}
					tmp.Parse(value)
					t.Temp = tmp
				} else if "TSMP" == label_string {
					tsmp := TotalSamples{}
					tsmp.Parse(value, &s)
				} else if "GYRO" == label_string {
					g := Gyrometer{}
					err := g.Parse(value, &s)
					if err != nil {
						return nil, err
					}
					t.Gyro = append(t.Gyro, g)
				} else if "GPSP" == label_string {
					g := GPSPrecision{}
					err := g.Parse(value)
					if err != nil {
						return nil, err
					}
					t.GpsAccuracy = g
				} else if "GPSF" == label_string {
					g := GPSFix{}
					err := g.Parse(value)
					if err != nil {
						return nil, err
					}
					t.GpsFix = g
				} else if "UNIT" == label_string {
					// this is a string of units like "rad/s", not sure if it changes
					//fmt.Printf("\tUN: %s\n", value)
				} else if "SIUN" == label_string {
					// this is the SI unit - also not sure if it changes

					//fmt.Printf("\tSI: %s\n", value)
				} else if "DVNM" == label_string {
					// device name, "Camera"
					//fmt.Printf("\tvals: %s\n", value)
				} else if "TICK" == label_string {
					//fmt.Printf("\tvals: %s\n", value)
				} else if "STNM" == label_string {
					//fmt.Printf("\tvals: %s\n", value)
				} else if "ISOG" == label_string {
					//fmt.Printf("\tvals: %s\n", value)
				} else if "SHUT" == label_string {
					//fmt.Printf("\tvals: %s\n", value)
				} else {
					//fmt.Printf("\tvalue is %v\n", value)
				}
			}
		}

		// pack into 4 bytes
		mod := length % 4
		if mod != 0 {
			seek := 4 - mod
			if _, err := io.CopyN(ioutil.Discard, f, seek); err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}
