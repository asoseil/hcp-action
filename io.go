package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	TargetEnv = "env"
	TargetOut = "out"
)

// Set an output var using env files
func setOutVar(name, value string) {
	githubOutput := os.Getenv("GITHUB_OUTPUT")
	if githubOutput == "" {
		log.Fatalf("GITHUB_OUTPUT environment variable is not set")
	}

	// Write to the environment file
	outputLine := fmt.Sprintf("%s=%s\n", name, value)
	file, err := os.OpenFile(githubOutput, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open GITHUB_OUTPUT file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close GITHUB_ENV file: %s", err)
		}
	}(file)

	_, err = file.WriteString(outputLine)
	if err != nil {
		log.Fatalf("Failed to write to GITHUB_OUTPUT file: %s", err)
	}
}

// Set an environment var using env files
func setEnvVar(name, value string) {
	githubEnv := os.Getenv("GITHUB_ENV")
	if githubEnv == "" {
		log.Fatalf("GITHUB_ENV environment variable is not set")
	}

	// Write to the environment file
	envLine := fmt.Sprintf("%s=%s\n", strings.ToUpper(name), value)
	file, err := os.OpenFile(githubEnv, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open GITHUB_ENV file: %s", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Failed to close GITHUB_ENV file: %s", err)
		}
	}(file)

	_, err = file.WriteString(envLine)
	if err != nil {
		log.Fatalf("Failed to write to GITHUB_ENV file: %s", err)
	}

	// Mask the value using ::add-mask command
	_, err = mask(value)
	if err != nil {
		log.Fatalf("Failed to mask environment variable: %s", err)
	}
}

func mask(value string) (int, error) {
	return fmt.Printf("::add-mask::%s\n", value)
}

type Config map[string][]string

func parseConfig(jsonInput string) (Config, error) {
	unescapedInput := strings.ReplaceAll(jsonInput, "\\", "")

	var config Config

	err := json.Unmarshal([]byte(unescapedInput), &config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
