package config

type Bridge struct {
	Name      string   `json:"name" mapstructure:"name"`
	From      string   `json:"from" mapstructure:"from"`
	To        string   `json:"to" mapstructure:"to"`
	FromRegex string   `json:"from_regex" mapstructure:"from_regex"`
	ToRegex   string   `json:"to_regex" mapstructure:"to_regex"`
	Endpoints []string `json:"endpoints" mapstructure:"endpoints"`
}
