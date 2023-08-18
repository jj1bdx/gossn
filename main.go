package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"log"
	// "net/http"
	"strconv"
	"strings"
	// "time"
)

var urlSWPCFlux string = "https://services.swpc.noaa.gov/json/f107_cm_flux.json"
var urlSWPCKp string = "https://services.swpc.noaa.gov/json/planetary_k_index_1m.json"
var urlSILSOEisn string = "https://sidc.be/SILSO/DATA/EISN/EISN_current.csv"

type Flux struct {
	Time string  `json:"time_tag"`
	Flux float64 `json:"flux"`
}

type Kp struct {
	Time string  `json:"time_tag"`
	Kp   float64 `json:"estimated_kp"`
}

func main() {

	fluxData := getURLBody(urlSWPCFlux)
	// kpData := getURLBody(urlSWPCKp)
	// silsoData := getURLBody(urlSILSOEisn)

	var fluxList []Flux
	if err := json.Unmarshal(fluxData, &fluxList); err != nil {
		log.Fatal(err)
	}
	fluxLatest := fluxList[0]
	fmt.Printf("%s : %f\n", fluxLatest.Time, fluxLatest.Flux)

	kpData := getURLBody(urlSWPCKp)

	var kpList []Kp
	if err := json.Unmarshal(kpData, &kpList); err != nil {
		log.Fatal(err)
	}
	kpLatest := kpList[len(kpList) - 1]
	fmt.Printf("%s : %f\n", kpLatest.Time, kpLatest.Kp)

	silsoData := getURLBody(urlSILSOEisn)
	eisnReader := csv.NewReader(strings.NewReader(string(silsoData[:])))
	eisnRecords, err := eisnReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	eisnLatestRecord := eisnRecords[len(eisnRecords) - 1]
	eisnYear := trimLeftSpace(eisnLatestRecord[0])
	eisnMonth := trimLeftSpace(eisnLatestRecord[1])
	eisnDay := trimLeftSpace(eisnLatestRecord[2])
	var eisnVal int
	if eisnVal, err = strconv.Atoi(trimLeftSpace(eisnLatestRecord[4])); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s:%s:%sT00:00:00 : %d\n",
		eisnYear, eisnMonth, eisnDay, eisnVal)
	
}
