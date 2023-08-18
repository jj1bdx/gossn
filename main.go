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

type Eisn struct {
	Time string
	Eisn float64
}

func main() {

	fluxData := getURLBody(urlSWPCFlux)
	var fluxList []Flux
	if err := json.Unmarshal(fluxData, &fluxList); err != nil {
		log.Fatal(err)
	}
	fluxLatest := fluxList[0]
	fmt.Printf("10.7cm Solar Flux Index (SFI): %s : %d\n",
		fluxLatest.Time, int(fluxLatest.Flux))
	fmt.Printf("Estimated SSN from SFI: %s : %d\n",
		fluxLatest.Time, estimatedSSN(fluxLatest.Flux))

	kpData := getURLBody(urlSWPCKp)

	var kpList []Kp
	if err := json.Unmarshal(kpData, &kpList); err != nil {
		log.Fatal(err)
	}
	kpLatest := kpList[len(kpList)-1]
	fmt.Printf("Estimated Kp: %s : %g\n", kpLatest.Time, kpLatest.Kp)

	silsoData := getURLBody(urlSILSOEisn)
	eisnReader := csv.NewReader(strings.NewReader(string(silsoData[:])))
	eisnRecords, err := eisnReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	eisnLatestRecord := eisnRecords[len(eisnRecords)-1]
	eisnYear := trimLeftSpace(eisnLatestRecord[0])
	eisnMonth := trimLeftSpace(eisnLatestRecord[1])
	eisnDay := trimLeftSpace(eisnLatestRecord[2])
	var eisnVal int
	if eisnVal, err = strconv.Atoi(trimLeftSpace(eisnLatestRecord[4])); err != nil {
		log.Fatal(err)
	}
	var eisnLatest Eisn
	eisnLatest.Time = fmt.Sprintf("%s-%s-%sT00:00:00",
		eisnYear, eisnMonth, eisnDay)
	eisnLatest.Eisn = float64(eisnVal)
	fmt.Printf("EISN (observed SSN): %s : %d\n", eisnLatest.Time, int(eisnLatest.Eisn))

	fmt.Println("SFI/Kp source: NOAA SWPC")
	fmt.Println("EISN source: WDC-SILSO, Royal Observatory of Belgium, Brussels")
}
