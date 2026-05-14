package dto

type GenerateImageRequest struct {
	Prompt  string `json:"prompt" binding:"required" example:"A friendly robot reading a book, watercolor style"`
	Model   string `json:"model,omitempty" example:"gpt-image-1"`
	Size    string `json:"size,omitempty" example:"1024x1024"`
	Quality string `json:"quality,omitempty" example:"auto"`
}

type GenerateImageResponse struct {
	B64JSON string `json:"b64_json,omitempty"`
	URL     string `json:"url,omitempty"`
	DataURL string `json:"data_url,omitempty"`
}
