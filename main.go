package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type RunRequest struct {
	Language  string     `json:"language"`
	Code      string     `json:"code"`
	TestCases []TestCase `json:"testCases"`
}

type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

type RunResponse struct {
	Passed bool   `json:"passed"`
	Output string `json:"output"`
	Error  string `json:"error"`
}

func main() {
	app := fiber.New()

	app.Post("/run", runCode)

	log.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func runCode(c *fiber.Ctx) error {
	var request RunRequest
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	results := make([]RunResponse, len(request.TestCases))
	for i, testCase := range request.TestCases {
		output, err := executeCode(request.Language, request.Code, testCase.Input)
		if err != nil {
			log.Printf("Error executing code: %v", err)
			results[i] = RunResponse{Passed: false, Output: "", Error: err.Error()}
			continue
		}

		passed := strings.TrimSpace(output) == testCase.Expected
		results[i] = RunResponse{Passed: passed, Output: output, Error: ""}
	}

	return c.JSON(results)
}

func executeCode(language, code, input string) (string, error) {
	var cmd *exec.Cmd
	switch language {
	case "javascript":
		log.Println("Executing JavaScript code")
		cmd = exec.Command("node", "-e", code)
	case "golang":
		log.Println("Executing Go code")
		tempFile, err := os.CreateTemp("", "*.go")
		if err != nil {
			log.Printf("Error creating temp file: %v", err)
			return "", err
		}
		defer os.Remove(tempFile.Name())

		if _, err := tempFile.WriteString(code); err != nil {
			log.Printf("Error writing to temp file: %v", err)
			return "", err
		}
		if err := tempFile.Close(); err != nil {
			log.Printf("Error closing temp file: %v", err)
			return "", err
		}

		cmd = exec.Command("go", "run", tempFile.Name())
	case "java":
		log.Println("Executing Java code")

		// Create a temporary directory
		tempDir, err := os.MkdirTemp("", "java")
		if err != nil {
			log.Printf("Error creating temp directory: %v", err)
			return "", err
		}
		defer os.RemoveAll(tempDir)

		// Create a file named Main.java in the temporary directory
		tempFile := filepath.Join(tempDir, "Main.java")
		if err := os.WriteFile(tempFile, []byte(code), 0644); err != nil {
			log.Printf("Error writing to temp file: %v", err)
			return "", err
		}

		// Compile the Java code
		compileCmd := exec.Command("javac", tempFile)
		var compileOut, compileErr bytes.Buffer
		compileCmd.Stdout = &compileOut
		compileCmd.Stderr = &compileErr
		log.Printf("Running command: %v", compileCmd)
		if err := compileCmd.Run(); err != nil {
			log.Printf("Error compiling Java code: %v\n%s", err, compileErr.String())
			return compileErr.String(), err
		}

		// Run the compiled Java code
		cmd = exec.Command("java", "-cp", tempDir, "Main")
	case "python":
		log.Println("Executing Python code")
		tempFile, err := os.CreateTemp("", "*.py")
		if err != nil {
			log.Printf("Error creating temp file: %v", err)
			return "", err
		}
		defer os.Remove(tempFile.Name())

		if _, err := tempFile.WriteString(code); err != nil {
			log.Printf("Error writing to temp file: %v", err)
			return "", err
		}
		if err := tempFile.Close(); err != nil {
			log.Printf("Error closing temp file: %v", err)
			return "", err
		}

		cmd = exec.Command("python3", tempFile.Name())
	default:
		err := errors.New("unsupported language")
		log.Printf("Error language: %v", err)
		return "", err
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = bytes.NewBufferString(input)

	log.Printf("Running command: %v", cmd)
	log.Printf("PATH: %v", os.Getenv("PATH"))
	if err := cmd.Run(); err != nil {
		log.Printf("Command error: %v\n%s", err, stderr.String())
		return stderr.String(), err
	}

	return stdout.String(), nil
}
