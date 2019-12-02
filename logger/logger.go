package logger

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/viper"
	"os"
)

// Infof prints a formatted info message to stdout in green
func Infof(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if viper.GetBool("colorize") {
		msg = color.New(color.FgGreen).Sprint(msg)
	}
	fmt.Fprintln(os.Stdout, msg)
}

// Info prints an info message to stdout in green
func Info(msg string) {
	if viper.GetBool("colorize") {
		msg = color.New(color.FgGreen).Sprint(msg)
	}
	fmt.Fprintln(os.Stdout, msg)
}

// Warnf prints a formatted warning message to stdout in yellow
func Warnf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if viper.GetBool("colorize") {
		msg = color.New(color.FgYellow).Sprint(msg)
	}
	fmt.Fprintln(os.Stdout, msg)
}

// Warn prints a warning message to stdout in yellow
func Warn(msg string) {
	if viper.GetBool("colorize") {
		msg = color.New(color.FgYellow).Sprint(msg)
	}
	fmt.Fprintln(os.Stdout, msg)
}

// Debugf prints a formatted debug message to stdout in yellow
func Debugf(format string, a ...interface{}) {
	if viper.GetBool("verbose") {
		msg := fmt.Sprintf(format, a...)
		if viper.GetBool("colorize") {
			msg = color.New(color.FgYellow).Sprint(msg)
		}
		fmt.Fprintln(os.Stdout, msg)
	}

}

// Debug prints a debug message to stdout in yellow
func Debug(msg string) {
	if viper.GetBool("verbose") {
		if viper.GetBool("colorize") {
			msg = color.New(color.FgYellow).Sprint(msg)
		}
		fmt.Fprintln(os.Stdout, msg)
	}

}

// Error prints an error message to stdout in red
func Error(err error) {
	msg := fmt.Sprintf("%v", err)

	if viper.GetBool("colorize") {
		msg = color.New(color.FgRed).Sprint(msg)
	}

	fmt.Fprintf(os.Stderr, "%+v\n", msg)
}

// Fatalf prints a formatted error message to stdout in red and exits
func Fatalf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if viper.GetBool("colorize") {
		msg = color.New(color.FgYellow).Sprint(msg)
	}
	fmt.Fprintln(os.Stdout, msg)
	os.Exit(1)
}

// Fatal prints an error message to stdout in red and exits
func Fatal(msg string) {
	if viper.GetBool("colorize") {
		msg = color.New(color.FgRed).Sprint(msg)
	}
	fmt.Fprintln(os.Stdout, msg)
	os.Exit(1)
}
