package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

//go:embed .env
var env embed.FS

func main() {
	err := loadEnv()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.POST("/", handleRequest)

	port := ":8000"
	log.Printf("Listening on port %s\n", port)

	router.Run(port)
}

func loadEnv() error {
	content, err := env.ReadFile(".env")
	if err != nil {
		return err
	}

	err = os.WriteFile(".env", content, 0755)
	if err != nil {
		return err
	}

	err = godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	return nil
}

func handleRequest(ctx *gin.Context) {
	var request TranslationRequest

	err := ctx.ShouldBind(&request)
	if err != nil {
		log.Printf("ShouldBind error: %v", err)
	}

	var fromLanguage, toLanguage, translate_text string
	if isEnglish(request.Text) {
		fromLanguage = request.Destination[1]
		toLanguage = request.Destination[0]
	} else {
		fromLanguage = request.Destination[0]
		toLanguage = request.Destination[1]
	}

	translate_text = fmt.Sprintf("请将下面的文本翻译为%s: %s", toLanguage, request.Text)
	translate_result := translate(translate_text)

	ctx.JSON(200, gin.H{
		"text":   request.Text,
		"from":   fromLanguage,
		"to":     toLanguage,
		"result": []string{translate_result},
	})
}

func translate(content string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Printf("genai.NewClient error: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(content))
	if err != nil {
		log.Printf("model.GenerateContent error: %v", err)
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

func isEnglish(s string) bool {
	var englishCount, chineseCount int
	for _, r := range s {
		if unicode.Is(unicode.Han, r) {
			chineseCount++
		} else if unicode.Is(unicode.Latin, r) {
			englishCount++
		}
	}

	if englishCount > chineseCount {
		return true
	} else {
		return false
	}
}
