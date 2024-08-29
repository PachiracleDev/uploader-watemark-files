package controllers

import (
	controllersModules "uploader/internal/presenters/controllers/upload"

	"go.uber.org/fx"
)

var RegisterControllers = fx.Options(
	controllersModules.UploadControllerModule,
)
