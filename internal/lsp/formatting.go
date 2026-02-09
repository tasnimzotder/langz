package lsp

import (
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (s *Server) textDocumentFormatting(ctx *glsp.Context, params *protocol.DocumentFormattingParams) ([]protocol.TextEdit, error) {
	uri := params.TextDocument.URI
	content, ok := s.documents[uri]
	if !ok {
		return nil, nil
	}

	tabSize := 4
	if v, ok := params.Options["tabSize"]; ok {
		if n, ok := v.(float64); ok {
			tabSize = int(n)
		}
	}
	insertSpaces := true
	if v, ok := params.Options["insertSpaces"]; ok {
		if b, ok := v.(bool); ok {
			insertSpaces = b
		}
	}

	formatted := formatSource(content, tabSize, insertSpaces)
	if formatted == content {
		return nil, nil
	}

	lines := strings.Split(content, "\n")
	lastLine := len(lines) - 1
	lastChar := len(lines[lastLine])

	return []protocol.TextEdit{{
		Range: protocol.Range{
			Start: protocol.Position{Line: 0, Character: 0},
			End: protocol.Position{
				Line:      protocol.UInteger(lastLine),
				Character: protocol.UInteger(lastChar),
			},
		},
		NewText: formatted,
	}}, nil
}

func formatSource(source string, tabSize int, insertSpaces bool) string {
	indent := "\t"
	if insertSpaces {
		indent = strings.Repeat(" ", tabSize)
	}

	lines := strings.Split(source, "\n")
	result := make([]string, len(lines))
	level := 0

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			result[i] = ""
			continue
		}
		if strings.HasPrefix(trimmed, "}") {
			level--
			if level < 0 {
				level = 0
			}
		}
		result[i] = strings.Repeat(indent, level) + trimmed
		if strings.HasSuffix(trimmed, "{") {
			level++
		}
	}

	return strings.Join(result, "\n")
}
