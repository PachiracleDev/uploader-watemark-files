package dtos

type (
	UploadFile struct {
		Privacy  string `json:"privacy" validate:"required"`
		Username string `json:"username"`
	}
)
