package main

import (
	"uploader/config"
	"uploader/pkg/aws"
	"uploader/pkg/http"
	"uploader/pkg/validator"

	ctrls "uploader/internal/presenters/controllers"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		http.HttpModule,
		validator.Module,
		aws.AwsModule,
		ctrls.RegisterControllers,
		fx.Invoke(http.RunHttpServer),
	).Run()
}
