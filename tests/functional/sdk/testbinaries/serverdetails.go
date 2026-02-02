// Copyright 2026 The MathWorks, Inc.

package testbinaries

type ServerDetails struct {
	binaryLocation string

	moduleName string

	name         string
	title        string
	instructions string
}

func (s ServerDetails) BinaryLocation() string {
	return s.binaryLocation
}

func (s ServerDetails) ModuleName() string {
	return s.moduleName
}

func (s ServerDetails) Name() string {
	return s.name
}

func (s ServerDetails) Title() string {
	return s.title
}

func (s ServerDetails) Instructions() string {
	return s.instructions
}
