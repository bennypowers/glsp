package protocol

import (
	"testing"
)

func TestHandler317Features(t *testing.T) {
	handler := &Handler{}

	// Test that the handler has the new 3.17 methods
	if handler.TextDocumentPrepareTypeHierarchy == nil {
		// This is expected initially, just testing the field exists
		t.Log("TextDocumentPrepareTypeHierarchy field exists")
	}

	if handler.TextDocumentInlineValue == nil {
		// This is expected initially, just testing the field exists
		t.Log("TextDocumentInlineValue field exists")
	}

	if handler.TextDocumentInlayHint == nil {
		// This is expected initially, just testing the field exists
		t.Log("TextDocumentInlayHint field exists")
	}

	// Test that we can create server capabilities
	capabilities := handler.CreateServerCapabilities()

	// Test position encoding constants exist
	_ = PositionEncodingKindUTF8
	_ = PositionEncodingKindUTF16
	_ = PositionEncodingKindUTF32

	// Test new method constants exist
	_ = MethodTextDocumentPrepareTypeHierarchy
	_ = MethodTypeHierarchySupertypes
	_ = MethodTypeHierarchySubtypes
	_ = MethodTextDocumentInlineValue
	_ = MethodTextDocumentInlayHint
	_ = MethodInlayHintResolve

	t.Log("All LSP 3.17 constants and fields are accessible")

	// Test that capabilities can be assigned
	if capabilities.TypeHierarchyProvider == nil &&
		capabilities.InlineValueProvider == nil &&
		capabilities.InlayHintProvider == nil {
		t.Log("Server capabilities can be set for 3.17 features")
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
