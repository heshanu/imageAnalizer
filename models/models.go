package models

type TextGenInput struct {
	Prompt     string  `form:"prompt" binding:"required"`
	Model      string  `form:"model"`
	MaxTokens  int     `form:"max_tokens"`
	Temperature float64 `form:"temperature"`
}