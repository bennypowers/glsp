package protocol

import (
	"encoding/json"
	"testing"

	protocol316 "github.com/bennypowers/glsp/protocol_3_16"
)

func TestPositionEncodingJSONSerialization(t *testing.T) {
	// Test position encoding kinds
	tests := []struct {
		name     string
		encoding PositionEncodingKind
		expected string
	}{
		{"UTF-8", PositionEncodingKindUTF8, `"utf-8"`},
		{"UTF-16", PositionEncodingKindUTF16, `"utf-16"`},
		{"UTF-32", PositionEncodingKindUTF32, `"utf-32"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.encoding)
			if err != nil {
				t.Fatalf("Failed to marshal %s: %v", tt.name, err)
			}
			if string(data) != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, string(data))
			}

			// Test unmarshaling
			var decoded PositionEncodingKind
			err = json.Unmarshal(data, &decoded)
			if err != nil {
				t.Fatalf("Failed to unmarshal %s: %v", tt.name, err)
			}
			if decoded != tt.encoding {
				t.Errorf("Expected %s, got %s", tt.encoding, decoded)
			}
		})
	}
}

func TestTypeHierarchyItemJSONSerialization(t *testing.T) {
	original := TypeHierarchyItem{
		Name:           "TestClass",
		Kind:           1,
		URI:            "file:///test.go",
		Range:          Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 10, Character: 0}},
		SelectionRange: Range{Start: Position{Line: 0, Character: 0}, End: Position{Line: 0, Character: 9}},
		Detail:         stringPtr("A test class"),
		Data:           map[string]any{"test": "data"},
	}

	// Test round-trip serialization
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal TypeHierarchyItem: %v", err)
	}

	var decoded TypeHierarchyItem
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal TypeHierarchyItem: %v", err)
	}

	// Verify fields
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: expected %s, got %s", original.Name, decoded.Name)
	}
	if decoded.Kind != original.Kind {
		t.Errorf("Kind mismatch: expected %d, got %d", original.Kind, decoded.Kind)
	}
	if decoded.URI != original.URI {
		t.Errorf("URI mismatch: expected %s, got %s", original.URI, decoded.URI)
	}
	if decoded.Detail == nil {
		t.Fatal("Detail should not be nil after unmarshaling")
	}
	if *decoded.Detail != *original.Detail {
		t.Errorf("Detail mismatch: expected %s, got %s", *original.Detail, *decoded.Detail)
	}
}

func TestInlayHintJSONSerialization(t *testing.T) {
	// Test with string label
	hintWithStringLabel := InlayHint{
		Position: Position{Line: 5, Character: 10},
		Label:    "i32",
		Kind:     &[]InlayHintKind{InlayHintKindType}[0],
		Tooltip:  "Type annotation",
	}

	data, err := json.Marshal(hintWithStringLabel)
	if err != nil {
		t.Fatalf("Failed to marshal InlayHint with string label: %v", err)
	}

	var decoded InlayHint
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal InlayHint with string label: %v", err)
	}

	if decoded.Position.Line != hintWithStringLabel.Position.Line {
		t.Errorf("Position.Line mismatch: expected %d, got %d", hintWithStringLabel.Position.Line, decoded.Position.Line)
	}

	// Test with label parts
	hintWithLabelParts := InlayHint{
		Position: Position{Line: 3, Character: 15},
		Label: []InlayHintLabelPart{
			{Value: "name", Tooltip: "Parameter name"},
			{Value: ": ", Tooltip: nil},
			{Value: "string", Tooltip: "Parameter type"},
		},
		Kind: &[]InlayHintKind{InlayHintKindParameter}[0],
	}

	data, err = json.Marshal(hintWithLabelParts)
	if err != nil {
		t.Fatalf("Failed to marshal InlayHint with label parts: %v", err)
	}

	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal InlayHint with label parts: %v", err)
	}
}

func TestInlineValueJSONSerialization(t *testing.T) {
	// Test InlineValueText
	inlineText := InlineValueText{
		Range: Range{Start: Position{Line: 5, Character: 0}, End: Position{Line: 5, Character: 10}},
		Text:  "42",
	}

	data, err := json.Marshal(inlineText)
	if err != nil {
		t.Fatalf("Failed to marshal InlineValueText: %v", err)
	}

	var decoded InlineValueText
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal InlineValueText: %v", err)
	}

	if decoded.Text != inlineText.Text {
		t.Errorf("Text mismatch: expected %s, got %s", inlineText.Text, decoded.Text)
	}

	// Test InlineValueVariableLookup
	variableLookup := InlineValueVariableLookup{
		Range:               Range{Start: Position{Line: 10, Character: 5}, End: Position{Line: 10, Character: 15}},
		VariableName:        stringPtr("myVar"),
		CaseSensitiveLookup: true,
	}

	data, err = json.Marshal(variableLookup)
	if err != nil {
		t.Fatalf("Failed to marshal InlineValueVariableLookup: %v", err)
	}

	var decodedLookup InlineValueVariableLookup
	err = json.Unmarshal(data, &decodedLookup)
	if err != nil {
		t.Fatalf("Failed to unmarshal InlineValueVariableLookup: %v", err)
	}

	if decodedLookup.VariableName == nil {
		t.Fatal("VariableName should not be nil after unmarshaling")
	}
	if *decodedLookup.VariableName != *variableLookup.VariableName {
		t.Errorf("VariableName mismatch: expected %s, got %s", *variableLookup.VariableName, *decodedLookup.VariableName)
	}
}

func TestNotebookDocumentJSONSerialization(t *testing.T) {
	notebook := NotebookDocument{
		URI:          "file:///notebook.ipynb",
		NotebookType: "jupyter-notebook",
		Version:      1,
		Cells: []NotebookCell{
			{
				Kind:     1, // NotebookCellKindMarkdown
				Document: "file:///cell1.md",
				Metadata: func() *any { m := any(map[string]any{"tags": []string{"markdown"}}); return &m }(),
			},
			{
				Kind:     2, // NotebookCellKindCode
				Document: "file:///cell2.py",
				ExecutionSummary: &ExecutionSummary{
					ExecutionOrder: 1,
					Success:        boolPtr(true),
				},
			},
		},
	}

	data, err := json.Marshal(notebook)
	if err != nil {
		t.Fatalf("Failed to marshal NotebookDocument: %v", err)
	}

	var decoded NotebookDocument
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal NotebookDocument: %v", err)
	}

	if decoded.NotebookType != notebook.NotebookType {
		t.Errorf("NotebookType mismatch: expected %s, got %s", notebook.NotebookType, decoded.NotebookType)
	}

	if len(decoded.Cells) != len(notebook.Cells) {
		t.Errorf("Cells length mismatch: expected %d, got %d", len(notebook.Cells), len(decoded.Cells))
	}

	if decoded.Cells[0].Kind != 1 { // NotebookCellKindMarkdown
		t.Errorf("First cell kind mismatch: expected %d, got %d", 1, decoded.Cells[0].Kind)
	}

	if decoded.Cells[1].ExecutionSummary == nil {
		t.Fatal("ExecutionSummary should not be nil after unmarshaling")
	}
	if decoded.Cells[1].ExecutionSummary.ExecutionOrder != 1 {
		t.Errorf("ExecutionOrder mismatch: expected 1, got %d", decoded.Cells[1].ExecutionSummary.ExecutionOrder)
	}
}

func TestDiagnosticPullModelJSONSerialization(t *testing.T) {
	// Test DocumentDiagnosticParams
	params := DocumentDiagnosticParams{
		TextDocument:     protocol316.TextDocumentIdentifier{URI: "file:///test.go"},
		Identifier:       stringPtr("go-lsp"),
		PreviousResultId: stringPtr("result-123"),
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal DocumentDiagnosticParams: %v", err)
	}

	var decoded DocumentDiagnosticParams
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal DocumentDiagnosticParams: %v", err)
	}

	if *decoded.Identifier != *params.Identifier {
		t.Errorf("Identifier mismatch: expected %s, got %s", *params.Identifier, *decoded.Identifier)
	}

	// Test FullDocumentDiagnosticReport
	report := FullDocumentDiagnosticReport{
		Kind:     string(DocumentDiagnosticReportKindFull),
		ResultID: stringPtr("result-456"),
		Items: []protocol316.Diagnostic{
			{
				Range:    protocol316.Range{Start: protocol316.Position{Line: 1, Character: 0}, End: protocol316.Position{Line: 1, Character: 5}},
				Message:  "Test diagnostic",
				Severity: func() *protocol316.DiagnosticSeverity { s := protocol316.DiagnosticSeverityError; return &s }(),
			},
		},
	}

	data, err = json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal FullDocumentDiagnosticReport: %v", err)
	}

	var decodedReport FullDocumentDiagnosticReport
	err = json.Unmarshal(data, &decodedReport)
	if err != nil {
		t.Fatalf("Failed to unmarshal FullDocumentDiagnosticReport: %v", err)
	}

	if decodedReport.Kind != report.Kind {
		t.Errorf("Kind mismatch: expected %s, got %s", report.Kind, decodedReport.Kind)
	}

	if len(decodedReport.Items) != 1 {
		t.Errorf("Items length mismatch: expected 1, got %d", len(decodedReport.Items))
	}
}

func TestClientCapabilitiesJSONSerialization(t *testing.T) {
	capabilities := ClientCapabilities{
		TextDocument: &TextDocumentClientCapabilities{
			Diagnostic: &DiagnosticClientCapabilities{
				DynamicRegistration:    true,
				RelatedDocumentSupport: false,
			},
			TypeHierarchy: &TypeHierarchyClientCapabilities{
				DynamicRegistration: boolPtr(true),
			},
			InlineValue: &InlineValueClientCapabilities{
				DynamicRegistration: boolPtr(false),
			},
			InlayHint: &InlayHintClientCapabilities{
				DynamicRegistration: boolPtr(true),
				ResolveSupport: &struct {
					Properties []string `json:"properties"`
				}{
					Properties: []string{"tooltip", "textEdits"},
				},
			},
		},
		NotebookDocument: &NotebookDocumentClientCapabilities{
			Synchronization: NotebookDocumentSyncClientCapabilities{
				DynamicRegistration:     boolPtr(true),
				ExecutionSummarySupport: boolPtr(true),
			},
		},
		General: &GeneralClientCapabilities{
			PositionEncodings: []PositionEncodingKind{
				PositionEncodingKindUTF8,
				PositionEncodingKindUTF16,
			},
		},
	}

	data, err := json.Marshal(capabilities)
	if err != nil {
		t.Fatalf("Failed to marshal ClientCapabilities: %v", err)
	}

	var decoded ClientCapabilities
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ClientCapabilities: %v", err)
	}

	// Verify key fields
	if decoded.TextDocument.Diagnostic.DynamicRegistration != true {
		t.Error("Diagnostic.DynamicRegistration should be true")
	}

	if len(decoded.General.PositionEncodings) != 2 {
		t.Errorf("Expected 2 position encodings, got %d", len(decoded.General.PositionEncodings))
	}

	if decoded.General.PositionEncodings[0] != PositionEncodingKindUTF8 {
		t.Errorf("Expected first encoding to be UTF-8, got %s", decoded.General.PositionEncodings[0])
	}
}

func TestWorkspaceDiagnosticJSONSerialization(t *testing.T) {
	// Test WorkspaceDiagnosticParams roundtrip
	params := WorkspaceDiagnosticParams{
		Identifier: stringPtr("go-lsp"),
		PreviousResultIds: []PreviousResultId{
			{URI: "file:///a.go", Value: "result-1"},
			{URI: "file:///b.go", Value: "result-2"},
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("Failed to marshal WorkspaceDiagnosticParams: %v", err)
	}

	var decodedParams WorkspaceDiagnosticParams
	err = json.Unmarshal(data, &decodedParams)
	if err != nil {
		t.Fatalf("Failed to unmarshal WorkspaceDiagnosticParams: %v", err)
	}

	if *decodedParams.Identifier != *params.Identifier {
		t.Errorf("Identifier mismatch: expected %s, got %s", *params.Identifier, *decodedParams.Identifier)
	}

	if len(decodedParams.PreviousResultIds) != 2 {
		t.Fatalf("PreviousResultIds length mismatch: expected 2, got %d", len(decodedParams.PreviousResultIds))
	}

	if decodedParams.PreviousResultIds[0].URI != "file:///a.go" {
		t.Errorf("PreviousResultIds[0].URI mismatch: expected file:///a.go, got %s", decodedParams.PreviousResultIds[0].URI)
	}

	if decodedParams.PreviousResultIds[1].Value != "result-2" {
		t.Errorf("PreviousResultIds[1].Value mismatch: expected result-2, got %s", decodedParams.PreviousResultIds[1].Value)
	}

	// Test WorkspaceDiagnosticReport roundtrip
	report := WorkspaceDiagnosticReport{
		Items: []WorkspaceDocumentDiagnosticReport{
			WorkspaceFullDocumentDiagnosticReport{
				FullDocumentDiagnosticReport: FullDocumentDiagnosticReport{
					Kind:     string(DocumentDiagnosticReportKindFull),
					ResultID: stringPtr("result-3"),
					Items: []protocol316.Diagnostic{
						{
							Range:   protocol316.Range{Start: protocol316.Position{Line: 1, Character: 0}, End: protocol316.Position{Line: 1, Character: 5}},
							Message: "unused variable",
						},
					},
				},
				URI:     "file:///a.go",
				Version: func() *protocol316.Integer { v := protocol316.Integer(3); return &v }(),
			},
		},
	}

	data, err = json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal WorkspaceDiagnosticReport: %v", err)
	}

	var decodedReport WorkspaceDiagnosticReport
	err = json.Unmarshal(data, &decodedReport)
	if err != nil {
		t.Fatalf("Failed to unmarshal WorkspaceDiagnosticReport: %v", err)
	}

	if len(decodedReport.Items) != 1 {
		t.Fatalf("Items length mismatch: expected 1, got %d", len(decodedReport.Items))
	}
}

func TestServerCapabilitiesJSONSerialization(t *testing.T) {
	capabilities := ServerCapabilities{
		DiagnosticProvider: DiagnosticOptions{
			InterFileDependencies: true,
			WorkspaceDiagnostics:  true,
			Identifier:            stringPtr("test-server"),
		},
		TypeHierarchyProvider: true,
		InlineValueProvider:   &InlineValueOptions{},
		InlayHintProvider: &InlayHintOptions{
			ResolveProvider: boolPtr(true),
		},
		NotebookDocumentSync: &NotebookDocumentSyncOptions{
			NotebookSelector: []NotebookSelector{
				{
					Notebook: "jupyter-notebook",
					Cells: []NotebookCellSelector{
						{Language: "python"},
						{Language: "markdown"},
					},
				},
			},
			Save: boolPtr(true),
		},
	}

	data, err := json.Marshal(capabilities)
	if err != nil {
		t.Fatalf("Failed to marshal ServerCapabilities: %v", err)
	}

	var decoded ServerCapabilities
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal ServerCapabilities: %v", err)
	}

	// Verify some key fields (the complex unmarshaling logic would need more comprehensive testing)
	if decoded.TypeHierarchyProvider != true {
		t.Error("TypeHierarchyProvider should be true")
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
