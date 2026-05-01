package protocol

import (
	"encoding/json"
	"testing"
)

func TestSemanticTokensClientCapabilitiesUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		check   func(t *testing.T, c SemanticTokensClientCapabilities)
		wantErr bool
	}{
		{
			name: "range is bool",
			json: `{"requests":{"range":true},"tokenTypes":[],"tokenModifiers":[],"formats":[]}`,
			check: func(t *testing.T, c SemanticTokensClientCapabilities) {
				if c.Requests.Range != true {
					t.Errorf("Range = %v, want true", c.Requests.Range)
				}
			},
		},
		{
			name: "range is empty struct",
			json: `{"requests":{"range":{}},"tokenTypes":[],"tokenModifiers":[],"formats":[]}`,
			check: func(t *testing.T, c SemanticTokensClientCapabilities) {
				if _, ok := c.Requests.Range.(struct{}); !ok {
					t.Errorf("Range type = %T, want struct{}", c.Requests.Range)
				}
			},
		},
		{
			name: "full is bool",
			json: `{"requests":{"full":true},"tokenTypes":[],"tokenModifiers":[],"formats":[]}`,
			check: func(t *testing.T, c SemanticTokensClientCapabilities) {
				if c.Requests.Full != true {
					t.Errorf("Full = %v, want true", c.Requests.Full)
				}
			},
		},
		{
			name: "full is SemanticDelta",
			json: `{"requests":{"full":{"delta":true}},"tokenTypes":[],"tokenModifiers":[],"formats":[]}`,
			check: func(t *testing.T, c SemanticTokensClientCapabilities) {
				sd, ok := c.Requests.Full.(SemanticDelta)
				if !ok {
					t.Fatalf("Full type = %T, want SemanticDelta", c.Requests.Full)
				}
				if sd.Delta == nil || !*sd.Delta {
					t.Error("Full.Delta should be true")
				}
			},
		},
		{
			name: "both nil",
			json: `{"requests":{},"tokenTypes":["type"],"tokenModifiers":["mod"],"formats":["relative"]}`,
			check: func(t *testing.T, c SemanticTokensClientCapabilities) {
				if c.Requests.Range != nil {
					t.Errorf("Range = %v, want nil", c.Requests.Range)
				}
				if c.Requests.Full != nil {
					t.Errorf("Full = %v, want nil", c.Requests.Full)
				}
				if len(c.TokenTypes) != 1 || c.TokenTypes[0] != "type" {
					t.Errorf("TokenTypes = %v", c.TokenTypes)
				}
			},
		},
		{
			name: "with dynamic registration",
			json: `{"dynamicRegistration":true,"requests":{},"tokenTypes":[],"tokenModifiers":[],"formats":[]}`,
			check: func(t *testing.T, c SemanticTokensClientCapabilities) {
				if c.DynamicRegistration == nil || !*c.DynamicRegistration {
					t.Error("DynamicRegistration should be true")
				}
			},
		},
		{
			name:    "range error",
			json:    `{"requests":{"range":"bad"},"tokenTypes":[],"tokenModifiers":[],"formats":[]}`,
			wantErr: true,
		},
		{
			name:    "full error",
			json:    `{"requests":{"full":"bad"},"tokenTypes":[],"tokenModifiers":[],"formats":[]}`,
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			json:    `{invalid`,
			wantErr: true,
		},
		{
			name:    "wrong JSON type",
			json:    `42`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c SemanticTokensClientCapabilities
			err := json.Unmarshal([]byte(tt.json), &c)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.check != nil {
				tt.check(t, c)
			}
		})
	}
}

