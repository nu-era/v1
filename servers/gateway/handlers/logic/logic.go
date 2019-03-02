package handlers

// Message represents the message that the client will receive and send
// back. Contains important information for alert as well as a 'receipt'
// to confirm delivery
type Message struct {
	ID        int64  `json:"id"`
	DeviceID  int64  `json:"device_id"`
	Message   string `json:"message"`
	CreatedAt string `json:"created_at"`
	EditedAt  string `json:"edited_at"`
	Status    bool   `json:"status"`
}
