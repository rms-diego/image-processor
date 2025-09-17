package validations

type Image struct {
	ID        string `db:"id" json:"id"`
	URL       string `db:"url" json:"url"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type ImagesFound = []Image

type ListImagesResponse struct {
	TotalImages int         `json:"total_images"`
	Data        ImagesFound `json:"data"`
}

type TransformImageReqBody struct {
	Resize  *resizeOptions `json:"resize,omitempty"`
	Crop    *cropOptions   `json:"crop,omitempty"`
	Rotate  *int           `json:"rotate,omitempty"`
	Format  *string        `json:"format,omitempty"`
	Filters *filterOptions `json:"filters,omitempty"`
}

type resizeOptions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type cropOptions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

type filterOptions struct {
	Grayscale bool `json:"grayscale"`
	Sepia     bool `json:"sepia"`
}

type TransformMessageQueue struct {
	ImageID string                `json:"image_id"`
	Payload TransformImageReqBody `json:"payload"`
}
