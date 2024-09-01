package validate_config

// https://www.kelche.co/blog/go/ini/

import (
	"fmt"
	"gopkg.in/ini.v1"
)

type RemoteConfigEntry struct {
	Type       string `yaml:"type"`
	Account    string `yaml:"account"`
	Key        string `yaml:"key"`
	HardDelete bool   `yaml:"hard_delete"`
}

type RemoteConfig struct {
	Remotes map[string]RemoteConfigEntry
}

func printIniFile(cfg *ini.File) {
	// For debugging
	for _, section := range cfg.Sections() {
		fmt.Printf("Section: [%s]\n", section.Name())
		for _, key := range section.Keys() {
			fmt.Printf("  %s = %s\n", key.Name(), key.String())
		}
	}
}

func ValidateConfig(configPath string) error {
	cfg, err := ini.Load(configPath)
	fmt.Printf("Reading from path: %s\n", configPath)

	if err != nil {
		return err
	}
	config, err := parseConfig(cfg)
	err = validateConfig(&config)
	if err != nil {
		return err
	}
	fmt.Println("Config is valid")
	return nil
}

func parseConfig(cfg *ini.File) (RemoteConfig, error) {
	mandatoryKeys := []string{"type", "account", "key"}
	config := RemoteConfig{
		Remotes: make(map[string]RemoteConfigEntry),
	}

	for _, section := range cfg.Sections() {

		if section.Name() == "DEFAULT" {
			continue
		}

		remoteName := section.Name()

		fmt.Printf("\nProcessing section: %v\n", section.Name())

		keyValues := make(map[string]interface{})
		fmt.Printf("")
		keyValues["hard_delete"] = section.Key("hard_delete").MustBool(false)

		for _, keyName := range mandatoryKeys {
			keyValue := section.Key(keyName).String()
			if keyValue == "" {
				return RemoteConfig{}, fmt.Errorf("Config failed for remote %s", remoteName)
			}
			keyValues[keyName] = keyValue
		}

		entry := RemoteConfigEntry{
			Type:       keyValues["type"].(string),
			Account:    keyValues["account"].(string),
			Key:        keyValues["key"].(string),
			HardDelete: keyValues["hard_delete"].(bool),
		}
		config.Remotes[remoteName] = entry
	}
	return config, nil
}

func validateConfig(config *RemoteConfig) error {
	for remoteName, remoteConfigEntry := range config.Remotes {
		if remoteConfigEntry.Type != "b2" {
			return fmt.Errorf("Config validation failed for %s", remoteName)
		}
		// TODO: More validation can go here
	}
	return nil
}
