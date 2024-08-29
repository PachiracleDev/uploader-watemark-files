package aws

import "go.uber.org/fx"

var AwsModule = fx.Module("aws_sdk", fx.Provide(NewSDKImplementation))
