package comment

import (
	"encoding/json"
	"fmt"
)

type Type string

const (
	RewardCommentType   Type = "reward"
	EligibleCommentType Type = "eligible"
)

type Identifier struct {
	Type    Type   `json:"type"`
	Version string `json:"version"`
}

func NewIdentifier(commentType Type, version string) Identifier {
	return Identifier{Type: commentType, Version: version}
}

func (i Identifier) String() (string, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("<!--%s-->", string(b)), err
}
