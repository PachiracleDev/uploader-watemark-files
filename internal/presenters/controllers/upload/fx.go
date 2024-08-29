package controllers

import (
	"go.uber.org/fx"
)

var UploadControllerModule = fx.Module(
	"upload_controller",
	//USECASES
	fx.Invoke(UploadController),
)
