// Copyright 2026 The MathWorks, Inc.

package testbinaries

import (
	"bytes"
	"embed"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/require"
)

//go:embed source
var source embed.FS

const goModName = "mathworks.com/sdk-module"
const goModTemplate = `module ` + goModName + `

replace github.com/matlab/matlab-mcp-core-server => {{ .SDKPath }}
`

func BuildEmptyServer(t *testing.T) ServerDetails {
	// Those string literals match the one in the source code
	return ServerDetails{
		binaryLocation: buildSDKServer(t, "empty_server"),

		moduleName: goModName,

		name:         "empty-server",
		title:        "Empty Server",
		instructions: "This is the Empty Server test binary",
	}
}

func BuildServerWithCustomTools(t *testing.T) ServerDetails {
	// Those string literals match the one in the source code
	return ServerDetails{
		binaryLocation: buildSDKServer(t, "server_with_custom_tools"),

		moduleName: goModName,

		name:         "server-with-custom-tools",
		title:        "Server With Custom Tools",
		instructions: "This is a test server with custom tools",
	}
}

func buildSDKServer(t *testing.T, serverFolder string) string {
	t.Helper()

	tempDir := t.TempDir()

	sourcePath := copySDKServerSourceFiles(t, serverFolder)

	cmd := exec.Command("go",
		"mod",
		"tidy",
	)
	cmd.Dir = sourcePath

	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy failed: %s", string(output))

	serverBinaryPath := filepath.Join(tempDir, "test_server")

	cmd = exec.Command("go", //nolint:gosec // Trusted variable
		"build",
		"-o",
		serverBinaryPath,
		".",
	)
	cmd.Dir = sourcePath

	output, err = cmd.CombinedOutput()
	require.NoError(t, err, "go build failed: %s", string(output))

	return serverBinaryPath
}

func copySDKServerSourceFiles(t *testing.T, serverFolder string) string {
	t.Helper()

	sourceDir := filepath.Join(t.TempDir(), "source", serverFolder)
	sourceFS, err := fs.Sub(source, "source/"+serverFolder)
	require.NoError(t, err)

	err = fs.WalkDir(sourceFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		targetPath := filepath.Join(sourceDir, path)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0750)
		}
		content, err := fs.ReadFile(sourceFS, path)
		if err != nil {
			return err
		}
		return os.WriteFile(targetPath, content, 0600)
	})
	require.NoError(t, err, "failed to copy embedded source files")

	sdkPath := getSDKPath(t)
	writeGoMod(t, sourceDir, sdkPath)

	return sourceDir
}

func getSDKPath(t *testing.T) string {
	t.Helper()

	cmd := exec.Command("go", "env", "GOMOD")
	output, err := cmd.Output()
	require.NoError(t, err, "failed to run 'go env GOMOD'")

	goModPath := strings.TrimSpace(string(output))
	require.NotEmpty(t, goModPath, "go env GOMOD returned empty path")

	return filepath.Dir(goModPath)
}

type goModData struct {
	SDKPath string
}

func writeGoMod(t *testing.T, dir string, sdkPath string) {
	t.Helper()

	tmpl, err := template.New("go.mod").Parse(goModTemplate)
	require.NoError(t, err, "failed to parse go.mod template")

	data := goModData{
		SDKPath: sdkPath,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	require.NoError(t, err, "failed to execute go.mod template")

	goModPath := filepath.Join(dir, "go.mod")
	err = os.WriteFile(goModPath, buf.Bytes(), 0600)
	require.NoError(t, err, "failed to write go.mod file")
}
