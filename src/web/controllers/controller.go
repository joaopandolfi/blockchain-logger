package controllers

import (
	"logger/models"
	"logger/web/server"
)

var SystemPermissions = []string{models.PermissionSystem}

// Controller public contract
type Controller interface {
	SetupRouter(s *server.Server)
}
