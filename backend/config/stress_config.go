package config

// StressConfig represents the configuration for stress tests
type StressConfig struct {
	TargetURL string           `json:"target_url"`
	Scenarios []ScenarioConfig `json:"scenarios"`
}

// ScenarioConfig represents the configuration for a single scenario
type ScenarioConfig struct {
	Name        string `json:"name"`
	Concurrency int    `json:"concurrency"`
	Duration    string `json:"duration"`
	Rate        int    `json:"rate"`
}
