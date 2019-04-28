package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type UnitSpec struct {
	Name    string          `json:"name"`
	Deps    []string        `json:"deps"`
	Type    ActionType      `json:"type"`
	RawSpec json.RawMessage `json:"spec"`
	Spec    ActionSpec
}

type ActionSpec struct {
	EnsureDir  *EnsureDirSpec
	EnsureFile *EnsureFileSpec
}

type ActionType string

const (
	ActionTypeEnsureDir  ActionType = "EnsureDir"
	ActionTypeEnsureFile ActionType = "EnsureFile"
)

type FileMetadata struct {
	Path  string `json:"path"`
	User  string `json:"user"`
	Group string `json:"group"`
	Mode  string `json:"mode"`
}

type DataSource struct {
	Path string `json:"path"`
	Data string `json:"data"`
}

type EnsureDirSpec struct {
	FileMetadata
}

type EnsureFileSpec struct {
	FileMetadata
	Contents DataSource `json:"contents"`
}

func ReadConfig(cb []byte) ([]*UnitSpec, error) {
	cbuf := bytes.NewBuffer(cb)
	decoder := json.NewDecoder(cbuf)
	decoder.DisallowUnknownFields()

	var units []*UnitSpec
	if err := decoder.Decode(&units); err != nil {
		return nil, err
	}

	for _, unit := range units {
		cbuf := bytes.NewBuffer(unit.RawSpec)
		decoder := json.NewDecoder(cbuf)
		decoder.DisallowUnknownFields()
		switch unit.Type {
		case ActionTypeEnsureDir:
			unit.Spec.EnsureDir = &EnsureDirSpec{}
			if err := decoder.Decode(&unit.Spec.EnsureDir); err != nil {
				return nil, err
			}
		case ActionTypeEnsureFile:
			unit.Spec.EnsureFile = &EnsureFileSpec{}
			if err := decoder.Decode(&unit.Spec.EnsureFile); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unknown action type: %q", unit.Type)
		}
	}
	return units, nil
}
