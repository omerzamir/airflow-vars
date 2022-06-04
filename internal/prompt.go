package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func YesNoPrompt(label string, def bool) bool {
	choices := "y/N"
	if def {
		choices = "Y/n"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		_, err := fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		cobra.CheckErr(err)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
