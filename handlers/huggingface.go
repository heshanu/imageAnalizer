package handlers

import (
	"bytes"
	"context"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/heshanu/go/models"
)

const huggingFaceAPI = "https://api-inference.huggingface.co/models/"

func GenerateTextHandler(c *gin.Context) {
	// Parse input
	var input models.TextGenInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Set default model
	if input.Model == "" {
		input.Model = "gpt2"
	}
	
	// Prepare payload
	payload := map[string]interface{}{
		"inputs": input.Prompt,
		"parameters": map[string]interface{}{
			"max_new_tokens":  input.MaxTokens,
			"temperature":     input.Temperature,
			"return_full_text": false,
		},
	}
	
	// Call Hugging Face API
	result, err := huggingFaceRequest(c.Request.Context(), input.Model, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Parse response
	var textResults []struct {
		GeneratedText string `json:"generated_text"`
	}
	if err := json.Unmarshal(result, &textResults); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}
	
	if len(textResults) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Empty response from API"})
		return
	}
	
	// Return to template
	c.HTML(http.StatusOK, "index.html", gin.H{
		"prompt": input.Prompt,
		"result": textResults[0].GeneratedText,
	})
}

func ClassifyImageHandler(c *gin.Context) {
	// Parse input
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file required"})
		return
	}
	defer file.Close()
	
	model := c.PostForm("model")
	if model == "" {
		model = "google/vit-base-patch16-224"
	}
	
	// Read image data
	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image"})
		return
	}
	
	// Call Hugging Face API
	result, err := huggingFaceImageRequest(c.Request.Context(), model, imageData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Parse response
	var classifications []struct {
		Label string  `json:"label"`
		Score float64 `json:"score"`
	}
	if err := json.Unmarshal(result, &classifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}
	
	// Return to template
	c.HTML(http.StatusOK, "index.html", gin.H{
		"image":        header.Filename,
		"classResults": classifications,
	})
}

func huggingFaceRequest(ctx context.Context, model string, payload interface{}) ([]byte, error) {
	apiKey := os.Getenv("HUGGING_FACE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not configured")
	}
	
	url := huggingFaceAPI + model
	jsonData, _ := json.Marshal(payload)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	
	return io.ReadAll(resp.Body)
}

func huggingFaceImageRequest(ctx context.Context, model string, imageData []byte) ([]byte, error) {
	apiKey := os.Getenv("HUGGING_FACE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key not configured")
	}
	
	url := huggingFaceAPI + model
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(imageData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/octet-stream")
	
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}
	
	return io.ReadAll(resp.Body)
}