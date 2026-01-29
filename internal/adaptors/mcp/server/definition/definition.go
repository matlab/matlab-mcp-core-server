// Copyright 2026 The MathWorks, Inc.

package definition

type Definition struct {
	name         string
	title        string
	instructions string
}

func New(name, title, instructions string) Definition {
	return Definition{
		name:         name,
		title:        title,
		instructions: instructions,
	}
}

func (d Definition) Name() string {
	return d.name
}

func (d Definition) Title() string {
	return d.title
}

func (d Definition) Instructions() string {
	return d.instructions
}
