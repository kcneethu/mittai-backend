package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// Get the absolute path of the executable
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to retrieve executable path:", err)
		return
	}

	// Get the directory name from the executable path
	projectName := filepath.Base(filepath.Dir(executablePath))

	fmt.Println("Project Name:", projectName)
}
