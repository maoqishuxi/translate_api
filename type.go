package main

// 上传类型
type TranslationRequest struct {
	Name        string   `json:"name"`
	Text        string   `json:"text"`
	Destination []string `json:"destination"`
	Source      string   `json:"source"`
}
