package dto

type Message struct {
	UUID        string       `json:"uuid"`
	Status      string       `json:"status"`
	From        string       `json:"from"`
	To          []string     `json:"to"`
	Subject     string       `json:"subject"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	Path string `json:"path"`
}
