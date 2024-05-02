package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	// Fields of your config file (e.g., API keys, database settings)
	Inet       []string `json:"inet"`
	IP         string   `json:"ip"`
	Port       int      `json:"port"`
	Interfaces []string `json:"interfaces"`
}

func GetConfig() (*Config, error) {
	homeDir, _ := os.UserHomeDir()
	filename := homeDir + "/.config/raspberry.json"

	config, err := readExistingConfig(filename)
	if err != nil && os.IsNotExist(err) {
		RunConfig()
	}

	return config, err
}

func readExistingConfig(filename string) (*Config, error) {
	// Read the file contents into a string
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON data into a Config struct
	config := &Config{}
	json.Unmarshal(fileContents, config)

	return config, nil
}

func writeNewConfig(config *Config, filename string) error {
	// Marshal the Config struct into JSON data
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

func RunConfig() {
	homeDir, _ := os.UserHomeDir()
	filename := homeDir + "/.config/raspberry.json"

	config, err := readExistingConfig(filename)
	if err != nil && os.IsNotExist(err) {
		config = &Config{}
	} else if err != nil {
		fmt.Println("Error reading existing configuration:", err)
		return
	}

	fmt.Println("Welcome to the configuration setup!")
	fmt.Println("Enter your network interfaces connected to the internet:")
	inet := readInput()

	fmt.Println("Enter your network interfaces able to monitor and inject:")
	interfaces := readInput()

	fmt.Println("Enter the IP you would like the webapp to be on:")
	ip := readInput()

	fmt.Println("Enter the port you would like the webapp to be on:")
	port := readInput()

	// Update the Config struct with the user's input
	config.Inet = strings.Split(inet, " ")
	config.Interfaces = strings.Split(interfaces, " ")
	config.IP = ip
	config.Port, _ = strconv.Atoi(port)

	err = writeNewConfig(config, filename)
	if err != nil {
		fmt.Println("Error writing new configuration:", err)
		return
	}

	fmt.Println("Configuration saved successfully!")
}

func readInput() string {
	var input string
	_, err := fmt.Scan(&input)
	if err != nil {
		panic(err)
	}
	return input
}
