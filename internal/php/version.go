package php

import "fmt"

// Version selects the semantics of the supported PHP subset, not a separate
// parser or evaluator. Official comparison runs pin the exact patch version.
type Version string

const (
	PHP83 Version = "8.3"
	PHP84 Version = "8.4"
	PHP85 Version = "8.5"
)

func ParseVersion(text string) (Version, error) {
	v := Version(text)
	switch v {
	case PHP83, PHP84, PHP85:
		return v, nil
	default:
		return "", fmt.Errorf("unsupported PHP version %q (expected 8.3, 8.4 or 8.5)", text)
	}
}
