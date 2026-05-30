package validator

import (
	"fmt"
	"strings"

	"github.com/Thitipong-PP/palet/internal/plugin"
)

// MissingArgs returns the names of required args that have no value.
func MissingArgs(cmd plugin.Command, vals map[string]string) []string {
	missing := make([]string, 0, len(cmd.Args))
	for _, a := range cmd.Args {
		if a.Required && strings.TrimSpace(vals[a.Name]) == "" {
			missing = append(missing, a.Name)
		}
	}
	return missing
}

func ValidateChoices(cmd plugin.Command, vals map[string]string) []string {
    var errs []string
    for _, a := range cmd.Args {
        if len(a.Choices) == 0 {
            continue
        }
        v := vals[a.Name]
        if v == "" {
            continue
        }
        valid := false
        for _, c := range a.Choices {
            if c == v { valid = true; break }
        }
        if !valid {
            errs = append(errs, fmt.Sprintf("%s: %q is not a valid choice", a.Name, v))
        }
    }
    return errs
}