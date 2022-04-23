package es

import (
	es7 "github.com/olivere/elastic/v7"
)

type Example struct {
	Client *es7.Client
}

// Index see Example
func (Example) Index() {}

// Search see Example
func (Example) Search() {}

// SearchFromSize see Example
func (Example) SearchFromSize() {}

// SearchMatch see example
func (Example) SearchMatch() {}

// SearchMatchPhrase see example
func (Example) SearchMatchPhrase() {}

// SearchBool see example
func (Example) SearchBool() {}

// SearchBoolFilter see example
func (Example) SearchBoolFilter() {}

// SearchGroupBy see example
func (Example) SearchGroupBy() {}
