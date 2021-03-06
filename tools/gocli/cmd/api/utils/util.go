package utils

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/cocoup/go-smart/tools/gocli/cmd/api/spec"
)

func GetGitName() string {
	cmd := exec.Command("git", "config", "user.name")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func GetGitEmail() string {
	cmd := exec.Command("git", "config", "user.email")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// ComponentName returns component name for typescript
func ComponentName(api *spec.ApiSpec) string {
	name := api.Service.Name
	if strings.HasSuffix(name, "-api") {
		return name[:len(name)-4] + "Components"
	}
	return name + "Components"
}

// WriteIndent writes tab spaces
func WriteIndent(writer io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		fmt.Fprint(writer, "\t")
	}
}

// RemoveComment filters comment content
func RemoveComment(line string) string {
	commentIdx := strings.Index(line, "//")
	if commentIdx >= 0 {
		return strings.TrimSpace(line[:commentIdx])
	}
	return strings.TrimSpace(line)
}
