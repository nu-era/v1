package alerts

// Alert represents the message that the client will receive and send
// back. Contains important information for alert as well as a 'receipt'
// to confirm delivery
type Alert struct {
	ID          int64  `json:"id"`
	Message     string `json:"message"`
	DeviceName  string `json:"deviceName"`
	Status      bool   `json:"status"`
	CreatedAt   string `json:"created_at"`
	EditedAt    string `json:"edited_at"`
	SendTime    string `json:"sendTime"`
	ReceiveTime string `json:"receiveTime"`
}

//AlertUpdates represents allowed updates to an alert
type AlertUpdates struct {
	Status  string `json:"status"`
	Message string `json:"status"`
}
