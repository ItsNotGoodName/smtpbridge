package config

type Bridge struct {
	From          string   `json:"from" mapstructure:"from"`
	To            string   `json:"to" mapstructure:"to"`
	FromRegex     string   `json:"from_regex" mapstructure:"from_regex"`
	ToRegex       string   `json:"to_regex" mapstructure:"to_regex"`
	MatchTemplate string   `json:"match_template" mapstructure:"match_template"`
	Endpoints     []string `json:"endpoints" mapstructure:"endpoints"`
}
