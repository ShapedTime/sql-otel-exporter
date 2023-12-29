package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

// Query represents the structure of each query in the YAML configuration.
type Query struct {
	SQL                  string            `yaml:"sql"`
	TrackingColumn       string            `yaml:"tracking_column,omitempty"`
	DefaultTrackingValue string            `yaml:"default_tracking_value,omitempty"`
	ValueColumn          string            `yaml:"value_column"`
	AttributeColumns     []string          `yaml:"attribute_columns,omitempty"`
	DataType             string            `yaml:"data_type"`
	ValueType            string            `yaml:"value_type"`
	Description          string            `yaml:"description,omitempty"`
	StaticAttributes     map[string]string `yaml:"static_attributes,omitempty"`
	IntervalSeconds      int               `yaml:"interval_seconds"`
	MetricName           string            `yaml:"metric_name"`
	Name                 string            `yaml:"name"`
	Buckets              []float64         `yaml:"buckets,omitempty"`
	IgnoreError          bool              `yaml:"ignore_error,omitempty"`
}

// Config represents the structure of the YAML configuration file.
type Config struct {
	Datasource string  `yaml:"datasource"`
	Driver     string  `yaml:"driver"`
	Queries    []Query `yaml:"queries"`
}

// ReadConfig reads a YAML configuration file and unmarshals it into a Config struct.
func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
