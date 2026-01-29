package logger

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	Red     = color.New(color.FgRed).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
)

func Info(format string, args ...interface{}) {
	prefix := Blue("[INFO]")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", prefix, msg)
}

func Success(format string, args ...interface{}) {
	prefix := Green("[SUCCESS]")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", prefix, msg)
}

func Warn(format string, args ...interface{}) {
	prefix := Yellow("[WARN]")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", prefix, msg)
}

func Error(format string, args ...interface{}) {
	prefix := Red("[ERROR]")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", prefix, msg)
}

func Debug(format string, args ...interface{}) {
	prefix := Magenta("[DEBUG]")
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("%s %s\n", prefix, msg)
}

func PrintBanner() {
	banner := `
   Chainmap - Modular Nmap Workflow
	`
	color.Cyan(banner)
}
