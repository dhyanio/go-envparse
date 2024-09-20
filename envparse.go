package envparse

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// EnvVar defines the structure to hold environment variable key and optional default value
type EnvVar struct {
	Key         string
	DefaultValue string
	Validate    func(string) bool // Optional validation function for the environment variable
}

// ParseEnvironment parses and sets the specified environment variables, with optional defaults and validation
func ParseEnvironment(vars []EnvVar, envFilePath string) {
	// Check if .env file exists at the given path
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		log.Printf(".env file not found at %s. Relying on global environment variables.", envFilePath)
	}

	for _, envVar := range vars {
		value := os.Getenv(envVar.Key)
		if value == "" {
			setEnvVariable(envVar.Key, envFilePath)
			value = os.Getenv(envVar.Key)
		}

		// Apply default if still empty after reading from .env file
		if value == "" && envVar.DefaultValue != "" {
			os.Setenv(envVar.Key, envVar.DefaultValue)
			value = envVar.DefaultValue
			log.Printf("Using default value for %s", envVar.Key)
		}

		// Validation check if a validation function is provided
		if envVar.Validate != nil && !envVar.Validate(value) {
			log.Printf("Validation failed for %s with value '%s'. Exiting.", envVar.Key, value)
			os.Exit(1)
		}

		// Final check if the variable is still empty, exit if required
		if os.Getenv(envVar.Key) == "" {
			log.Printf("Could not resolve a %s environment variable. Exiting.", envVar.Key)
			os.Exit(1)
		} else {
			log.Printf("Successfully loaded %s", envVar.Key)
		}
	}
}

// setEnvVariable sets the environment variable by reading the .env file
func setEnvVariable(env string, envFilePath string) {
	// Open the .env file and scan for the variable
	file, err := os.Open(envFilePath)
	if err != nil {
		log.Printf("Error opening .env file: %v", err)
		return
	}
	defer file.Close()

	lookInFile := bufio.NewScanner(file)
	lookInFile.Split(bufio.ScanLines)

	for lookInFile.Scan() {
		parts := strings.SplitN(lookInFile.Text(), "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		if strings.EqualFold(key, env) { // Case-insensitive comparison
			os.Setenv(key, value)
		}
	}
}

// Example usage of ParseEnvironment function
func main() {
	// Define the variables to be parsed, including optional default values and validation
	envVars := []EnvVar{
		{
			Key:          "CLIENT_ID",
			DefaultValue: "default_client_id",
			Validate: func(value string) bool {
				return len(value) > 0 // Ensure it is non-empty
			},
		},
		{
			Key:          "CLIENT_SECRET",
			DefaultValue: "default_client_secret",
			Validate: func(value string) bool {
				return len(value) >= 8 // Ensure it has at least 8 characters
			},
		},
		{
			Key:          "ISSUER",
			DefaultValue: "https://default-issuer.com",
			Validate: func(value string) bool {
				return strings.HasPrefix(value, "https://") // Ensure it starts with https
			},
		},
	}

	// Call the parser with a custom .env file path
	ParseEnvironment(envVars, ".env")
}
