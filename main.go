/*
Copyright 2016, RadiantBlue Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/venicegeo/pzsvc-sdk-go/s3"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// LasHeader is the LAS header.
type LasHeader struct {
	FileSignature                [4]byte
	FileSourceID                 uint16
	GlobalEncoding               uint16
	GUID1                        uint32
	GUID2                        uint16
	GUID3                        uint16
	GUID4                        [8]uint8
	VersionMajor, VersionMinor   uint8
	SystemIdentifier             [32]byte
	GeneratingSoftware           [32]byte
	FileCreationDayOfYear        uint16
	FileCreationYear             uint16
	HeaderSize                   uint16
	OffsetToPointData            uint32
	NumberOfVLRs                 uint32
	PointDataRecordFormat        uint8
	PointDataRecordLength        uint16
	LegacyNumberOfPointRecords   uint32
	LegacyNumberOfPointsByReturn [5]uint32
	XScale, YScale, ZScale       float64
	XOffset, YOffset, ZOffset    float64
	MaxX, MinX                   float64
	MaxY, MinY                   float64
	MaxZ, MinZ                   float64
	// missing some waveform packet stuff for 1.4
}

// Format0 is Point Data Record Format 0.
type Format0 struct {
	X, Y, Z        int32
	Intensity      uint16
	Foo            byte
	Classification uint8
	ScanAngleRank  byte
	UserData       uint8
	PointSourceID  uint16
}

// Format1 is Point Data Record Format 1.
type Format1 struct {
	Format0
	GPSTime float64
}

// Format3 is Point Date Record Format 3.
type Format3 struct {
	Format1
	Red, Blue, Green uint16
}

// ReadLas reads an LAS file.
func ReadLas(fname string) (h LasHeader, p []Format1, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return h, p, err
	}
	defer f.Close()

	binary.Read(f, binary.LittleEndian, &h)
	fmt.Printf("%+v\n", h)

	f.Seek(int64(h.OffsetToPointData), 0)
	if h.PointDataRecordFormat == 0 {
		fmt.Println("format = 0")
		points := make([]Format0, h.LegacyNumberOfPointRecords)
		binary.Read(f, binary.LittleEndian, points)
		fmt.Printf("%+v\n", points[:5])
		fmt.Printf("%+v\n", points[len(points)-5:])
	} else if h.PointDataRecordFormat == 1 {
		fmt.Println("format = 1")
		p = make([]Format1, h.LegacyNumberOfPointRecords)
		binary.Read(f, binary.LittleEndian, p)
		fmt.Printf("%+v\n", p[:5])
		fmt.Printf("%+v\n", p[len(p)-5:])
	} else if h.PointDataRecordFormat == 3 {
		fmt.Println("format = 3")
		points := make([]Format3, h.LegacyNumberOfPointRecords)
		binary.Read(f, binary.LittleEndian, points)
		fmt.Printf("%+v\n", points[:5])
		fmt.Printf("%+v\n", points[len(points)-5:])
	}
	return h, p, nil
}

func main() {
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
		fmt.Fprintf(w, "Hi!")
	})

	router.GET("/info", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var inputName string
		var fileIn *os.File

		// Split the source S3 key string, interpreting the last element as the
		// input filename. Create the input file, throwing 500 on error.
		inputName = s3.ParseFilenameFromKey("pointcloud/samp11-utm.las")
		fileIn, err := os.Create(inputName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer fileIn.Close()

		// // Download the source data from S3, throwing 500 on error.
		// err = s3.Download(fileIn, "venicegeo-sample-data", "pointcloud/samp11-utm.las")
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
		//
		// h, _, err := ReadLas(inputName)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// if err := json.NewEncoder(w).Encode(h); err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// }
	})

	var defaultPort = os.Getenv("PORT")
	if defaultPort == "" {
		defaultPort = "8080"
	}

	log.Println("Starting on ", defaultPort)
	if err := http.ListenAndServe(":"+defaultPort, router); err != nil {
		log.Fatal(err)
	}
}
