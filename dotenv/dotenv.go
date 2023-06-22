package dotenv

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func Load() {
	// Get the current file path
	_, currentFilePath, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("error getting current file path")
	}

	// Get the project root directory
	rootDir := filepath.Join(filepath.Dir(currentFilePath), "..")

	// Construct the absolute path of the .env file
	envPath := filepath.Join(rootDir, ".env")

	bytes, err := ioutil.ReadFile(envPath)
	if err != nil {
		log.Fatalf("error parsing file: %v", err.Error())
	}

	lines := strings.Split(string(bytes), "\n")

	for i, line := range lines {
		if len(strings.TrimSpace(line)) > 0 && !strings.HasPrefix(line, "#") {
			arr := strings.Split(line, "=")
			if len(arr) != 2 {
				log.Fatalf("invalid format at line %v\n%v", i+1, line)
			}

			key, value := arr[0], arr[1]
			os.Setenv(key, value)
		}
	}
}
