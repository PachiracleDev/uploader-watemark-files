package dtos

type (
	UploadFile struct {
		Privacy  string `json:"privacy" validate:"required"`
		Username string `json:"username"`
	}

	UploadCollectionFile struct {
		Username     string `json:"username"`
		TokenUpload  string `json:"tokenUpload"`
		CollectionId string `json:"collectionId"`
	}
)
