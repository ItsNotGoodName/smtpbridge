package config

type Endpoint struct {
	Name               string            `json:"name" mapstructure:"name"`
	TextDisable        bool              `json:"text_disable" mapstructure:"text_disable"`
	TextTemplate       string            `json:"text_template" mapstructure:"text_template"`
	AttachmentsDisable bool              `json:"attachments_disable" mapstructure:"attachments_disable"`
	Type               string            `json:"type" mapstructure:"type"`
	Config             map[string]string `json:"config" mapstructure:"config"`
}
