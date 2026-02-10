package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// builtinSignatures maps function names to their signature information.
var builtinSignatures = map[string]protocol.SignatureInformation{
	"fetch": {
		Label: "fetch(url, method:, body:, headers:, timeout:, retries:)",
		Documentation: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "HTTP request via curl. Sets `_status`, `_body`, `_headers`.",
		},
		Parameters: []protocol.ParameterInformation{
			{Label: "url", Documentation: "The URL to request"},
			{Label: "method:", Documentation: "HTTP method (GET, POST, PUT, PATCH, DELETE)"},
			{Label: "body:", Documentation: "Request body data"},
			{Label: "headers:", Documentation: "Request headers map"},
			{Label: "timeout:", Documentation: "Max seconds to wait"},
			{Label: "retries:", Documentation: "Number of retry attempts"},
		},
	},
	"json_get": {
		Label: "json_get(data, path)",
		Documentation: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "Extract a value from JSON using a jq path. Requires `jq`.",
		},
		Parameters: []protocol.ParameterInformation{
			{Label: "data", Documentation: "JSON string to query"},
			{Label: "path", Documentation: "jq path expression (e.g. \".name\")"},
		},
	},
	"print": {
		Label: "print(args...)",
		Parameters: []protocol.ParameterInformation{
			{Label: "args...", Documentation: "Values to print"},
		},
	},
	"write": {
		Label: "write(path, content)",
		Parameters: []protocol.ParameterInformation{
			{Label: "path", Documentation: "File path to write to"},
			{Label: "content", Documentation: "Content to write"},
		},
	},
	"read": {
		Label: "read(path)",
		Parameters: []protocol.ParameterInformation{
			{Label: "path", Documentation: "File path to read"},
		},
	},
	"exec": {
		Label: "exec(command)",
		Parameters: []protocol.ParameterInformation{
			{Label: "command", Documentation: "Shell command to execute"},
		},
	},
	"copy": {
		Label: "copy(src, dst)",
		Parameters: []protocol.ParameterInformation{
			{Label: "src", Documentation: "Source file path"},
			{Label: "dst", Documentation: "Destination file path"},
		},
	},
	"move": {
		Label: "move(src, dst)",
		Parameters: []protocol.ParameterInformation{
			{Label: "src", Documentation: "Source file path"},
			{Label: "dst", Documentation: "Destination file path"},
		},
	},
	"chmod": {
		Label: "chmod(path, mode)",
		Parameters: []protocol.ParameterInformation{
			{Label: "path", Documentation: "File path"},
			{Label: "mode", Documentation: "Permission mode (e.g. 755)"},
		},
	},
	"chown": {
		Label: "chown(path, owner)",
		Parameters: []protocol.ParameterInformation{
			{Label: "path", Documentation: "File path"},
			{Label: "owner", Documentation: "Owner (user or user:group)"},
		},
	},
}

func (s *Server) textDocumentSignatureHelp(ctx *glsp.Context, params *protocol.SignatureHelpParams) (result *protocol.SignatureHelp, err error) {
	defer recoverErr(&err)
	uri := params.TextDocument.URI
	content, ok := s.documents[uri]
	if !ok {
		return nil, nil
	}

	line := int(params.Position.Line) + 1
	col := int(params.Position.Character) + 1

	return getSignatureHelp(content, line, col), nil
}

// getSignatureHelp returns signature information for the function call enclosing
// the cursor, or nil if the cursor is not inside a known function call.
func getSignatureHelp(source string, line, col int) *protocol.SignatureHelp {
	funcName, paramIdx := findEnclosingFuncCall(source, line, col)
	if funcName == "" {
		return nil
	}

	sig, ok := builtinSignatures[funcName]
	if !ok {
		return nil
	}

	activeParam := protocol.UInteger(paramIdx)
	zero := protocol.UInteger(0)
	return &protocol.SignatureHelp{
		Signatures:      []protocol.SignatureInformation{sig},
		ActiveSignature: &zero,
		ActiveParameter: &activeParam,
	}
}
