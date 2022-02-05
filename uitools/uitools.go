// Package uitools contains helper utils for CLI UI.
package uitools

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// InputString waits user input on stdin and returns it without new line symbol.
func InputString(msg string) (string, error) {
	fmt.Print(msg, " ")
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return s[:len(s)-1], nil
}

// InputStringWithDefault calls InputString and returns defaultValue if InputString has finished without error and has return empty string.
func InputStringWithDefault(msg string, defaultValue string) (string, error) {
	msg = msg + " (default=" + defaultValue + ")"
	s, err := InputString(msg)
	if err != nil {
		return "", err
	}
	if s == "" {
		s = defaultValue
	}
	return s, nil
}

// InputYesNo waits on stdin user answer to Yes/No question.
// Returns true on yes and false on no.
func InputYesNo(msg string, yesByDefault bool) (bool, error) {
	question := "y/N"
	if yesByDefault {
		question = "Y/n"
	}
	s, err := InputString(fmt.Sprintf("%s [%s]", msg, question))
	if err != nil {
		return false, err
	}

	s = strings.ToLower(s)
	switch s {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	case "":
		return yesByDefault, nil
	}

	return false, fmt.Errorf("unknown input: %s", s)
}
