package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FilterList struct {
	ApFilter *Filter `json:"ApFilter"`
	ClFilter *Filter `json:"ClFilter"` }

type Filter struct {
	IsWhitelist bool     `json:"isWhitelist"`
	Identifier  []string `json:"identifier"`
}

// Create new Filterlists
func newFilterList() *FilterList {
  // Check if safed filters are available
  filter, err := checkStorage()
  if err == nil {
    return filter
  }

  // If no filters were safed, create new ones
	return &FilterList{
		ApFilter: newFilter(),
    ClFilter: newFilter(),
	}
}

func newFilter() *Filter {
	return &Filter{
		IsWhitelist: false,
		Identifier:  make([]string, 0),
	}
}


// Switch the lists from White to Black lists
func (fl *Filter) Switch() {
	fl.IsWhitelist = !fl.IsWhitelist
}

// Add and remove accesspoints and clients
func (fl *Filter) Add(id string) {
	fl.Identifier = append(fl.Identifier, id)
}

func (fl *Filter) Delete(id string) {
	for i, filter := range fl.Identifier {
		if id == filter {
			fl.Identifier = removeElementInPlace(fl.Identifier, i)
			return
		}
	}
}

func (fl *Filter) Reset(){
  fl.Identifier = make([]string, 0)
}

// Apply the Filterlists
func (fl *Filter) IsAllowed(id string) bool {
	var exists bool
	for _, filter := range fl.Identifier {
		if id == filter{
			exists = true
			break
		}
	}
	return fl.IsWhitelist == exists
}

// Cleanup Function for FilterList
func (fl *FilterList) cleanup() {
  	// Marshal the struct into a byte slice
	jsonData, err := json.Marshal(fl)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Open a file for writing
	file, err := os.Create("/usr/local/raspberry/FilterList.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close() // Close the file on exit

	// Write the JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}

func checkStorage() (*FilterList, error) {
  // Open the file for reading
  file, err := os.Open("/usr/local/raspberry/FilterList.json")
  if err != nil {
    return nil, fmt.Errorf("Error opening file: %w", err)
  }
  defer file.Close() // Close the file on exit

  // Create a new FilterList object
  var fl FilterList

  // Decode the JSON data from the file
  decoder := json.NewDecoder(file)
  err = decoder.Decode(&fl)
  if err != nil {
    return nil, fmt.Errorf("Error decoding JSON: %w", err)
  }

  return &fl, nil
}

// Helperfunctions
func removeElementInPlace(slice []string, index int) []string {
	if index >= 0 && index < len(slice) {
		for i := index; i < len(slice)-1; i++ {
			slice[i] = slice[i+1]
		}
		slice = slice[:len(slice)-1] // Reduce slice length by 1
	}
	return slice
}
