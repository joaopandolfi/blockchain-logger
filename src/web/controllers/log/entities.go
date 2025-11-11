package log

import (
	"logger/models"
	"strings"
)

type payload struct {
	Data     map[string]interface{}
	SystemID string
	Tags     []string
}

func (p *payload) ToLog() *models.Log {
	return &models.Log{
		SystemID: p.SystemID,
		Payload:  p.Data,
		Tags:     strings.Join(p.Tags, models.TAG_SEPARATOR),
	}
}
