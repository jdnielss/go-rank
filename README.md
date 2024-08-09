# Code Runner API

This project provides an API to run code in various languages (JavaScript, Go, Java, and Python) based on provided unit tests. It leverages Docker to create a secure and isolated environment for executing the code.

## Table of Contents

- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Building the Docker Image](#building-the-docker-image)
- [Running the Docker Container](#running-the-docker-container)
- [API Endpoints](#api-endpoints)
- [Example Requests](#example-requests)

## Getting Started

These instructions will help you set up the project on your local machine.

### Prerequisites

- Docker
- Go (if you plan to run the application outside of Docker)

### Building the Docker Image

To build the Docker image, run the following command:

```sh
docker build -t code-runner .
```

Run the Docker container:
```shell
docker run -p 8080:8080 code-runner
```
This will start the server on port 8080.

---

# API Endpoints

## POST /run

Executes the provided code and runs it against the provided test cases.

Request Body

```json
{
  "language": "javascript",
  "code": "console.log('Hello, World!');",
  "testCases": [
    {
      "input": "",
      "expected": "Hello, World!"
    },
    {
      "input": "Another test case",
      "expected": "Another test case"
    }
  ]
}
```
- language: The programming language of the code (supported: javascript, golang, java, python)
- code: The code to execute
- testCases: An array of test cases, each containing:
- input: The input to pass to the code
- expected: The expected output

Response Body
```json
[
    {
        "passed": true,
        "output": "Hello, World!",
        "error": ""
    },
    {
        "passed": true,
        "output": "Another test case",
        "error": ""
    }
]
```

- passed: Indicates whether the test case passed 
- output: The actual output of the code 
- error: Any error encountered during execution

# Example Requests

## JavaScript Example

```shell
curl -X POST http://localhost:8080/run \
     -H "Content-Type: application/json" \
     -d '{
         "language": "javascript",
         "code": "console.log(\"Hello, World!\");",
         "testCases": [
             {
                 "input": "",
                 "expected": "Hello, World!"
             }
         ]
     }'
```

## Go Example

```shell
curl -X POST http://localhost:8080/run \
     -H "Content-Type: application/json" \
     -d '{
         "language": "golang",
         "code": "package main\nimport \"fmt\"\nfunc main() { var input string; fmt.Scanln(&input); fmt.Println(input) }",
         "testCases": [
             {
                 "input": "Hello, World!",
                 "expected": "Hello, World!"
             }
         ]
     }'
```

## Java Example

```shell
curl -X POST http://localhost:8080/run \
     -H "Content-Type: application/json" \
     -d '{
         "language": "java",
         "code": "import java.util.Scanner;\n\npublic class Main {\n    public static void main(String[] args) {\n        Scanner scanner = new Scanner(System.in);\n        String input = scanner.nextLine();\n        System.out.println(input);\n    }\n}",
         "testCases": [
             {
                 "input": "Hello, World!",
                 "expected": "Hello, World!"
             }
         ]
     }'
```

## Python Example

```shell
curl -X POST http://localhost:8080/run \
     -H "Content-Type: application/json" \
     -d '{
         "language": "python",
         "code": "print(input())",
         "testCases": [
             {
                 "input": "Hello, World!",
                 "expected": "Hello, World!"
             }
         ]
     }'
```

# License

This project is licensed under the MIT License.