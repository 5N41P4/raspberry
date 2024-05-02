package modules

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/5N41P4/rpine/internal/data"
)

func ParseCSV(path string) ([]data.AppAP, []data.AppClient, error) {
	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return []data.AppAP{}, []data.AppClient{}, err
	}
	defer file.Close()

	// Read the file content
	dump, err := io.ReadAll(file)
	if err != nil {
		return []data.AppAP{}, []data.AppClient{}, err
	}

	dump_str := string(dump)
	// Replace endline with just an  \n
	dump_str = strings.Replace(dump_str, ", \r\n", ", \n", -1)
	dump_str = strings.Replace(dump_str, ",\r\n", ",\n", -1)
	dump_split := strings.SplitN(dump_str, "\r\n", 4)

	// Extract the two parts of the csv
	dump_aps := dump_split[2]
	dump_clients := dump_split[3]
	dump_clients = strings.SplitN(dump_clients, "\r\n", 2)[1]

	// End of dirty hack, fill the structs
	reader_aps := csv.NewReader(strings.NewReader(dump_aps))
	reader_aps.Comma = ','
	reader_aps.TrimLeadingSpace = true

	reader_clients := csv.NewReader(strings.NewReader(dump_clients))
	reader_clients.Comma = ','
	reader_clients.TrimLeadingSpace = true

	// Start with the APs
	var aps []data.AppAP
	for {
		record, err := reader_aps.Read()
		if err == io.EOF {
			break
		}
		// Okay, fill an AP struct then append to the dump

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
	for {
		record, err := reader_clients.Read()
		if err == io.EOF {
			break
		}

		// Okay, fill a Client struct then append to the dump
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
