package protocol

import (
	"testing"

	"github.com/tliron/glsp"
)

func TestHandler317Features(t *testing.T) {
	// Compile-time check that 3.17 constants exist
	_ = PositionEncodingKindUTF8
	_ = PositionEncodingKindUTF16
	_ = PositionEncodingKindUTF32
	_ = MethodTextDocumentPrepareTypeHierarchy
	_ = MethodTypeHierarchySupertypes
	_ = MethodTypeHierarchySubtypes
	_ = MethodTextDocumentInlineValue
	_ = MethodTextDocumentInlayHint
	_ = MethodInlayHintResolve
	_ = MethodTextDocumentDiagnostic
	_ = MethodWorkspaceDiagnostic
	_ = MethodWorkspaceDiagnosticRefresh

	// Capabilities without handlers should have nil providers
	handler := &Handler{}
	capabilities := handler.CreateServerCapabilities()
	if capabilities.TypeHierarchyProvider != nil {
		t.Error("TypeHierarchyProvider should be nil without handler")
	}
	if capabilities.InlineValueProvider != nil {
		t.Error("InlineValueProvider should be nil without handler")
	}
	if capabilities.InlayHintProvider != nil {
		t.Error("InlayHintProvider should be nil without handler")
	}
	if capabilities.DiagnosticProvider != nil {
		t.Error("DiagnosticProvider should be nil without handler")
	}

	// Capabilities with handlers should set providers
	handler.TextDocumentPrepareTypeHierarchy = func(_ *glsp.Context, _ *TypeHierarchyPrepareParams) ([]TypeHierarchyItem, error) {
		return nil, nil
	}
	handler.TextDocumentDiagnostic = func(_ *glsp.Context, _ *DocumentDiagnosticParams) (any, error) {
		return nil, nil
	}
	handler.WorkspaceDiagnostic = func(_ *glsp.Context, _ *WorkspaceDiagnosticParams) (*WorkspaceDiagnosticReport, error) {
		return nil, nil
	}
	capabilities = handler.CreateServerCapabilities()

	if capabilities.TypeHierarchyProvider == nil {
		t.Error("TypeHierarchyProvider should be set when handler is registered")
	}
	if capabilities.DiagnosticProvider == nil {
		t.Fatal("DiagnosticProvider should be set when handler is registered")
	}
	diagOpts, ok := capabilities.DiagnosticProvider.(DiagnosticOptions)
	if !ok {
		t.Fatal("DiagnosticProvider should be DiagnosticOptions")
	}
	if !diagOpts.WorkspaceDiagnostics {
		t.Error("WorkspaceDiagnostics should be true when WorkspaceDiagnostic handler is set")
	}
}

func TestPositionEncoding(t *testing.T) {
	// Test position encoding kinds
	utf8 := PositionEncodingKindUTF8
	utf16 := PositionEncodingKindUTF16
	utf32 := PositionEncodingKindUTF32

	if utf8 != "utf-8" {
		t.Errorf("Expected utf-8, got %s", utf8)
	}
	if utf16 != "utf-16" {
		t.Errorf("Expected utf-16, got %s", utf16)
	}
	if utf32 != "utf-32" {
		t.Errorf("Expected utf-32, got %s", utf32)
	}
}

func TestTypeHierarchyItem(t *testing.T) {
	// Test that we can create a TypeHierarchyItem
	item := TypeHierarchyItem{
		Name:           "TestClass",
		Kind:           1, // Assuming this corresponds to a class
		URI:            "file:///test.go",
		Range:          Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 10, Character: 0}},
		SelectionRange: Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 0, Character: 9}},
	}

	if item.Name != "TestClass" {
		t.Errorf("Expected TestClass, got %s", item.Name)
	}
}

func TestInlayHint(t *testing.T) {
	// Test that we can create an InlayHint
	hint := InlayHint{
		Position: Position{Line: 5, Character: 10},
		Label:    "i32",
		Kind:     &[]InlayHintKind{InlayHintKindType}[0],
	}

	if hint.Position.Line != 5 {
		t.Errorf("Expected line 5, got %d", hint.Position.Line)
	}
}

func TestHandlerInitialization(t *testing.T) {
	handler := &Handler{}

	// Test handler initialization state
	if handler.IsInitialized() {
		t.Error("Handler should not be initialized by default")
	}

	handler.SetInitialized(true)
	if !handler.IsInitialized() {
		t.Error("Handler should be initialized after SetInitialized(true)")
	}
}
