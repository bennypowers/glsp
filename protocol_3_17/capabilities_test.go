package protocol

import (
	"testing"

	"github.com/bennypowers/glsp"
	protocol316 "github.com/bennypowers/glsp/protocol_3_16"
)

func TestCreateServerCapabilitiesTextDocumentSync(t *testing.T) {
	t.Run("open+close sets OpenClose", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentDidOpen = func(_ *glsp.Context, _ *DidOpenTextDocumentParams) error { return nil }
		caps := h.CreateServerCapabilities()
		opts, ok := caps.TextDocumentSync.(*protocol316.TextDocumentSyncOptions)
		if !ok {
			t.Fatalf("type = %T, want *TextDocumentSyncOptions", caps.TextDocumentSync)
		}
		if opts.OpenClose == nil || !*opts.OpenClose {
			t.Error("OpenClose should be true")
		}
	})

	t.Run("close only sets OpenClose", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentDidClose = func(_ *glsp.Context, _ *DidCloseTextDocumentParams) error { return nil }
		caps := h.CreateServerCapabilities()
		opts := caps.TextDocumentSync.(*protocol316.TextDocumentSyncOptions)
		if opts.OpenClose == nil || !*opts.OpenClose {
			t.Error("OpenClose should be true for close-only")
		}
	})

	t.Run("didChange sets Change to Incremental", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentDidChange = func(_ *glsp.Context, _ *DidChangeTextDocumentParams) error { return nil }
		caps := h.CreateServerCapabilities()
		opts := caps.TextDocumentSync.(*protocol316.TextDocumentSyncOptions)
		if opts.Change == nil || *opts.Change != TextDocumentSyncKindIncremental {
			t.Errorf("Change = %v, want Incremental", opts.Change)
		}
	})

	t.Run("willSave", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentWillSave = func(_ *glsp.Context, _ *WillSaveTextDocumentParams) error { return nil }
		caps := h.CreateServerCapabilities()
		opts := caps.TextDocumentSync.(*protocol316.TextDocumentSyncOptions)
		if opts.WillSave == nil || !*opts.WillSave {
			t.Error("WillSave should be true")
		}
	})

	t.Run("willSaveWaitUntil", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentWillSaveWaitUntil = func(_ *glsp.Context, _ *WillSaveTextDocumentParams) ([]TextEdit, error) {
			return nil, nil
		}
		caps := h.CreateServerCapabilities()
		opts := caps.TextDocumentSync.(*protocol316.TextDocumentSyncOptions)
		if opts.WillSaveWaitUntil == nil || !*opts.WillSaveWaitUntil {
			t.Error("WillSaveWaitUntil should be true")
		}
	})

	t.Run("didSave", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentDidSave = func(_ *glsp.Context, _ *DidSaveTextDocumentParams) error { return nil }
		caps := h.CreateServerCapabilities()
		opts := caps.TextDocumentSync.(*protocol316.TextDocumentSyncOptions)
		if opts.Save == nil {
			t.Error("Save should be set")
		}
	})

	t.Run("all sync combined reuses options", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentDidOpen = func(_ *glsp.Context, _ *DidOpenTextDocumentParams) error { return nil }
		h.TextDocumentDidChange = func(_ *glsp.Context, _ *DidChangeTextDocumentParams) error { return nil }
		h.TextDocumentWillSave = func(_ *glsp.Context, _ *WillSaveTextDocumentParams) error { return nil }
		h.TextDocumentDidSave = func(_ *glsp.Context, _ *DidSaveTextDocumentParams) error { return nil }
		h.TextDocumentDidClose = func(_ *glsp.Context, _ *DidCloseTextDocumentParams) error { return nil }
		caps := h.CreateServerCapabilities()
		opts := caps.TextDocumentSync.(*protocol316.TextDocumentSyncOptions)
		if opts.OpenClose == nil || !*opts.OpenClose {
			t.Error("OpenClose")
		}
		if opts.Change == nil || *opts.Change != TextDocumentSyncKindIncremental {
			t.Error("Change")
		}
		if opts.WillSave == nil || !*opts.WillSave {
			t.Error("WillSave")
		}
		if opts.Save == nil {
			t.Error("Save")
		}
	})
}

