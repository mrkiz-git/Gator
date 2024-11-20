package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	db_url            string `json:"db_url"`
	current_user_name string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error retrieving home directory: %v", err)
		return "", err
	}

	// Construct the full path to the configuration file
	configFilePath := filepath.Join(homeDir, "configFileName")
	log.Printf("Config file path resolved to: %s", configFilePath)

	return configFilePath, nil
}

func write(config *Config) error {

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Convert the Config struct to JSON
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Println("Failed to marshal JSON: %v", err)
		return err
	}

	// Create or open the file
	file, err := os.Create(filePath) // Creates the file or overwrites if it exists
	if err != nil {
		log.Println("Failed to create file: %v", err)
		return err
	}
	defer file.Close() // Ensure the file is closed when we're done

	// Write JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Println("Failed to write to file: %v", err)
		return err
	}

	log.Println("JSON written successfully to config.json")
	return nil
}

func (c *Config) Read() (*Config, error) {

	filePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading file at %s: %v", filePath, err)
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(fileContent, &cfg); err != nil {
		log.Printf("Failed to decode JSON: %v", err)
		return nil, err
	}

	log.Printf("Config loaded successfully: %+v", cfg)
	return &cfg, nil
}

func (c *Config) SetUser(userName string) error {

	log.Printf("Setting current_user_name ")
	c.current_user_name = userName

	if err := write(c); err != nil {
		log.Printf("Failed to write JSON json file: %v", err)
		return err
	}

	return nil
}
