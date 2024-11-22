package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

// getConfigFilePath returns the full path to the configuration file in the user's home directory.
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error retrieving home directory: %v", err)
		return "", err
	}

	// Construct the full path to the configuration file
	configFilePath := filepath.Join(homeDir, configFileName)
	log.Printf("Config file path resolved to: %s", configFilePath)

	return configFilePath, nil
}

// writeConfig writes the given Config object to the configuration file in JSON format.
func writeConfig(config *Config) error {
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// Convert the Config struct to JSON
	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return err
	}

	// Create or open the file
	file, err := os.Create(filePath) // Overwrites the file if it exists
	if err != nil {
		log.Printf("Failed to create file: %v", err)
		return err
	}
	defer file.Close()

	// Write JSON data to the file
	_, err = file.Write(jsonData)
	if err != nil {
		log.Printf("Failed to write to file: %v", err)
		return err
	}

	log.Println("JSON written successfully to", configFileName)
	return nil
}

// LoadConfig reads and returns the configuration from the configuration file.
func LoadConfig() (*Config, error) {
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

// SetUser updates the CurrentUserName in the configuration and saves it to the file.
func (c *Config) SetUser(userName string) error {
	log.Printf("Setting CurrentUserName to %s", userName)
	c.CurrentUserName = userName

	if err := writeConfig(c); err != nil {
		log.Printf("Failed to write config to file: %v", err)
		return err
	}

	return nil
}
