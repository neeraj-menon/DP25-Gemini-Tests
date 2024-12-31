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

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()
	// Access your API key as an environment variable
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	file, err := client.UploadFileFromPath(ctx, "ML_BOOK_250.pdf", nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer client.DeleteFile(ctx, file.Name)

	fd := genai.FileData{URI: file.URI}

	argcc := &genai.CachedContent{
		Model:             "gemini-1.5-flash-001",
		SystemInstruction: genai.NewUserContent(genai.Text("You are an expert analyzing transcripts.")),
		Contents:          []*genai.Content{genai.NewUserContent(fd)},
	}
	cc, err := client.CreateCachedContent(ctx, argcc)
	if err != nil {
		log.Fatal(err)
	}
	defer client.DeleteCachedContent(ctx, cc.Name)

	// Create the request.
	req := []genai.Part{
		genai.Text("Please summarize this transcript"),
	}

	model := client.GenerativeModelFromCachedContent(cc)

	// Generate content.
	resp, err := model.GenerateContent(ctx, req...)
	if err != nil {
		panic(err)
	}

	// Handle the response of generated text.
	for _, c := range resp.Candidates {
		if c.Content != nil {
			fmt.Println(*c.Content)
		}
	}
}
