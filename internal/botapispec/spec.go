// Package botapispec holds the machine-readable Bot API spec that telebot
// implements, generated from https://core.telegram.org/bots/api by ./gen.
package botapispec

import (
	_ "embed"
	"encoding/json"
)

//go:embed api.json
var raw []byte

// Spec is the Bot API at one version.
type Spec struct {
	Version     string   `json:"version"`
	ReleaseDate string   `json:"release_date"`
	Methods     []Method `json:"methods"`
	Types       []Type   `json:"types"`
}

// Method is a Bot API method.
type Method struct {
	Name        string  `json:"name"`
	Section     string  `json:"section"`
	Description string  `json:"description"`
	Params      []Field `json:"params,omitempty"`
}

// Type is a Bot API object. Abstract objects list their Subtypes.
type Type struct {
	Name        string   `json:"name"`
	Section     string   `json:"section"`
	Description string   `json:"description"`
	Fields      []Field  `json:"fields,omitempty"`
	Subtypes    []string `json:"subtypes,omitempty"`
}

// Field is a field of a Type or a parameter of a Method.
type Field struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// Load decodes the vendored spec.
func Load() (*Spec, error) {
	var s Spec
	return &s, json.Unmarshal(raw, &s)
}
