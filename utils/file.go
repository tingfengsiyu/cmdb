package utils

import (
	"github.com/elfgzp/ssh"
	"os"
	"strings"
)

// FilePath replace ~ -> $HOME
func FilePath(path string) string {
	path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	return path
}

// FileExited check file exited
func FileExited(path string) bool {
	info, err := os.Stat(FilePath(path))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDirector IsDir
func IsDirector(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// If If
func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// SessIO SessIO
func SessIO(sess *ssh.Session) ssh.Session {
	var stdio ssh.Session
	if sess != nil {
		stdio = *sess
	} else {
		stdio = nil
	}
	return stdio
}
