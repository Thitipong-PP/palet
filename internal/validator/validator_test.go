package validator

import (
	"testing"

	"github.com/Thitipong-PP/palet/internal/plugin"
)

// ── MissingArgs ───────────────────────────────────────────────────────────────

func TestMissingArgs_NoneRequired(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{
			{Name: "branch", Required: false},
		},
	}
	missing := MissingArgs(cmd, map[string]string{})
	if len(missing) != 0 {
		t.Errorf("expected 0 missing, got %v", missing)
	}
}

func TestMissingArgs_AllProvided(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{
			{Name: "message", Required: true},
			{Name: "branch", Required: true},
		},
	}
	vals := map[string]string{"message": "fix bug", "branch": "main"}
	missing := MissingArgs(cmd, vals)
	if len(missing) != 0 {
		t.Errorf("expected 0 missing, got %v", missing)
	}
}

func TestMissingArgs_SomeMissing(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{
			{Name: "message", Required: true},
			{Name: "branch", Required: true},
			{Name: "author", Required: false},
		},
	}
	vals := map[string]string{"message": "fix bug"}
	missing := MissingArgs(cmd, vals)
	if len(missing) != 1 || missing[0] != "branch" {
		t.Errorf("expected [branch], got %v", missing)
	}
}

func TestMissingArgs_WhitespaceCountsAsMissing(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{{Name: "message", Required: true}},
	}
	vals := map[string]string{"message": "   "}
	missing := MissingArgs(cmd, vals)
	if len(missing) != 1 {
		t.Errorf("whitespace-only value should count as missing, got %v", missing)
	}
}

func TestMissingArgs_EmptyStringCountsAsMissing(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{{Name: "message", Required: true}},
	}
	missing := MissingArgs(cmd, map[string]string{"message": ""})
	if len(missing) != 1 {
		t.Errorf("empty string should count as missing, got %v", missing)
	}
}

func TestMissingArgs_OptionalWithNoValue(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{{Name: "author", Required: false}},
	}
	missing := MissingArgs(cmd, map[string]string{})
	if len(missing) != 0 {
		t.Errorf("optional arg with no value should not be missing, got %v", missing)
	}
}

func TestMissingArgs_EmptyArgs(t *testing.T) {
	cmd := plugin.Command{Args: nil}
	missing := MissingArgs(cmd, map[string]string{})
	if len(missing) != 0 {
		t.Errorf("no args means nothing missing, got %v", missing)
	}
}

func TestMissingArgs_AllMissing(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{
			{Name: "a", Required: true},
			{Name: "b", Required: true},
			{Name: "c", Required: true},
		},
	}
	missing := MissingArgs(cmd, map[string]string{})
	if len(missing) != 3 {
		t.Errorf("expected 3 missing, got %v", missing)
	}
}

// ── ValidateChoices ───────────────────────────────────────────────────────────

func TestValidateChoices_NoChoicesAlwaysOK(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{{Name: "branch", Type: "string", Choices: nil}},
	}
	errs := ValidateChoices(cmd, map[string]string{"branch": "anything"})
	if len(errs) != 0 {
		t.Errorf("no choices defined → should always pass, got %v", errs)
	}
}

func TestValidateChoices_ValidChoice(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{{Name: "env", Choices: []string{"dev", "staging", "prod"}}},
	}
	errs := ValidateChoices(cmd, map[string]string{"env": "staging"})
	if len(errs) != 0 {
		t.Errorf("valid choice should pass, got %v", errs)
	}
}

func TestValidateChoices_InvalidChoice(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{{Name: "env", Choices: []string{"dev", "staging", "prod"}}},
	}
	errs := ValidateChoices(cmd, map[string]string{"env": "production"})
	if len(errs) != 1 {
		t.Errorf("expected 1 error, got %v", errs)
	}
}

func TestValidateChoices_EmptyValueSkipped(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{{Name: "env", Choices: []string{"dev", "prod"}}},
	}
	errs := ValidateChoices(cmd, map[string]string{"env": ""})
	if len(errs) != 0 {
		t.Errorf("empty value should be skipped for choice validation, got %v", errs)
	}
}

func TestValidateChoices_MultipleArgs(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{
			{Name: "env", Choices: []string{"dev", "prod"}},
			{Name: "tier", Choices: []string{"free", "paid"}},
		},
	}
	errs := ValidateChoices(cmd, map[string]string{"env": "bad", "tier": "alsoBad"})
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %v", errs)
	}
}

func TestValidateChoices_UnsetArgNotChecked(t *testing.T) {
	cmd := plugin.Command{
		Args: []plugin.Arg{{Name: "env", Choices: []string{"dev", "prod"}}},
	}
	// arg not present in vals at all
	errs := ValidateChoices(cmd, map[string]string{})
	if len(errs) != 0 {
		t.Errorf("unset arg should not trigger choice validation, got %v", errs)
	}
}