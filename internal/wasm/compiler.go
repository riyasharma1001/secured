package wasm

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type Compiler struct {
	TinyGoPath string
	OutputPath string
}

func NewCompiler(tinyGoPath, outputPath string) *Compiler {
	return &Compiler{
		TinyGoPath: tinyGoPath,
		OutputPath: outputPath,
	}
}

// CompileToWASM compiles Go code to WebAssembly
func (c *Compiler) CompileToWASM(sourceCode []byte) ([]byte, error) {
	// Create temporary source file
	tempDir := filepath.Join(c.OutputPath, "temp")
	sourceFile := filepath.Join(tempDir, "main.go")
	wasmFile := filepath.Join(tempDir, "main.wasm")

	// Write source code to temp file
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, err
	}
	if err := os.WriteFile(sourceFile, sourceCode, 0644); err != nil {
		return nil, err
	}

	// Prepare tinygo command
	cmd := exec.Command(c.getTinyGoCmd(),
		"build",
		"-o", wasmFile,
		"-target=wasm",
		sourceFile,
	)

	// Capture command output
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run compilation
	if err := cmd.Run(); err != nil {
		return nil, errors.New(stderr.String())
	}

	// Read compiled WASM
	wasmBytes, err := os.ReadFile(wasmFile)
	if err != nil {
		return nil, err
	}

	// Cleanup
	os.RemoveAll(tempDir)

	return wasmBytes, nil
}

// getTinyGoCmd returns the correct tinygo command based on OS
func (c *Compiler) getTinyGoCmd() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(c.TinyGoPath, "tinygo.exe")
	}
	return filepath.Join(c.TinyGoPath, "tinygo")
}

// GenerateWASMLoader generates the JavaScript loader for WASM
func (c *Compiler) GenerateWASMLoader(wasmURL string) []byte {
	loader := `
        const wasmLoader = async () => {
            try {
                const wasmInstance = await WebAssembly.instantiateStreaming(
                    fetch("` + wasmURL + `"),
                    {
                        env: {
                            // Add any required environment functions
                        }
                    }
                );
                return wasmInstance.instance.exports;
            } catch (err) {
                console.error("Failed to load WASM:", err);
                return null;
            }
        };
    `
	return []byte(loader)
}

// Helper function to check if TinyGo is installed
func (c *Compiler) CheckTinyGoInstallation() error {
	cmd := exec.Command(c.getTinyGoCmd(), "version")
	if err := cmd.Run(); err != nil {
		return errors.New("TinyGo not found. Please install TinyGo first")
	}
	return nil
}
