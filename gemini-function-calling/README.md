# Gemini File Writing Chatbot

A Go-based chatbot that uses Google's Gemini AI to understand natural language requests and create text files. The chatbot leverages Gemini's function calling capabilities to process user requests and create files in a specified directory.

## Features

- Natural language processing for file creation commands
- Automatic file extension handling (.txt)
- Organized file storage in a dedicated results directory
- Support for multi-line content using \n notation
- Robust error handling and user feedback

## Prerequisites

- Go 1.22 or higher
- Google Gemini API Key
- Environment variables setup

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd gemini-function-calling
```

2. Install dependencies:
```bash
go mod tidy
```

3. Create a `.env` file in the project root:
```env
GENAI_API_KEY=your_api_key_here
```

## Project Structure

```
gemini-function-calling/
├── main.go         # Main application and chat logic
├── tool.go         # File operations and function definitions
├── results/        # Output directory for generated files
├── go.mod          # Go module file
├── go.sum          # Go module checksum
└── .env            # Environment variables
```

## Application Flow and Function Calling Method

### 1. Initialization
```go
// Initialize Gemini client and model
genaiApp := &App{}
genaiApp.client, err = NewClient(apiKey, context.Background())
genaiApp.model = NewModel(genaiApp.client, GenaiModel)
genaiApp.model.Tools = []*genai.Tool{FileTool}  // Register our file writing tool
genaiApp.cs = genaiApp.model.StartChat()        // Start chat session
```

### 2. Tool Definition (tool.go)
```go
// Define schema for file writing function
var fileWriteSchema = &genai.Schema{
    Type: genai.TypeObject,
    Properties: map[string]*genai.Schema{
        "fileName": {
            Type:        genai.TypeString,
            Description: "The name of the file to write to",
        },
        "content": {
            Type:        genai.TypeString,
            Description: "The text content to write to the file",
        },
    },
    Required: []string{"fileName", "content"},
}

// Register the tool with Gemini
var FileTool = &genai.Tool{
    FunctionDeclarations: []*genai.FunctionDeclaration{
        {
            Name:        "file_write",
            Description: "write a text file to user local file system",
            Parameters:  fileWriteSchema,
        },
    },
}
```

### 3. Function Call Flow

1. **User Input Processing**
   ```go
   input, _ := reader.ReadString('\n')
   response, err := genaiApp.cs.SendMessage(context.Background(), genai.Text(input))
   ```

2. **Response Handling and Function Mapping**
   ```go
   // Extract function calls from response
   for _, part := range resp.Candidates[0].Content.Parts {
       functionCall, ok := part.(genai.FunctionCall)
       if ok {
           // Process function call
           switch functionCall.Name {
           case "file_write":
               // Extract parameters
               fileName, fileNameOk := functionCall.Args["fileName"].(string)
               content, contentOk := functionCall.Args["content"].(string)
               
               // Validate parameters
               if !fileNameOk || fileName == "" {
                   funcResponse["error"] = "expected non-empty string at fileName"
                   break
               }
               if !contentOk || content == "" {
                   funcResponse["error"] = "expected non-empty string at content"
                   break
               }
               
               // Programmatically call WriteDesktop with extracted parameters
               err := WriteDesktop(fileName, content)
               if err != nil {
                   funcResponse["error"] = "could not write file"
               } else {
                   funcResponse["result"] = "file successfully written"
               }
           }
       }
   }
   ```

3. **File Writing Implementation**
   ```go
   func WriteDesktop(fileName string, content string) error {
       currentDir, _ := os.Getwd()
       resultsDir := filepath.Join(currentDir, "results")
       fileName = fileName + ".txt"
       fullPath := filepath.Join(resultsDir, fileName)
       return os.WriteFile(fullPath, []byte(content), 0644)
   }
   ```

### 4. Complete Flow Example

1. User sends: "Create a file named 'notes' with content 'Hello, World!'"
2. Gemini processes the request and identifies it as a file writing operation
3. Gemini generates a function call with the "file_write" function name:
   ```json
   {
       "name": "file_write",
       "args": {
           "fileName": "notes",
           "content": "Hello, World!"
       }
   }
   ```
4. Application receives the function call and:
   - Identifies the function name as "file_write"
   - Extracts fileName and content from the args
   - **Programmatically calls WriteDesktop(fileName, content)**
   - Returns success/failure response
5. Gemini generates a natural language response to the user
6. User receives confirmation of file creation

The key aspect here is that the application acts as a bridge between Gemini's function calls and actual Go functions. When Gemini returns a "file_write" function call, the application automatically maps this to the WriteDesktop function, passing the extracted parameters. This mapping is done in the buildResponse function, which checks the function name and routes the call to the appropriate implementation.

## Error Handling

The application handles various scenarios:
- Invalid or missing API key
- Malformed file names or content
- File system permission issues
- Directory access problems
- Invalid function calls from Gemini

## Technical Implementation

### Function Call Format
```json
{
    "name": "file_write",
    "args": {
        "fileName": "example",
        "content": "Sample content"
    }
}
```

### Response Format
```json
{
    "name": "Function_Call",
    "response": {
        "result": "file successfully written"
    }
}
```

## Future Enhancements

1. Support for different file formats
2. File reading capabilities
3. File modification operations
4. Directory organization features
5. Content formatting options

## How Function Calling Works

The project demonstrates function calling with Gemini AI through three main components:

1. **Schema Definition** (tool.go)
```go
var fileWriteSchema = &genai.Schema{
    Type: genai.TypeObject,
    Properties: map[string]*genai.Schema{
        "fileName": {
            Type:        genai.TypeString,
            Description: "The name of the file to write to",
        },
        "content": {
            Type:        genai.TypeString,
            Description: "The text content to write to the file",
        },
    },
    Required: []string{"fileName", "content"},
}
```

2. **Function Registration** (tool.go)
```go
var FileTool = &genai.Tool{
    FunctionDeclarations: []*genai.FunctionDeclaration{
        {
            Name:        "file_write",
            Description: "write a text file to user local file system",
            Parameters:  fileWriteSchema,
        },
    },
}
```

3. **File Writing Implementation** (tool.go)
```go
func WriteDesktop(fileName string, content string) error {
    // Gets current working directory
    // Creates file in results directory
    // Handles file writing operations
}
```

### Function Call Flow

1. User provides a natural language request
2. Request is processed by Gemini AI
3. Gemini identifies the intent and extracts parameters
4. Function call is made with fileName and content
5. File is created in the results directory
6. Response is sent back to user

## Usage Examples

```
You: Create a file named 'meeting-notes' with content 'Team meeting at 2 PM\nDiscuss Q1 goals'
Bot: I'll create a file named 'meeting-notes.txt' with your meeting notes
File 'meeting-notes.txt' created successfully in results directory

You: Write a file called 'shopping' with 'Milk\nBread\nEggs'
Bot: I'll create a file named 'shopping.txt' with your shopping list
File 'shopping.txt' created successfully in results directory
