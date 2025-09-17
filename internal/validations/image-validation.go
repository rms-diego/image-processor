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