func TestCreateServerCapabilitiesSimpleProviders(t *testing.T) {
	tests := []struct {
		name  string
		setup func(h *Handler)
		check func(t *testing.T, c ServerCapabilities)
	}{
		{"Completion", func(h *Handler) {
			h.TextDocumentCompletion = func(_ *glsp.Context, _ *CompletionParams) (any, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.CompletionProvider == nil {
				t.Error("CompletionProvider")
			}
		}},
		{"Hover", func(h *Handler) {
			h.TextDocumentHover = func(_ *glsp.Context, _ *HoverParams) (*Hover, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.HoverProvider != true {
				t.Error("HoverProvider")
			}
		}},
		{"SignatureHelp", func(h *Handler) {
			h.TextDocumentSignatureHelp = func(_ *glsp.Context, _ *SignatureHelpParams) (*SignatureHelp, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.SignatureHelpProvider == nil {
				t.Error("SignatureHelpProvider")
			}
		}},
		{"Declaration", func(h *Handler) {
			h.TextDocumentDeclaration = func(_ *glsp.Context, _ *DeclarationParams) (any, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.DeclarationProvider != true {
				t.Error("DeclarationProvider")
			}
		}},
		{"Definition", func(h *Handler) {
			h.TextDocumentDefinition = func(_ *glsp.Context, _ *DefinitionParams) (any, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.DefinitionProvider != true {
				t.Error("DefinitionProvider")
			}
		}},
		{"TypeDefinition", func(h *Handler) {
			h.TextDocumentTypeDefinition = func(_ *glsp.Context, _ *TypeDefinitionParams) (any, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.TypeDefinitionProvider != true {
				t.Error("TypeDefinitionProvider")
			}
		}},
		{"Implementation", func(h *Handler) {
			h.TextDocumentImplementation = func(_ *glsp.Context, _ *ImplementationParams) (any, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.ImplementationProvider != true {
				t.Error("ImplementationProvider")
			}
		}},
		{"References", func(h *Handler) {
			h.TextDocumentReferences = func(_ *glsp.Context, _ *ReferenceParams) ([]Location, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.ReferencesProvider != true {
				t.Error("ReferencesProvider")
			}
		}},
		{"DocumentHighlight", func(h *Handler) {
			h.TextDocumentDocumentHighlight = func(_ *glsp.Context, _ *DocumentHighlightParams) ([]DocumentHighlight, error) {
				return nil, nil
			}
		}, func(t *testing.T, c ServerCapabilities) {
			if c.DocumentHighlightProvider != true {
				t.Error("DocumentHighlightProvider")
			}
		}},
		{"DocumentSymbol", func(h *Handler) {
			h.TextDocumentDocumentSymbol = func(_ *glsp.Context, _ *DocumentSymbolParams) (any, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.DocumentSymbolProvider != true {
				t.Error("DocumentSymbolProvider")
			}
		}},
		{"CodeAction", func(h *Handler) {
			h.TextDocumentCodeAction = func(_ *glsp.Context, _ *CodeActionParams) (any, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.CodeActionProvider != true {
				t.Error("CodeActionProvider")
			}
		}},
		{"CodeLens", func(h *Handler) {
			h.TextDocumentCodeLens = func(_ *glsp.Context, _ *CodeLensParams) ([]CodeLens, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.CodeLensProvider == nil {
				t.Error("CodeLensProvider")
			}
		}},
		{"DocumentLink", func(h *Handler) {
			h.TextDocumentDocumentLink = func(_ *glsp.Context, _ *DocumentLinkParams) ([]DocumentLink, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.DocumentLinkProvider == nil {
				t.Error("DocumentLinkProvider")
			}
		}},
		{"Color", func(h *Handler) {
			h.TextDocumentColor = func(_ *glsp.Context, _ *DocumentColorParams) ([]ColorInformation, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.ColorProvider != true {
				t.Error("ColorProvider")
			}
		}},
		{"Formatting", func(h *Handler) {
			h.TextDocumentFormatting = func(_ *glsp.Context, _ *DocumentFormattingParams) ([]TextEdit, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.DocumentFormattingProvider != true {
				t.Error("DocumentFormattingProvider")
			}
		}},
		{"RangeFormatting", func(h *Handler) {
			h.TextDocumentRangeFormatting = func(_ *glsp.Context, _ *DocumentRangeFormattingParams) ([]TextEdit, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.DocumentRangeFormattingProvider != true {
				t.Error("DocumentRangeFormattingProvider")
			}
		}},
		{"OnTypeFormatting", func(h *Handler) {
			h.TextDocumentOnTypeFormatting = func(_ *glsp.Context, _ *DocumentOnTypeFormattingParams) ([]TextEdit, error) {
				return nil, nil
			}
		}, func(t *testing.T, c ServerCapabilities) {
			if c.DocumentOnTypeFormattingProvider == nil {
				t.Error("DocumentOnTypeFormattingProvider")
			}
		}},
		{"Rename", func(h *Handler) {
			h.TextDocumentRename = func(_ *glsp.Context, _ *RenameParams) (*WorkspaceEdit, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.RenameProvider != true {
				t.Error("RenameProvider")
			}
		}},
		{"FoldingRange", func(h *Handler) {
			h.TextDocumentFoldingRange = func(_ *glsp.Context, _ *FoldingRangeParams) ([]FoldingRange, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.FoldingRangeProvider != true {
				t.Error("FoldingRangeProvider")
			}
		}},
		{"ExecuteCommand", func(h *Handler) {
			h.WorkspaceExecuteCommand = func(_ *glsp.Context, _ *ExecuteCommandParams) (any, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.ExecuteCommandProvider == nil {
				t.Error("ExecuteCommandProvider")
			}
		}},
		{"SelectionRange", func(h *Handler) {
			h.TextDocumentSelectionRange = func(_ *glsp.Context, _ *SelectionRangeParams) ([]SelectionRange, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.SelectionRangeProvider != true {
				t.Error("SelectionRangeProvider")
			}
		}},
		{"LinkedEditingRange", func(h *Handler) {
			h.TextDocumentLinkedEditingRange = func(_ *glsp.Context, _ *LinkedEditingRangeParams) (*LinkedEditingRanges, error) {
				return nil, nil
			}
		}, func(t *testing.T, c ServerCapabilities) {
			if c.LinkedEditingRangeProvider != true {
				t.Error("LinkedEditingRangeProvider")
			}
		}},
		{"CallHierarchy", func(h *Handler) {
			h.TextDocumentPrepareCallHierarchy = func(_ *glsp.Context, _ *CallHierarchyPrepareParams) ([]CallHierarchyItem, error) {
				return nil, nil
			}
		}, func(t *testing.T, c ServerCapabilities) {
			if c.CallHierarchyProvider != true {
				t.Error("CallHierarchyProvider")
			}
		}},
		{"Moniker", func(h *Handler) {
			h.TextDocumentMoniker = func(_ *glsp.Context, _ *MonikerParams) ([]Moniker, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.MonikerProvider != true {
				t.Error("MonikerProvider")
			}
		}},
		{"WorkspaceSymbol", func(h *Handler) {
			h.WorkspaceSymbol = func(_ *glsp.Context, _ *WorkspaceSymbolParams) ([]SymbolInformation, error) { return nil, nil }
		}, func(t *testing.T, c ServerCapabilities) {
			if c.WorkspaceSymbolProvider != true {
				t.Error("WorkspaceSymbolProvider")
			}
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{}
			tt.setup(h)
			caps := h.CreateServerCapabilities()
			tt.check(t, caps)
		})
	}
}

func TestCreateServerCapabilitiesSemanticTokens(t *testing.T) {
	t.Run("full only", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentSemanticTokensFull = func(_ *glsp.Context, _ *SemanticTokensParams) (*SemanticTokens, error) {
			return nil, nil
		}
		caps := h.CreateServerCapabilities()
		opts := caps.SemanticTokensProvider.(*protocol316.SemanticTokensOptions)
		if opts.Full != true {
			t.Errorf("Full = %v, want true", opts.Full)
		}
	})

	t.Run("full with delta", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentSemanticTokensFull = func(_ *glsp.Context, _ *SemanticTokensParams) (*SemanticTokens, error) {
			return nil, nil
		}
		h.TextDocumentSemanticTokensFullDelta = func(_ *glsp.Context, _ *SemanticTokensDeltaParams) (any, error) {
			return nil, nil
		}
		caps := h.CreateServerCapabilities()
		opts := caps.SemanticTokensProvider.(*protocol316.SemanticTokensOptions)
		sd, ok := opts.Full.(*protocol316.SemanticDelta)
		if !ok {
			t.Fatalf("Full type = %T, want *SemanticDelta", opts.Full)
		}
		if sd.Delta == nil || !*sd.Delta {
			t.Error("Full.Delta should be true")
		}
	})

	t.Run("range only", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentSemanticTokensRange = func(_ *glsp.Context, _ *SemanticTokensRangeParams) (any, error) {
			return nil, nil
		}
		caps := h.CreateServerCapabilities()
		opts := caps.SemanticTokensProvider.(*protocol316.SemanticTokensOptions)
		if opts.Range != true {
			t.Errorf("Range = %v, want true", opts.Range)
		}
	})

	t.Run("full and range share options", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentSemanticTokensFull = func(_ *glsp.Context, _ *SemanticTokensParams) (*SemanticTokens, error) {
			return nil, nil
		}
		h.TextDocumentSemanticTokensRange = func(_ *glsp.Context, _ *SemanticTokensRangeParams) (any, error) {
			return nil, nil
		}
		caps := h.CreateServerCapabilities()
		opts := caps.SemanticTokensProvider.(*protocol316.SemanticTokensOptions)
		if opts.Full != true {
			t.Error("Full should be true")
		}
		if opts.Range != true {
			t.Error("Range should be true")
		}
	})
}

func TestCreateServerCapabilitiesWorkspaceFileOps(t *testing.T) {
	t.Run("didCreate", func(t *testing.T) {
		h := &Handler{}
		h.WorkspaceDidCreateFiles = func(_ *glsp.Context, _ *CreateFilesParams) error { return nil }
		caps := h.CreateServerCapabilities()
		if caps.Workspace == nil || caps.Workspace.FileOperations == nil || caps.Workspace.FileOperations.DidCreate == nil {
			t.Error("DidCreate should be set")
		}
	})

	t.Run("willCreate", func(t *testing.T) {
		h := &Handler{}
		h.WorkspaceWillCreateFiles = func(_ *glsp.Context, _ *CreateFilesParams) (*WorkspaceEdit, error) { return nil, nil }
		caps := h.CreateServerCapabilities()
		if caps.Workspace == nil || caps.Workspace.FileOperations == nil || caps.Workspace.FileOperations.WillCreate == nil {
			t.Error("WillCreate should be set")
		}
	})

	t.Run("didRename", func(t *testing.T) {
		h := &Handler{}
		h.WorkspaceDidRenameFiles = func(_ *glsp.Context, _ *RenameFilesParams) error { return nil }
		caps := h.CreateServerCapabilities()
		if caps.Workspace.FileOperations.DidRename == nil {
			t.Error("DidRename should be set")
		}
	})

	t.Run("willRename", func(t *testing.T) {
		h := &Handler{}
		h.WorkspaceWillRenameFiles = func(_ *glsp.Context, _ *RenameFilesParams) (*WorkspaceEdit, error) { return nil, nil }
		caps := h.CreateServerCapabilities()
		if caps.Workspace.FileOperations.WillRename == nil {
			t.Error("WillRename should be set")
		}
	})

	t.Run("didDelete", func(t *testing.T) {
		h := &Handler{}
		h.WorkspaceDidDeleteFiles = func(_ *glsp.Context, _ *DeleteFilesParams) error { return nil }
		caps := h.CreateServerCapabilities()
		if caps.Workspace.FileOperations.DidDelete == nil {
			t.Error("DidDelete should be set")
		}
	})

	t.Run("willDelete", func(t *testing.T) {
		h := &Handler{}
		h.WorkspaceWillDeleteFiles = func(_ *glsp.Context, _ *DeleteFilesParams) (*WorkspaceEdit, error) { return nil, nil }
		caps := h.CreateServerCapabilities()
		if caps.Workspace.FileOperations.WillDelete == nil {
			t.Error("WillDelete should be set")
		}
	})

	t.Run("multiple ops share workspace", func(t *testing.T) {
		h := &Handler{}
		h.WorkspaceDidCreateFiles = func(_ *glsp.Context, _ *CreateFilesParams) error { return nil }
		h.WorkspaceDidDeleteFiles = func(_ *glsp.Context, _ *DeleteFilesParams) error { return nil }
		caps := h.CreateServerCapabilities()
		if caps.Workspace.FileOperations.DidCreate == nil {
			t.Error("DidCreate")
		}
		if caps.Workspace.FileOperations.DidDelete == nil {
			t.Error("DidDelete")
		}
	})
}

func TestCreateServerCapabilities317Features(t *testing.T) {
	t.Run("diagnostic without workspace", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentDiagnostic = func(_ *glsp.Context, _ *DocumentDiagnosticParams) (any, error) { return nil, nil }
		caps := h.CreateServerCapabilities()
		opts, ok := caps.DiagnosticProvider.(DiagnosticOptions)
		if !ok {
			t.Fatalf("DiagnosticProvider type = %T", caps.DiagnosticProvider)
		}
		if !opts.InterFileDependencies {
			t.Error("InterFileDependencies should be true")
		}
		if opts.WorkspaceDiagnostics {
			t.Error("WorkspaceDiagnostics should be false without workspace handler")
		}
	})

	t.Run("diagnostic with workspace", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentDiagnostic = func(_ *glsp.Context, _ *DocumentDiagnosticParams) (any, error) { return nil, nil }
		h.WorkspaceDiagnostic = func(_ *glsp.Context, _ *WorkspaceDiagnosticParams) (*WorkspaceDiagnosticReport, error) {
			return nil, nil
		}
		caps := h.CreateServerCapabilities()
		opts := caps.DiagnosticProvider.(DiagnosticOptions)
		if !opts.WorkspaceDiagnostics {
			t.Error("WorkspaceDiagnostics should be true")
		}
	})

	t.Run("typeHierarchy", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentPrepareTypeHierarchy = func(_ *glsp.Context, _ *TypeHierarchyPrepareParams) ([]TypeHierarchyItem, error) {
			return nil, nil
		}
		caps := h.CreateServerCapabilities()
		if caps.TypeHierarchyProvider != true {
			t.Error("TypeHierarchyProvider")
		}
	})

	t.Run("inlineValue", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentInlineValue = func(_ *glsp.Context, _ *InlineValueParams) ([]InlineValue, error) { return nil, nil }
		caps := h.CreateServerCapabilities()
		if caps.InlineValueProvider != true {
			t.Error("InlineValueProvider")
		}
	})

	t.Run("inlayHint returns InlayHintOptions", func(t *testing.T) {
		h := &Handler{}
		h.TextDocumentInlayHint = func(_ *glsp.Context, _ *InlayHintParams) ([]InlayHint, error) { return nil, nil }
		caps := h.CreateServerCapabilities()
		if _, ok := caps.InlayHintProvider.(*InlayHintOptions); !ok {
			t.Errorf("InlayHintProvider type = %T, want *InlayHintOptions", caps.InlayHintProvider)
		}
	})

	t.Run("notebook via didOpen only", func(t *testing.T) {
		h := &Handler{}
		h.NotebookDocumentDidOpen = func(_ *glsp.Context, _ *DidOpenNotebookDocumentParams) error { return nil }
		caps := h.CreateServerCapabilities()
		if caps.NotebookDocumentSync == nil {
			t.Error("NotebookDocumentSync should be set")
		}
	})

	t.Run("notebook via didClose only", func(t *testing.T) {
		h := &Handler{}
		h.NotebookDocumentDidClose = func(_ *glsp.Context, _ *DidCloseNotebookDocumentParams) error { return nil }
		caps := h.CreateServerCapabilities()
		if caps.NotebookDocumentSync == nil {
			t.Error("NotebookDocumentSync should be set")
		}
	})
}
