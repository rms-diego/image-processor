package validations

type Image struct {
	ID        string `db:"id" json:"id"`
	URL       string `db:"url" json:"url"`
	CreatedAt string `db:"created_at" json:"created_at,omitempty"`
	S3Key     string `db:"s3_key" json:"s3_key,omitempty"`
	UserID    string `db:"user_id" json:"user_id,omitempty"`
}

type ManyImages = []Image

type ListImagesResponse struct {
	TotalImages int        `json:"total_images"`
	Data        ManyImages `json:"data"`
}

type TransformImageReqBody struct {
	Resize  *resizeOptions `json:"resize,omitempty"`
	Crop    *cropOptions   `json:"crop,omitempty"`
	Rotate  *int           `json:"rotate,omitempty"`
	Format  *string        `json:"format,omitempty"`
	Filters *filterOptions `json:"filters,omitempty"`
	Quality *int           `json:"quality,omitempty"`
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
	S3Key   string                `json:"s3_key"`
	Payload TransformImageReqBody `json:"payload"`
}
