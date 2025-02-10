package main

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("templates/*.html")

	// Serve static files (CSS)
	r.Static("/static", "./static")

	// Handle GET request
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{"output": "", "code": ""})
	})

	// Handle POST request (Code Execution)
	r.POST("/", func(c *gin.Context) {
		code := c.PostForm("code")

		// Save code to a temporary Go file
		tmpFile := "temp.go"
		os.WriteFile(tmpFile, []byte(code), 0644)

		// Execute the Go file
		cmd := exec.Command("go", "run", tmpFile)
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()

		// Delete temp file after execution
		os.Remove(tmpFile)

		output := out.String()
		if err != nil {
			output += "\nError: " + err.Error()
		}

		// Return output to frontend
		c.HTML(200, "index.html", gin.H{"output": output, "code": code})
	})

	// Run server
	r.Run(":8080")
}

// go mod init mini-compiler
// go get github.com/gin-gonic/gin
// go run main.go
