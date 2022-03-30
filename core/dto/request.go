package dto

type MessageHandleRequest struct {
	Subject     string
	From        string
	To          map[string]struct{}
	Text        string
	Attachments []AttachmentHandleRequest
}

type AttachmentHandleRequest struct {
	Name string
	Data []byte
}

func (r *MessageHandleRequest) AddAttachment(name string, data []byte) {
	r.Attachments = append(r.Attachments, AttachmentHandleRequest{name, data})
}

type MessageGetRequest struct {
	ID int64
}

type MessageListRequest struct {
	Cursor    int64
	Ascending bool
	Limit     int
}

type MessageListResponse struct {
	Messages   []Message `json:"messages"`
	NextCursor int64     `json:"next_cursor"`
	HasMore    bool      `json:"has_more"`
}

type InfoResponse struct {
	EventsCount      int   `json:"events_count"`
	MessagesCount    int   `json:"messages_count"`
	AttachmentsCount int   `json:"attachments_count"`
	AttachmentsSize  int64 `json:"attachments_size"`
}

type VersionResponse struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
	BuiltBy string `json:"built_by"`
}

type EventListRequest struct {
	Ascending bool
	Page      int
	Limit     int
	EntityID  int64
}

type EventListResponse struct {
	Events   []Event `json:"events"`
	MaxPages int     `json:"max_pages"`
}

type MessageDeleteRequest struct {
	ID int64
}

type SMTPLoginRequest struct {
	Username string
	Password string
}
