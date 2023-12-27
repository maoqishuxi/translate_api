package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", handleRequest)

	router.Run()
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func handleRequest(ctx *gin.Context) {
	var request TranslationRequest

	err := ctx.ShouldBind(&request)
	if err != nil {
		log.Fatal(err)
	}

	translate_text := fmt.Sprintf("请将下面的文本翻译为%s: %s", request.Destination[0], request.Text)

	translate_result := translate(translate_text)

	ctx.JSON(200, gin.H{
		"text":   request.Text,
		"from":   request.Destination[1],
		"to":     request.Destination[0],
		"result": []string{translate_result},
	})
}

func translate(content string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(content))
	if err != nil {
		log.Fatal(err)
	}

	ret := printResponse(resp)
	return ret
}

func printResponse(resp *genai.GenerateContentResponse) string {
	result := ""
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				result += fmt.Sprintf("%s", part)
			}
		}
	}

	return result
}
