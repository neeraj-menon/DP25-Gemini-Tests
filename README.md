# Gemini API Exploration Projects

This repository contains a collection of projects exploring different capabilities of Google's Gemini AI API using Go. Each project demonstrates specific features and use cases of the Gemini API.

## Projects Overview

### 1. Gemini Chat (`gemini-chat/`)
A basic chat application that demonstrates:
- Direct interaction with Gemini AI
- Basic conversation handling
- Message streaming
- Simple prompt engineering
- Error handling and response processing

### 2. Gemini Context Caching (`gemini-context-caching/`)
Explores context management in conversations:
- Conversation history caching
- Context window management
- Memory optimization techniques
- Stateful conversations
- Session management

### 3. Gemini Document Understanding (`gemini-doc-understanding/`)
Demonstrates Gemini's document processing capabilities:
- Text document analysis
- Content extraction
- Document summarization
- Information retrieval
- Natural language understanding

### 4. Gemini Function Calling (`gemini-function-calling/`)
Showcases Gemini's function calling capabilities:
- Natural language to function mapping
- File creation based on user commands
- Structured data extraction
- Custom tool implementation
- Error handling and response management

## Common Features Across Projects

- Environment variable management using `.env` files
- Go module dependency management
- Error handling and logging
- Clean code architecture
- Documentation and examples

## Prerequisites

For all projects:
- Go 1.22 or higher
- Google Gemini API Key
- Basic understanding of Go programming
- Familiarity with API concepts

## Setup Instructions

1. Clone the repository:
```bash
git clone <repository-url>
cd DP25-Gemini-Tests
```

2. Set up environment variables:
Create a `.env` file in each project directory with:
```env
GENAI_API_KEY=your_api_key_here
```

3. Install dependencies:
```bash
cd <project-directory>
go mod tidy
```

## Project Structure

```
DP25-Gemini-Tests/
├── gemini-chat/               # Basic chat implementation
├── gemini-context-caching/    # Context management demo
├── gemini-doc-understanding/  # Document processing features
├── gemini-function-calling/   # Function calling capabilities
├── .gitignore                # Git ignore rules
└── README.md                 # This file
```

## Security

- All projects use environment variables for sensitive data
- API keys are never committed to the repository
- Proper error handling for security-related issues

## Contributing

Feel free to contribute to any of the projects:
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request


