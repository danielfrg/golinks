package links

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

var shortPattern = regexp.MustCompile(`^\w[\w\-\.]*$`)

func Normalize(short string) string { return strings.ReplaceAll(strings.ToLower(short), "-", "") }

func ValidateShort(short string) error {
	if Normalize(short) == "_" {
		return errors.New("The name _ is reserved for internal routes")
	}
	if !shortPattern.MatchString(short) {
		return errors.New("Use letters, numbers, underscores, hyphens or periods; start with a letter, number or underscore")
	}
	return nil
}

func ValidateDestination(destination string) error {
	if strings.TrimSpace(destination) == "" {
		return errors.New("Enter a destination")
	}
	if hasControls(destination) {
		return errors.New("Destination cannot contain control characters")
	}
	if _, err := parse(destination); err != nil {
		return errors.New("Invalid Go template: " + err.Error())
	}
	if strings.Contains(destination, "{{") {
		return nil
	}
	_, err := destinationURL(destination)
	return err
}

func hasControls(s string) bool { return strings.ContainsFunc(s, unicode.IsControl) }

func destinationURL(s string) (*url.URL, error) {
	if hasControls(s) {
		return nil, errors.New("Destination cannot contain control characters")
	}
	u, err := url.Parse(s)
	if err != nil {
		return nil, errors.New("Destination is not a valid URL")
	}
	if u.Scheme != "" && u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("Destination must use HTTP, HTTPS or a relative URL")
	}
	if (u.Scheme == "http" || u.Scheme == "https") && u.Host == "" {
		return nil, errors.New("Destination must include a hostname")
	}
	if s == "" {
		return nil, errors.New("Destination expanded to an empty URL")
	}
	return u, nil
}
