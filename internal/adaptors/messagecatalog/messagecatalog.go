// Copyright 2025 The MathWorks, Inc.

package messagecatalog

import "github.com/matlab/matlab-mcp-core-server/internal/messages"

type MessageCatalog struct {
	catalog *messages.Catalog
}

func New() *MessageCatalog {
	return &MessageCatalog{catalog: messages.NewCatalog(messages.Locale_en_US)}
}

func (m *MessageCatalog) Get(key messages.MessageKey) string {
	return m.catalog.Get(key)
}

func (m *MessageCatalog) GetFromGeneralError(err error) (string, bool) {
	return messages.FromGeneralError(m.catalog, err)
}

func (m *MessageCatalog) GetFromError(err messages.Error) string {
	return messages.FromError(m.catalog, err)
}
