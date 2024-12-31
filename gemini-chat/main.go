package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

func printResponse(res *genai.GenerateContentResponse) {
	if res == nil || len(res.Candidates) == 0 {
		fmt.Println("No response from the model.")
		return
	}
	// Print the text of the first candidate response
	if text, ok := res.Candidates[0].Content.Parts[0].(genai.Text); ok {
		fmt.Printf("Response: %s\n", string(text))
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down...")
		cancel()
	}()

	// Create a new Gemini client
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}

	// Upload the PDF file (if needed)
	var file *genai.File
	if _, err := os.Stat("ML_BOOK_250.pdf"); !os.IsNotExist(err) {
		file, err = client.UploadFileFromPath(ctx, "ML_BOOK_250.pdf", nil)
		if err != nil {
			log.Println("Error uploading PDF:", err)
			return
		}
		defer client.DeleteFile(ctx, file.Name)
	}

	// Create cached content (if needed)
	var cachedContent *genai.CachedContent
	if file != nil {
		fd := genai.FileData{URI: file.URI}
		// Adjusted System Instruction:
		adjustedSystemInstruction := "You are an expert in machine learning. Analyze the content of the provided PDF document. I will ask you questions about specific topics or pages within this document. Use the information from the PDF to answer my questions accurately."

		adjustedContents := []*genai.Content{genai.NewUserContent(fd)}
		// Consider adding more content here to reach the minimum token count

		argcc := &genai.CachedContent{
			Model:             "gemini-1.5-flash-001",
			SystemInstruction: genai.NewUserContent(genai.Text(adjustedSystemInstruction)),
			Contents:          adjustedContents,
		}
		cachedContent, err = client.CreateCachedContent(ctx, argcc)
		if err != nil {
			log.Fatalf("Failed to create cached content: %v", err)
		}
		defer client.DeleteCachedContent(ctx, cachedContent.Name)
	}

	// Initialize chat model
	var model *genai.GenerativeModel
	if cachedContent != nil {
		model = client.GenerativeModelFromCachedContent(cachedContent)
		log.Println("Using model with cached content.")
	} else {
		model = client.GenerativeModel("gemini-1.5-flash")
		log.Println("Using standard model without cached content.")
	}
	
	cs := model.StartChat()

	// Set initial chat history
	cs.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text("Hello, I have loaded a machine learning document. Please help me understand its contents."),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text("I'll help you understand the machine learning document. What would you like to know about it?"),
			},
			Role: "model",
		},
	}

	// Start the chat and read user input
	fmt.Println("Chat started. Type your questions about the document (press Ctrl+C to exit)")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		userInput := strings.TrimSpace(scanner.Text())
		if userInput == "" {
			continue // Skip empty lines
		}

		// Provide specific instructions in the query
		query := fmt.Sprintf("Using the PDF you were provided, %s", userInput)
		res, err := cs.SendMessage(ctx, genai.Text(query))
		if err != nil {
			log.Fatalf("Error sending message: %v", err)
		}

		// Print the response content
		printResponse(res)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
