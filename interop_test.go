package glsp_test

// Integration tests proving protocol version type boundaries are correct:
// - protocol_3_17 types are self-contained (no protocol_3_16 import needed)
// - protocol_3_16 types remain usable independently
// - shared aliases are assignment-compatible across versions

import (
	"encoding/json"
	"testing"

	protocol316 "github.com/bennypowers/glsp/protocol_3_16"
	protocol317 "github.com/bennypowers/glsp/protocol_3_17"
)

// TestProtocol317SelfContained verifies that constructing 3.17 diagnostic
// structs requires no protocol_3_16 imports.
func TestProtocol317SelfContained(t *testing.T) {
	diag := protocol317.Diagnostic{
		Range: protocol317.Range{
			Start: protocol317.Position{Line: 0, Character: 0},
			End:   protocol317.Position{Line: 0, Character: 10},
		},
		Message: "something went wrong",
	}

	report := protocol317.FullDocumentDiagnosticReport{
		Kind:  string(protocol317.DocumentDiagnosticReportKindFull),
		Items: []protocol317.Diagnostic{diag},
	}

	params := protocol317.PublishDiagnosticsParams{
		URI:         "file:///main.go",
		Diagnostics: []protocol317.Diagnostic{diag},
	}

	// Round-trip through JSON to verify wire format is identical
	reportJSON, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("marshal FullDocumentDiagnosticReport: %v", err)
	}
	var decodedReport protocol317.FullDocumentDiagnosticReport
	if err := json.Unmarshal(reportJSON, &decodedReport); err != nil {
		t.Fatalf("unmarshal FullDocumentDiagnosticReport: %v", err)
	}
	if len(decodedReport.Items) != 1 || decodedReport.Items[0].Message != "something went wrong" {
		t.Errorf("unexpected decoded report: %+v", decodedReport)
	}

	paramsJSON, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal PublishDiagnosticsParams: %v", err)
	}
	var decodedParams protocol317.PublishDiagnosticsParams
	if err := json.Unmarshal(paramsJSON, &decodedParams); err != nil {
		t.Fatalf("unmarshal PublishDiagnosticsParams: %v", err)
	}
	if decodedParams.URI != "file:///main.go" {
		t.Errorf("URI mismatch: %s", decodedParams.URI)
	}
}

// TestProtocol316SelfContained verifies protocol_3_16 works independently.
func TestProtocol316SelfContained(t *testing.T) {
	diag := protocol316.Diagnostic{
		Range: protocol316.Range{
			Start: protocol316.Position{Line: 1, Character: 0},
			End:   protocol316.Position{Line: 1, Character: 5},
		},
		Message: "error in 3.16 code",
	}

	params := protocol316.PublishDiagnosticsParams{
		URI:         "file:///old.go",
		Diagnostics: []protocol316.Diagnostic{diag},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded protocol316.PublishDiagnosticsParams
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(decoded.Diagnostics) != 1 || decoded.Diagnostics[0].Message != "error in 3.16 code" {
		t.Errorf("unexpected decoded params: %+v", decoded)
	}
}

// TestSharedAliasesCompatible verifies that types aliased in 3.17 from 3.16
// are assignment-compatible (MarkupKind, DocumentUri, Integer, etc.).
func TestSharedAliasesCompatible(t *testing.T) {
	// MarkupKind is an alias: protocol317.MarkupKind = protocol316.MarkupKind
	var mk317 protocol317.MarkupKind = protocol316.MarkupKindMarkdown
	var mk316 protocol316.MarkupKind = protocol317.MarkupKindPlainText
	if mk317 != "markdown" {
		t.Errorf("MarkupKind alias broken: %s", mk317)
	}
	if mk316 != "plaintext" {
		t.Errorf("MarkupKind alias broken: %s", mk316)
	}

	// DocumentUri is alias = string in both
	var uri317 protocol317.DocumentUri = "file:///a.go"
	var uri316 protocol316.DocumentUri = uri317
	if uri316 != "file:///a.go" {
		t.Errorf("DocumentUri alias broken: %s", uri316)
	}

	// Integer is alias = int32 in both
	var i317 protocol317.Integer = 42
	var i316 protocol316.Integer = i317
	if i316 != 42 {
		t.Errorf("Integer alias broken: %d", i316)
	}
}

// TestWireFormatCompatibility verifies that 3.16 and 3.17 diagnostics produce
// identical JSON (they share the same wire format for common fields).
func TestWireFormatCompatibility(t *testing.T) {
	diag316 := protocol316.PublishDiagnosticsParams{
		URI: "file:///shared.go",
		Diagnostics: []protocol316.Diagnostic{
			{
				Range:   protocol316.Range{Start: protocol316.Position{Line: 0, Character: 0}, End: protocol316.Position{Line: 0, Character: 5}},
				Message: "shared error",
			},
		},
	}

	diag317 := protocol317.PublishDiagnosticsParams{
		URI: "file:///shared.go",
		Diagnostics: []protocol317.Diagnostic{
			{
				Range:   protocol317.Range{Start: protocol317.Position{Line: 0, Character: 0}, End: protocol317.Position{Line: 0, Character: 5}},
				Message: "shared error",
			},
		},
	}

	json316, err := json.Marshal(diag316)
	if err != nil {
		t.Fatalf("marshal 316: %v", err)
	}
	json317, err := json.Marshal(diag317)
	if err != nil {
		t.Fatalf("marshal 317: %v", err)
	}

	if string(json316) != string(json317) {
		t.Errorf("wire format differs:\n316: %s\n317: %s", json316, json317)
	}
}
