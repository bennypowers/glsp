package glsp_test

// These tests verify that protocol_3_17 types are assignment-compatible
// with protocol_3_16 types wherever the handler expects 3.16 types.
// This is critical because protocol_3_17.Handler embeds protocol_3_16.Handler,
// so handler func signatures use protocol_3_16 param types.

import (
	"testing"

	"github.com/bennypowers/glsp"
	protocol316 "github.com/bennypowers/glsp/protocol_3_16"
	protocol317 "github.com/bennypowers/glsp/protocol_3_17"
)

// TestHandlerFuncAssignment verifies that consumers can set handler functions
// using protocol_3_17 types directly, without importing protocol_3_16.
func TestHandlerFuncAssignment(t *testing.T) {
	handler := &protocol317.Handler{}

	// TextDocumentDidOpen: param type must be assignable
	handler.TextDocumentDidOpen = func(ctx *glsp.Context, params *protocol317.DidOpenTextDocumentParams) error {
		_ = params.TextDocument.URI // must be accessible
		return nil
	}

	// TextDocumentCompletion: param type must be assignable
	handler.TextDocumentCompletion = func(ctx *glsp.Context, params *protocol317.CompletionParams) (any, error) {
		_ = params.TextDocument.URI
		_ = params.Position.Line
		return nil, nil
	}

	// TextDocumentHover: param type must be assignable
	handler.TextDocumentHover = func(ctx *glsp.Context, params *protocol317.HoverParams) (*protocol317.Hover, error) {
		return &protocol317.Hover{
			Contents: protocol317.MarkupContent{
				Kind:  protocol317.MarkupKindMarkdown,
				Value: "hello",
			},
		}, nil
	}

	// TextDocumentDefinition
	handler.TextDocumentDefinition = func(ctx *glsp.Context, params *protocol317.DefinitionParams) (any, error) {
		return protocol317.Location{
			URI:   "file:///a.go",
			Range: protocol317.Range{Start: protocol317.Position{Line: 0, Character: 0}, End: protocol317.Position{Line: 0, Character: 5}},
		}, nil
	}

	// TextDocumentCodeAction
	handler.TextDocumentCodeAction = func(ctx *glsp.Context, params *protocol317.CodeActionParams) (any, error) {
		return []protocol317.CodeAction{
			{Title: "fix", Kind: strPtr(string(protocol317.CodeActionKindQuickFix))},
		}, nil
	}

	// TextDocumentDidChange
	handler.TextDocumentDidChange = func(ctx *glsp.Context, params *protocol317.DidChangeTextDocumentParams) error {
		return nil
	}

	// TextDocumentDidClose
	handler.TextDocumentDidClose = func(ctx *glsp.Context, params *protocol317.DidCloseTextDocumentParams) error {
		return nil
	}

	// TextDocumentReferences
	handler.TextDocumentReferences = func(ctx *glsp.Context, params *protocol317.ReferenceParams) ([]protocol317.Location, error) {
		return nil, nil
	}

	// SetTrace
	handler.SetTrace = func(ctx *glsp.Context, params *protocol317.SetTraceParams) error {
		_ = params.Value
		return nil
	}

	// Workspace
	handler.WorkspaceDidChangeConfiguration = func(ctx *glsp.Context, params *protocol317.DidChangeConfigurationParams) error {
		return nil
	}

	handler.WorkspaceDidChangeWatchedFiles = func(ctx *glsp.Context, params *protocol317.DidChangeWatchedFilesParams) error {
		return nil
	}

	_ = handler
}

// TestDiagnosticTypeCompatibility verifies that Diagnostic and its field types
// (Range, Position, DiagnosticSeverity) are the same type across versions.
func TestDiagnosticTypeCompatibility(t *testing.T) {
	// Build a diagnostic using 317 types
	diag317 := protocol317.Diagnostic{
		Range: protocol317.Range{
			Start: protocol317.Position{Line: 1, Character: 0},
			End:   protocol317.Position{Line: 1, Character: 5},
		},
		Message:  "error",
		Severity: diagSevPtr(protocol317.DiagnosticSeverityError),
	}

	// Must be assignable to 316 Diagnostic (same type via alias)
	var diag316 protocol316.Diagnostic = diag317
	_ = diag316

	// And back
	var roundtrip protocol317.Diagnostic = diag316
	if roundtrip.Message != "error" {
		t.Errorf("roundtrip failed: %s", roundtrip.Message)
	}

	// Slices must be assignable
	var slice316 []protocol316.Diagnostic = []protocol317.Diagnostic{diag317}
	_ = slice316
}

// TestStructuralTypeCompatibility verifies key structural types
// are identical across versions.
func TestStructuralTypeCompatibility(t *testing.T) {
	// TextDocumentIdentifier
	var tdid protocol316.TextDocumentIdentifier = protocol317.TextDocumentIdentifier{URI: "file:///a.go"}
	_ = tdid

	// CompletionItem
	var ci protocol316.CompletionItem = protocol317.CompletionItem{Label: "foo"}
	_ = ci

	// TextEdit
	var te protocol316.TextEdit = protocol317.TextEdit{
		Range:   protocol317.Range{Start: protocol317.Position{Line: 0, Character: 0}, End: protocol317.Position{Line: 0, Character: 5}},
		NewText: "bar",
	}
	_ = te

	// WorkspaceEdit
	var we protocol316.WorkspaceEdit = protocol317.WorkspaceEdit{}
	_ = we

	// MarkupKind constants must be the same type
	var mk protocol316.MarkupKind = protocol317.MarkupKindMarkdown
	if mk != "markdown" {
		t.Errorf("MarkupKind mismatch: %s", mk)
	}

	// TraceValue constants
	var tv protocol316.TraceValue = protocol317.TraceValueOff
	if tv != "off" {
		t.Errorf("TraceValue mismatch: %s", tv)
	}

	// CodeActionKind
	var cak protocol316.CodeActionKind = protocol317.CodeActionKindQuickFix
	_ = cak

	// MessageType
	var mt protocol316.MessageType = protocol317.MessageTypeError
	_ = mt
}

// TestMethodConstantCompat verifies method constants are usable from both packages.
func TestMethodConstantCompat(t *testing.T) {
	// These should be the same string value regardless of origin
	if protocol317.MethodTextDocumentCompletion != protocol316.MethodTextDocumentCompletion {
		t.Error("MethodTextDocumentCompletion mismatch")
	}
	if protocol317.MethodTextDocumentHover != protocol316.MethodTextDocumentHover {
		t.Error("MethodTextDocumentHover mismatch")
	}
	if protocol317.MethodTextDocumentDidOpen != protocol316.MethodTextDocumentDidOpen {
		t.Error("MethodTextDocumentDidOpen mismatch")
	}
}

func strPtr(s string) *string { return &s }

func diagSevPtr(s protocol317.DiagnosticSeverity) *protocol317.DiagnosticSeverity { return &s }
