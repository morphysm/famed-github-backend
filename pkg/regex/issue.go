package regex

import (
	"fmt"
	"regexp"
)

// FindRightOfKey returns the string that in line right of the given key.
// If the key is not found or no string is in line right of the given key an empty string is returned.
func FindRightOfKey(text string, key string) (string, error) {
	r, err := regexp.Compile(fmt.Sprintf("\\**%s\\**[ \t]*([^\n\r]*)", key))
	if err != nil {
		return "", err
	}

	matches := r.FindStringSubmatch(text)
	if len(matches) == 0 {
		return "", fmt.Errorf("no matches found for %s", key)
	}
	if len(matches) == 1 {
		return "", fmt.Errorf("no in line token group right of %s found", key)
	}

	return matches[1], nil
}
