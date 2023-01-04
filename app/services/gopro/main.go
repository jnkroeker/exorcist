package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jnkroeker/exorcist/business/telemetry"
)

type data struct {
	Data []telemetry.TelemOut `json:"data"`
}

func main() {
	inName := flag.String("i", "", "Required: telemetry file")
	outName := flag.String("o", "", "Required: json file to write")
	flag.Parse()

	if *inName == "" {
		flag.Usage()
		return
	}

	telemFile, err := os.Open(*inName)
	if err != nil {
		fmt.Printf("cannot access telemetry file %s.\n", *inName)
		os.Exit(1)
	}
	defer telemFile.Close()

	var d data

	t := &telemetry.Telem{}
	t_prev := &telemetry.Telem{}

	for {
		t, err = telemetry.Read(telemFile)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading telemetry file", err)
			os.Exit(1)
		} else if err == io.EOF || t == nil {
			break
		}

		if t_prev.IsZero() {
			*t_prev = *t
			t.Clear()
			continue
		}

		t_prev.FillTimes(t.Time.Time)

		telems := t_prev.OutputJson()

		d.Data = append(d.Data, telems...)

		*t_prev = *t
		t = &telemetry.Telem{}
	}

	jsonFile, err := os.Create(*outName)
	if err != nil {
		fmt.Printf("Cannot make output file %s.\n", *outName)
		os.Exit(1)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Cannot close json file %s: %s", file.Name(), err)
			os.Exit(1)
		}
	}(jsonFile)

	if err := json.NewEncoder(jsonFile).Encode(d); err != nil {
		fmt.Println("Error encoding output json", err)
		os.Exit(1)
	}
}
