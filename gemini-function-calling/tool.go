package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

var fileWriteSchema = &genai.Schema{
	Type: genai.TypeObject,
	Properties: map[string]*genai.Schema{
		"fileName": {
			Type:        genai.TypeString,
			Description: "The name of the file to write to. Do not include extension, it will be automatically added (.txt)",
		},
		"content": {
			Type:        genai.TypeString,
			Description: "The text content to write to the file",
		},
	},
	Required: []string{"fileName", "content"},
}

var FileTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		{
			Name:        "file_write",
			Description: "write a text file to user local file system with specified name and content.",
			Parameters:  fileWriteSchema,
		},
	},
}

func WriteDesktop(fileName string, content string) error {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create path to results directory relative to current directory
	resultsDir := filepath.Join(currentDir, "results")

	// Create the results directory if it doesn't exist
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		return err
	}

	// Add .txt extension and create full path
	fileName = fileName + ".txt"
	fullPath := filepath.Join(resultsDir, fileName)

	// Format content
	formattedContent := strings.ReplaceAll(content, "\\n", "\n")

	// Write the file
	err = os.WriteFile(fullPath, []byte(formattedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}
