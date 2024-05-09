package modules

import (
	"encoding/csv"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/5N41P4/raspberry/internal/data"
)

func ParseCSV(path string) ([]data.AppAP, []data.AppClient, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return []data.AppAP{}, []data.AppClient{}, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	// Read the file content
	records, err := reader.ReadAll()
	if err != nil {
		return []data.AppAP{}, []data.AppClient{}, err
	}

	// Find the row that separates the access points and clients data
	var separatorIndex int
	for i, record := range records {
		if strings.Contains(record[0], "Station MAC") {
			separatorIndex = i
			break
		}
	}

	// Split the records into two slices
	if len(records) <= 0 {
		return []data.AppAP{}, []data.AppClient{}, errors.New("nothing to parse")
	}
	apRecords := records[1:separatorIndex]
	clientRecords := records[separatorIndex+1:]

	// Start with the APs
	var aps []data.AppAP
	for _, record := range apRecords {
		// Okay, fill an AP struct then append to the dump
		for i, str := range record {
			record[i] = strings.Trim(str, " ")
		}

		channel, _ := strconv.Atoi(record[3])
		speed, _ := strconv.Atoi(record[4])
		power, _ := strconv.Atoi(record[8])
		beacons, _ := strconv.Atoi(record[9])
		ivs, _ := strconv.Atoi(record[10])
		idlen, _ := strconv.Atoi(record[12])

		cur_ap := data.AppAP{
			Bssid:   record[0],
			First:   record[1],
			Last:    record[2],
			Channel: channel,
			Speed:   speed,
			Privacy: record[5],
			Cipher:  record[6],
			Auth:    record[7],
			Power:   power,
			Beacons: beacons,
			IVs:     ivs,
			Lan:     strings.Replace(record[11], " ", "", -1), // Clean blanks
			IdLen:   idlen,
			Essid:   record[13],
			Key:     record[14],
		}

		aps = append(aps, cur_ap)
	}

	// Continue with the clients
	var cls []data.AppClient
	for _, record := range clientRecords {
		// Okay, fill a Client struct then append to the dump
		for i, str := range record {
			record[i] = strings.Trim(str, " ")
		}

		power, _ := strconv.Atoi(record[3])
		packets, _ := strconv.Atoi(record[4])

		cur_client := data.AppClient{
			Station: record[0],
			First:   record[1],
			Last:    record[2],
			Power:   power,
			Packets: packets,
			Bssid:   record[5],
			Probed:  record[6],
		}

		cls = append(cls, cur_client)
	}

	return aps, cls, nil
}
