package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

// var testDataDir = "path/to/test/data"

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func runDocumentUnderstanding() {
	ctx := context.Background()
	// Access your API key as an environment variable
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	file, err := client.UploadFileFromPath(ctx, "ML_BOOK.pdf", nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer client.DeleteFile(ctx, file.Name)

	resp, err := model.GenerateContent(ctx, genai.Text("What is in page 457 of this document?"), genai.FileData{URI: file.URI})
	if err != nil {
		log.Println(err)
		return
	}

	printResponse(resp)
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}

func main() {
	runDocumentUnderstanding()
}
