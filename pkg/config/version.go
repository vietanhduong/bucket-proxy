package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

var (
	gitCommit = "unknown"
	version   = "unreleased"
	buildDate = "unknown"
	arch      = fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	userAgent = fmt.Sprintf("%s/%s", filepath.Base(os.Args[0]), version)
)

func GitCommit() string {
	return gitCommit
}

func Version() string {
	return version
}

func BuildDate() string {
	return buildDate
}

func Arch() string {
	return arch
}

func PrintVersion() {
	fmt.Fprintf(os.Stdout, "Version: %s\n", Version())
	fmt.Fprintf(os.Stdout, "Git Commit: %s\n", GitCommit())
	fmt.Fprintf(os.Stdout, "Build Date: %s\n", BuildDate())
	fmt.Fprintf(os.Stdout, "Architecture: %s\n", Arch())
}

func UserAgent() string { return userAgent }
