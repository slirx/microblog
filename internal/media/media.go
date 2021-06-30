package media

type ImagesResponse struct {
	Images map[string]string `json:"images"`
}

type UploadImageResponse struct {
	URL string `json:"url"`
}
