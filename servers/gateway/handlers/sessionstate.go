package handlers

import (
	"time"

	"github.com/New-Era/servers/gateway/models/devices"
)

//SessionState of start time and user
type SessionState struct {
	StartTime time.Time       `json:"startTime"`
	Device    *devices.Device `json:"device"`
}
