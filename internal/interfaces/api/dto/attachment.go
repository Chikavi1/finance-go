package dto

type AttachmentResponse struct {
	ID            string  `json:"id"`
	TransactionID *string `json:"transaction_id,omitempty"`
	Filename      string  `json:"filename"`
	OriginalName  string  `json:"original_name"`
	MimeType      string  `json:"mime_type"`
	Size          int64   `json:"size"`
	URL           string  `json:"url"`
	CreatedAt     string  `json:"created_at"`
}
