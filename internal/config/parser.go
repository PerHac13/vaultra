package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Parser struct {}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseFile(filePath string) (*ConfigType, error) {

	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("config file not found: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg ConfigType

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &cfg, nil

}

func (p *Parser) ParseDir(dirPath string) ([]*ConfigType, error){
	var configs []*ConfigType

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" && ext != ".json" {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())
		cfg, err := p.ParseFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse config file %s: %v", filePath, err)
		}

		configs = append(configs, cfg)
	}

	return configs, nil
}