func TestServerCapabilitiesUnmarshalJSON(t *testing.T) {
	unmarshalField := func(t *testing.T, field, value string) ServerCapabilities {
		t.Helper()
		data := `{` + `"` + field + `":` + value + `}`
		var sc ServerCapabilities
		if err := json.Unmarshal([]byte(data), &sc); err != nil {
			t.Fatalf("unmarshal %s=%s: %v", field, value, err)
		}
		return sc
	}

	expectErr := func(t *testing.T, field, value string) {
		t.Helper()
		data := `{` + `"` + field + `":` + value + `}`
		var sc ServerCapabilities
		if err := json.Unmarshal([]byte(data), &sc); err == nil {
			t.Fatalf("expected error for %s=%s", field, value)
		}
	}

	// Empty object
	t.Run("empty", func(t *testing.T) {
		var sc ServerCapabilities
		if err := json.Unmarshal([]byte(`{}`), &sc); err != nil {
			t.Fatal(err)
		}
	})

	// Invalid JSON
	t.Run("invalid JSON", func(t *testing.T) {
		var sc ServerCapabilities
		if err := json.Unmarshal([]byte(`{invalid`), &sc); err == nil {
			t.Fatal("expected error")
		}
	})

	// Wrong JSON type (hits the outer else branch in UnmarshalJSON)
	t.Run("wrong JSON type", func(t *testing.T) {
		var sc ServerCapabilities
		if err := json.Unmarshal([]byte(`42`), &sc); err == nil {
			t.Fatal("expected error")
		}
	})

	// TextDocumentSync: TextDocumentSyncOptions | TextDocumentSyncKind
	t.Run("textDocumentSync options", func(t *testing.T) {
		sc := unmarshalField(t, "textDocumentSync", `{"openClose":true}`)
		if _, ok := sc.TextDocumentSync.(TextDocumentSyncOptions); !ok {
			t.Errorf("type = %T, want TextDocumentSyncOptions", sc.TextDocumentSync)
		}
	})
	t.Run("textDocumentSync kind", func(t *testing.T) {
		sc := unmarshalField(t, "textDocumentSync", `1`)
		if _, ok := sc.TextDocumentSync.(TextDocumentSyncKind); !ok {
			t.Errorf("type = %T, want TextDocumentSyncKind", sc.TextDocumentSync)
		}
	})
	t.Run("textDocumentSync error", func(t *testing.T) {
		expectErr(t, "textDocumentSync", `"bad"`)
	})

	// bool | Options pattern (2 variants)
	boolOrOptions := []struct {
		field       string
		optionsJSON string
	}{
		{"hoverProvider", `{}`},
		{"definitionProvider", `{}`},
		{"referencesProvider", `{}`},
		{"documentHighlightProvider", `{}`},
		{"documentSymbolProvider", `{}`},
		{"codeActionProvider", `{}`},
		{"documentFormattingProvider", `{}`},
		{"documentRangeFormattingProvider", `{}`},
		{"renameProvider", `{}`},
		{"workspaceSymbolProvider", `{}`},
	}

	for _, tc := range boolOrOptions {
		t.Run(tc.field+" bool", func(t *testing.T) {
			sc := unmarshalField(t, tc.field, `true`)
			_ = sc
		})
		t.Run(tc.field+" options", func(t *testing.T) {
			sc := unmarshalField(t, tc.field, tc.optionsJSON)
			_ = sc
		})
		t.Run(tc.field+" error", func(t *testing.T) {
			expectErr(t, tc.field, `"bad"`)
		})
	}

	// bool | Options | RegistrationOptions pattern (3 variants)
	boolOrOptionsOrReg := []struct {
		field       string
		optionsJSON string
		regJSON     string
	}{
		{"declarationProvider", `{}`, `{"id":"x"}`},
		{"typeDefinitionProvider", `{}`, `{"id":"x"}`},
		{"implementationProvider", `{}`, `{"id":"x"}`},
		{"colorProvider", `{}`, `{"id":"x"}`},
		{"foldingRangeProvider", `{}`, `{"id":"x"}`},
		{"selectionRangeProvider", `{}`, `{"id":"x"}`},
		{"linkedEditingRangeProvider", `{}`, `{"id":"x"}`},
		{"callHierarchyProvider", `{}`, `{"id":"x"}`},
		{"monikerProvider", `{}`, `{"id":"x"}`},
		{"typeHierarchyProvider", `{}`, `{"id":"x"}`},
		{"inlineValueProvider", `{}`, `{"id":"x"}`},
		{"inlayHintProvider", `{}`, `{"id":"x"}`},
	}

	for _, tc := range boolOrOptionsOrReg {
		t.Run(tc.field+" bool", func(t *testing.T) {
			sc := unmarshalField(t, tc.field, `true`)
			_ = sc
		})
		t.Run(tc.field+" options", func(t *testing.T) {
			sc := unmarshalField(t, tc.field, tc.optionsJSON)
			_ = sc
		})
		t.Run(tc.field+" registration", func(t *testing.T) {
			sc := unmarshalField(t, tc.field, tc.regJSON)
			_ = sc
		})
		t.Run(tc.field+" error", func(t *testing.T) {
			expectErr(t, tc.field, `"bad"`)
		})
	}

	// Options | RegistrationOptions pattern (2 variants, no bool)
	t.Run("semanticTokensProvider options", func(t *testing.T) {
		sc := unmarshalField(t, "semanticTokensProvider", `{"legend":{"tokenTypes":[],"tokenModifiers":[]},"full":true}`)
		_ = sc
	})
	t.Run("semanticTokensProvider registration", func(t *testing.T) {
		sc := unmarshalField(t, "semanticTokensProvider", `{"id":"x","legend":{"tokenTypes":[],"tokenModifiers":[]},"full":true}`)
		_ = sc
	})
	t.Run("semanticTokensProvider error", func(t *testing.T) {
		expectErr(t, "semanticTokensProvider", `"bad"`)
	})

	t.Run("diagnosticProvider options", func(t *testing.T) {
		sc := unmarshalField(t, "diagnosticProvider", `{"interFileDependencies":true,"workspaceDiagnostics":false}`)
		_ = sc
	})
	t.Run("diagnosticProvider registration", func(t *testing.T) {
		sc := unmarshalField(t, "diagnosticProvider", `{"id":"x","interFileDependencies":true,"workspaceDiagnostics":false}`)
		_ = sc
	})
	t.Run("diagnosticProvider error", func(t *testing.T) {
		expectErr(t, "diagnosticProvider", `"bad"`)
	})

	t.Run("notebookDocumentSync options", func(t *testing.T) {
		sc := unmarshalField(t, "notebookDocumentSync", `{"notebookSelector":[{"notebook":"*"}]}`)
		_ = sc
	})
	t.Run("notebookDocumentSync registration", func(t *testing.T) {
		sc := unmarshalField(t, "notebookDocumentSync", `{"id":"x","notebookSelector":[{"notebook":"*"}]}`)
		_ = sc
	})
	t.Run("notebookDocumentSync error", func(t *testing.T) {
		expectErr(t, "notebookDocumentSync", `"bad"`)
	})

	// Directly-typed fields (no raw JSON parsing)
	t.Run("completionProvider", func(t *testing.T) {
		sc := unmarshalField(t, "completionProvider", `{"triggerCharacters":["."]}`)
		if sc.CompletionProvider == nil {
			t.Fatal("CompletionProvider should not be nil")
		}
	})
	t.Run("signatureHelpProvider", func(t *testing.T) {
		sc := unmarshalField(t, "signatureHelpProvider", `{}`)
		if sc.SignatureHelpProvider == nil {
			t.Fatal("SignatureHelpProvider should not be nil")
		}
	})
	t.Run("codeLensProvider", func(t *testing.T) {
		sc := unmarshalField(t, "codeLensProvider", `{}`)
		if sc.CodeLensProvider == nil {
			t.Fatal("CodeLensProvider should not be nil")
		}
	})
	t.Run("documentLinkProvider", func(t *testing.T) {
		sc := unmarshalField(t, "documentLinkProvider", `{}`)
		if sc.DocumentLinkProvider == nil {
			t.Fatal("DocumentLinkProvider should not be nil")
		}
	})
	t.Run("executeCommandProvider", func(t *testing.T) {
		sc := unmarshalField(t, "executeCommandProvider", `{"commands":["cmd"]}`)
		if sc.ExecuteCommandProvider == nil {
			t.Fatal("ExecuteCommandProvider should not be nil")
		}
	})
	t.Run("experimental", func(t *testing.T) {
		sc := unmarshalField(t, "experimental", `{"custom": true}`)
		if sc.Experimental == nil {
			t.Fatal("Experimental should not be nil")
		}
	})
	t.Run("workspace", func(t *testing.T) {
		sc := unmarshalField(t, "workspace", `{"workspaceFolders":{"supported":true}}`)
		if sc.Workspace == nil {
			t.Fatal("Workspace should not be nil")
		}
	})
}
