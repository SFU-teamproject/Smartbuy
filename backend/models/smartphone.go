package models

type Smartphone struct {
	ID           int      `json:"id"`
	Model        string   `json:"model"`
	Producer     string   `json:"producer"`
	Memory       int      `json:"memory"`
	Ram          int      `json:"ram"`
	DisplaySize  float32  `json:"display_size"`
	Price        int      `json:"price"`
	RatingsSum   int      `json:"ratings_sum"`
	RatingsCount int      `json:"ratings_count"`
	ImagePath    string   `json:"image_path"`
	Description  string   `json:"description"`
	Reviews      []Review `json:"reviews,omitempty"`
}
