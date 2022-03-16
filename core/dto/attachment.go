package dto

type Attachment struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	File string `json:"file"`
	Type string `json:"type"`
}
