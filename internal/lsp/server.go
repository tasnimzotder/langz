package lsp

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	glspServer "github.com/tliron/glsp/server"

	_ "github.com/tliron/commonlog/simple"
)

const serverName = "langz-lsp"
const serverVersion = "0.1.0"

// Server holds LSP server state.
type Server struct {
	handler   protocol.Handler
	documents map[protocol.DocumentUri]string
}

// NewServer creates a new LSP server with diagnostics and hover support.
func NewServer() *Server {
	s := &Server{
		documents: make(map[protocol.DocumentUri]string),
	}

	s.handler.Initialize = s.initialize
	s.handler.Initialized = s.initialized
	s.handler.Shutdown = s.shutdown
	s.handler.SetTrace = s.setTrace
	s.handler.TextDocumentDidOpen = s.textDocumentDidOpen
	s.handler.TextDocumentDidChange = s.textDocumentDidChange
	s.handler.TextDocumentDidClose = s.textDocumentDidClose
	s.handler.TextDocumentHover = s.textDocumentHover
	s.handler.TextDocumentSignatureHelp = s.textDocumentSignatureHelp
	s.handler.TextDocumentCompletion = s.textDocumentCompletion
	s.handler.TextDocumentDefinition = s.textDocumentDefinition
	s.handler.TextDocumentDocumentSymbol = s.textDocumentDocumentSymbol
	s.handler.TextDocumentFormatting = s.textDocumentFormatting

	return s
}

// Run starts the LSP server on stdio.
func (s *Server) Run() {
	srv := glspServer.NewServer(&s.handler, serverName, false)
	srv.RunStdio()
}

func (s *Server) initialize(ctx *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := s.handler.CreateServerCapabilities()

	// Use full document sync (client sends entire content on every change)
	syncKind := protocol.TextDocumentSyncKindFull
	capabilities.TextDocumentSync = syncKind
	capabilities.CompletionProvider = &protocol.CompletionOptions{}
	capabilities.SignatureHelpProvider = &protocol.SignatureHelpOptions{
		TriggerCharacters:   []string{"(", ","},
		RetriggerCharacters: []string{",", " "},
	}
	capabilities.DefinitionProvider = true
	capabilities.DocumentSymbolProvider = true
	capabilities.DocumentFormattingProvider = true

	version := serverVersion
	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    serverName,
			Version: &version,
		},
	}, nil
}

func (s *Server) initialized(ctx *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func (s *Server) shutdown(ctx *glsp.Context) error {
	return nil
}

func (s *Server) setTrace(ctx *glsp.Context, params *protocol.SetTraceParams) error {
	return nil
}

func (s *Server) textDocumentDidOpen(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	uri := params.TextDocument.URI
	content := params.TextDocument.Text
	s.documents[uri] = content
	s.publishDiagnostics(ctx, uri, content)
	return nil
}

func (s *Server) textDocumentDidChange(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	uri := params.TextDocument.URI

	// With TextDocumentSyncKindFull, we get the entire document content
	for _, change := range params.ContentChanges {
		switch c := change.(type) {
		case protocol.TextDocumentContentChangeEventWhole:
			s.documents[uri] = c.Text
		case protocol.TextDocumentContentChangeEvent:
			// Incremental change — shouldn't happen with full sync, but handle gracefully
			s.documents[uri] = c.Text
		}
	}

	if content, ok := s.documents[uri]; ok {
		s.publishDiagnostics(ctx, uri, content)
	}
	return nil
}

func (s *Server) textDocumentDidClose(ctx *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	uri := params.TextDocument.URI
	delete(s.documents, uri)

	// Clear diagnostics for the closed file
	ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []protocol.Diagnostic{},
	})
	return nil
}
