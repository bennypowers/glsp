package protocol

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/bennypowers/glsp"
)

func ctx(method string, params string) *glsp.Context {
	return &glsp.Context{
		Method: method,
		Params: json.RawMessage(params),
		Notify: func(string, any) {},
	}
}

func TestHandleBeforeInitialization(t *testing.T) {
	h := &Handler{}
	_, _, _, err := h.Handle(ctx(MethodTextDocumentCompletion, `{}`))
	if err == nil {
		t.Fatal("expected error for non-initialize method before initialization")
	}
}

func TestHandleInitialize(t *testing.T) {
	h := &Handler{}
	h.Initialize = func(c *glsp.Context, p *InitializeParams) (any, error) {
		return &InitializeResult{}, nil
	}
	r, vm, vp, err := h.Handle(ctx(MethodInitialize, `{"capabilities":{}}`))
	if err != nil {
		t.Fatal(err)
	}
	if !vm || !vp {
		t.Errorf("validMethod=%v validParams=%v", vm, vp)
	}
	if r == nil {
		t.Error("result should not be nil")
	}
	if !h.IsInitialized() {
		t.Error("handler should be initialized after successful Initialize")
	}
}

func TestHandleInitializeError(t *testing.T) {
	h := &Handler{}
	h.Initialize = func(c *glsp.Context, p *InitializeParams) (any, error) {
		return nil, errors.New("init failed")
	}
	_, _, _, err := h.Handle(ctx(MethodInitialize, `{"capabilities":{}}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if h.IsInitialized() {
		t.Error("should not be initialized after failed Initialize")
	}
}

func TestHandleShutdown(t *testing.T) {
	h := &Handler{}
	h.SetInitialized(true)
	h.Shutdown = func(c *glsp.Context) error { return nil }
	_, vm, vp, err := h.Handle(ctx(MethodShutdown, `null`))
	if err != nil {
		t.Fatal(err)
	}
	if !vm || !vp {
		t.Errorf("validMethod=%v validParams=%v", vm, vp)
	}
	if h.IsInitialized() {
		t.Error("should not be initialized after Shutdown")
	}
}

func TestHandleExit(t *testing.T) {
	h := &Handler{}
	h.SetInitialized(true)
	h.Exit = func(c *glsp.Context) error { return nil }
	_, vm, vp, err := h.Handle(ctx(MethodExit, `null`))
	if err != nil {
		t.Fatal(err)
	}
	if !vm || !vp {
		t.Errorf("validMethod=%v validParams=%v", vm, vp)
	}
}

func TestHandleCustomRequest(t *testing.T) {
	h := &Handler{}
	h.SetInitialized(true)
	h.CustomRequest = CustomRequestHandlers{
		"custom/test": {
			Func: func(c *glsp.Context, p json.RawMessage) (any, error) {
				return "ok", nil
			},
		},
	}
	r, vm, vp, err := h.Handle(ctx("custom/test", `{}`))
	if err != nil {
		t.Fatal(err)
	}
	if !vm || !vp {
		t.Errorf("validMethod=%v validParams=%v", vm, vp)
	}
	if r != "ok" {
		t.Errorf("result = %v, want ok", r)
	}
}

func TestHandleNoParamsMethods(t *testing.T) {
	h := &Handler{}
	h.SetInitialized(true)

	h.WorkspaceSemanticTokensRefresh = func(c *glsp.Context) error { return nil }
	_, vm, vp, err := h.Handle(ctx(MethodWorkspaceSemanticTokensRefresh, `null`))
	if err != nil {
		t.Fatal(err)
	}
	if !vm || !vp {
		t.Errorf("validMethod=%v validParams=%v", vm, vp)
	}
}

func TestHandleUnknownMethod(t *testing.T) {
	h := &Handler{}
	h.SetInitialized(true)
	_, vm, _, _ := h.Handle(ctx("unknown/method", `{}`))
	if vm {
		t.Error("validMethod should be false for unknown method")
	}
}

type handleMethodTest struct {
	name      string
	method    string
	setup     func(h *Handler)
	validJSON string
}

func Test3_17Methods(t *testing.T) {
	tests := []handleMethodTest{
		{"TextDocumentDiagnostic", MethodTextDocumentDiagnostic,
			func(h *Handler) {
				h.TextDocumentDiagnostic = func(c *glsp.Context, p *DocumentDiagnosticParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"WorkspaceDiagnostic", MethodWorkspaceDiagnostic,
			func(h *Handler) {
				h.WorkspaceDiagnostic = func(c *glsp.Context, p *WorkspaceDiagnosticParams) (*WorkspaceDiagnosticReport, error) {
					return &WorkspaceDiagnosticReport{}, nil
				}
			}, `{"previousResultIds":[]}`},
		{"PrepareTypeHierarchy", MethodTextDocumentPrepareTypeHierarchy,
			func(h *Handler) {
				h.TextDocumentPrepareTypeHierarchy = func(c *glsp.Context, p *TypeHierarchyPrepareParams) ([]TypeHierarchyItem, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TypeHierarchySupertypes", MethodTypeHierarchySupertypes,
			func(h *Handler) {
				h.TypeHierarchySupertypes = func(c *glsp.Context, p *TypeHierarchySupertypesParams) ([]TypeHierarchyItem, error) {
					return nil, nil
				}
			}, `{"item":{"name":"A","kind":5,"uri":"file:///a.go","range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"selectionRange":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}}`},
		{"TypeHierarchySubtypes", MethodTypeHierarchySubtypes,
			func(h *Handler) {
				h.TypeHierarchySubtypes = func(c *glsp.Context, p *TypeHierarchySubtypesParams) ([]TypeHierarchyItem, error) {
					return nil, nil
				}
			}, `{"item":{"name":"A","kind":5,"uri":"file:///a.go","range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"selectionRange":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}}`},
		{"TextDocumentInlineValue", MethodTextDocumentInlineValue,
			func(h *Handler) {
				h.TextDocumentInlineValue = func(c *glsp.Context, p *InlineValueParams) ([]InlineValue, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"context":{"frameId":1,"stoppedLocation":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}}`},
		{"TextDocumentInlayHint", MethodTextDocumentInlayHint,
			func(h *Handler) {
				h.TextDocumentInlayHint = func(c *glsp.Context, p *InlayHintParams) ([]InlayHint, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}`},
		{"InlayHintResolve", MethodInlayHintResolve,
			func(h *Handler) {
				h.InlayHintResolve = func(c *glsp.Context, p *InlayHint) (*InlayHint, error) { return p, nil }
			}, `{"position":{"line":0,"character":0},"label":"test"}`},
		{"NotebookDocumentDidOpen", MethodNotebookDocumentDidOpen,
			func(h *Handler) {
				h.NotebookDocumentDidOpen = func(c *glsp.Context, p *DidOpenNotebookDocumentParams) error { return nil }
			}, `{"notebookDocument":{"uri":"file:///nb.ipynb","notebookType":"jupyter","version":1,"cells":[]},"cellTextDocuments":[]}`},
		{"NotebookDocumentDidChange", MethodNotebookDocumentDidChange,
			func(h *Handler) {
				h.NotebookDocumentDidChange = func(c *glsp.Context, p *DidChangeNotebookDocumentParams) error { return nil }
			}, `{"notebookDocument":{"uri":"file:///nb.ipynb","version":2},"change":{}}`},
		{"NotebookDocumentDidSave", MethodNotebookDocumentDidSave,
			func(h *Handler) {
				h.NotebookDocumentDidSave = func(c *glsp.Context, p *DidSaveNotebookDocumentParams) error { return nil }
			}, `{"notebookDocument":{"uri":"file:///nb.ipynb"}}`},
		{"NotebookDocumentDidClose", MethodNotebookDocumentDidClose,
			func(h *Handler) {
				h.NotebookDocumentDidClose = func(c *glsp.Context, p *DidCloseNotebookDocumentParams) error { return nil }
			}, `{"notebookDocument":{"uri":"file:///nb.ipynb"},"cellTextDocuments":[]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/nil handler", func(t *testing.T) {
			h := &Handler{}
			h.SetInitialized(true)
			_, vm, _, _ := h.Handle(ctx(tt.method, tt.validJSON))
			if vm {
				t.Error("validMethod should be false when handler is nil")
			}
		})
		t.Run(tt.name+"/valid", func(t *testing.T) {
			h := &Handler{}
			h.SetInitialized(true)
			tt.setup(h)
			_, vm, vp, err := h.Handle(ctx(tt.method, tt.validJSON))
			if err != nil {
				t.Fatal(err)
			}
			if !vm || !vp {
				t.Errorf("validMethod=%v validParams=%v", vm, vp)
			}
		})
		t.Run(tt.name+"/invalid params", func(t *testing.T) {
			h := &Handler{}
			h.SetInitialized(true)
			tt.setup(h)
			_, vm, vp, _ := h.Handle(ctx(tt.method, `{invalid`))
			if !vm {
				t.Error("validMethod should be true")
			}
			if vp {
				t.Error("validParams should be false for invalid JSON")
			}
		})
	}
}

func Test3_16Methods(t *testing.T) {
	tests := []handleMethodTest{
		{"CancelRequest", MethodCancelRequest,
			func(h *Handler) { h.CancelRequest = func(c *glsp.Context, p *CancelParams) error { return nil } },
			`{"id":1}`},
		{"Progress", MethodProgress,
			func(h *Handler) { h.Progress = func(c *glsp.Context, p *ProgressParams) error { return nil } },
			`{"token":1,"value":{}}`},
		{"Initialized", MethodInitialized,
			func(h *Handler) { h.Initialized = func(c *glsp.Context, p *InitializedParams) error { return nil } },
			`{}`},
		{"LogTrace", MethodLogTrace,
			func(h *Handler) { h.LogTrace = func(c *glsp.Context, p *LogTraceParams) error { return nil } },
			`{"message":"test"}`},
		{"SetTrace", MethodSetTrace,
			func(h *Handler) { h.SetTrace = func(c *glsp.Context, p *SetTraceParams) error { return nil } },
			`{"value":"off"}`},
		{"WindowWorkDoneProgressCancel", MethodWindowWorkDoneProgressCancel,
			func(h *Handler) {
				h.WindowWorkDoneProgressCancel = func(c *glsp.Context, p *WorkDoneProgressCancelParams) error { return nil }
			}, `{"token":1}`},
		{"WorkspaceDidChangeWorkspaceFolders", MethodWorkspaceDidChangeWorkspaceFolders,
			func(h *Handler) {
				h.WorkspaceDidChangeWorkspaceFolders = func(c *glsp.Context, p *DidChangeWorkspaceFoldersParams) error { return nil }
			}, `{"event":{"added":[],"removed":[]}}`},
		{"WorkspaceDidChangeConfiguration", MethodWorkspaceDidChangeConfiguration,
			func(h *Handler) {
				h.WorkspaceDidChangeConfiguration = func(c *glsp.Context, p *DidChangeConfigurationParams) error { return nil }
			}, `{"settings":{}}`},
		{"WorkspaceDidChangeWatchedFiles", MethodWorkspaceDidChangeWatchedFiles,
			func(h *Handler) {
				h.WorkspaceDidChangeWatchedFiles = func(c *glsp.Context, p *DidChangeWatchedFilesParams) error { return nil }
			}, `{"changes":[]}`},
		{"WorkspaceSymbol", MethodWorkspaceSymbol,
			func(h *Handler) {
				h.WorkspaceSymbol = func(c *glsp.Context, p *WorkspaceSymbolParams) ([]SymbolInformation, error) { return nil, nil }
			}, `{"query":"test"}`},
		{"WorkspaceExecuteCommand", MethodWorkspaceExecuteCommand,
			func(h *Handler) {
				h.WorkspaceExecuteCommand = func(c *glsp.Context, p *ExecuteCommandParams) (any, error) { return nil, nil }
			}, `{"command":"cmd"}`},
		{"TextDocumentDidOpen", MethodTextDocumentDidOpen,
			func(h *Handler) {
				h.TextDocumentDidOpen = func(c *glsp.Context, p *DidOpenTextDocumentParams) error { return nil }
			}, `{"textDocument":{"uri":"file:///a.go","languageId":"go","version":1,"text":""}}`},
		{"TextDocumentDidChange", MethodTextDocumentDidChange,
			func(h *Handler) {
				h.TextDocumentDidChange = func(c *glsp.Context, p *DidChangeTextDocumentParams) error { return nil }
			}, `{"textDocument":{"uri":"file:///a.go","version":2},"contentChanges":[]}`},
		{"TextDocumentWillSave", MethodTextDocumentWillSave,
			func(h *Handler) {
				h.TextDocumentWillSave = func(c *glsp.Context, p *WillSaveTextDocumentParams) error { return nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"reason":1}`},
		{"TextDocumentWillSaveWaitUntil", MethodTextDocumentWillSaveWaitUntil,
			func(h *Handler) {
				h.TextDocumentWillSaveWaitUntil = func(c *glsp.Context, p *WillSaveTextDocumentParams) ([]TextEdit, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"reason":1}`},
		{"TextDocumentDidSave", MethodTextDocumentDidSave,
			func(h *Handler) {
				h.TextDocumentDidSave = func(c *glsp.Context, p *DidSaveTextDocumentParams) error { return nil }
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"TextDocumentDidClose", MethodTextDocumentDidClose,
			func(h *Handler) {
				h.TextDocumentDidClose = func(c *glsp.Context, p *DidCloseTextDocumentParams) error { return nil }
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"TextDocumentCompletion", MethodTextDocumentCompletion,
			func(h *Handler) {
				h.TextDocumentCompletion = func(c *glsp.Context, p *CompletionParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"CompletionItemResolve", MethodCompletionItemResolve,
			func(h *Handler) {
				h.CompletionItemResolve = func(c *glsp.Context, p *CompletionItem) (*CompletionItem, error) { return p, nil }
			}, `{"label":"test"}`},
		{"TextDocumentHover", MethodTextDocumentHover,
			func(h *Handler) {
				h.TextDocumentHover = func(c *glsp.Context, p *HoverParams) (*Hover, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentSignatureHelp", MethodTextDocumentSignatureHelp,
			func(h *Handler) {
				h.TextDocumentSignatureHelp = func(c *glsp.Context, p *SignatureHelpParams) (*SignatureHelp, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentDeclaration", MethodTextDocumentDeclaration,
			func(h *Handler) {
				h.TextDocumentDeclaration = func(c *glsp.Context, p *DeclarationParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentDefinition", MethodTextDocumentDefinition,
			func(h *Handler) {
				h.TextDocumentDefinition = func(c *glsp.Context, p *DefinitionParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentTypeDefinition", MethodTextDocumentTypeDefinition,
			func(h *Handler) {
				h.TextDocumentTypeDefinition = func(c *glsp.Context, p *TypeDefinitionParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentImplementation", MethodTextDocumentImplementation,
			func(h *Handler) {
				h.TextDocumentImplementation = func(c *glsp.Context, p *ImplementationParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentReferences", MethodTextDocumentReferences,
			func(h *Handler) {
				h.TextDocumentReferences = func(c *glsp.Context, p *ReferenceParams) ([]Location, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0},"context":{"includeDeclaration":false}}`},
		{"TextDocumentDocumentHighlight", MethodTextDocumentDocumentHighlight,
			func(h *Handler) {
				h.TextDocumentDocumentHighlight = func(c *glsp.Context, p *DocumentHighlightParams) ([]DocumentHighlight, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentDocumentSymbol", MethodTextDocumentDocumentSymbol,
			func(h *Handler) {
				h.TextDocumentDocumentSymbol = func(c *glsp.Context, p *DocumentSymbolParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"TextDocumentCodeAction", MethodTextDocumentCodeAction,
			func(h *Handler) {
				h.TextDocumentCodeAction = func(c *glsp.Context, p *CodeActionParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"context":{"diagnostics":[]}}`},
		{"CodeActionResolve", MethodCodeActionResolve,
			func(h *Handler) {
				h.CodeActionResolve = func(c *glsp.Context, p *CodeAction) (*CodeAction, error) { return p, nil }
			}, `{"title":"test"}`},
		{"TextDocumentCodeLens", MethodTextDocumentCodeLens,
			func(h *Handler) {
				h.TextDocumentCodeLens = func(c *glsp.Context, p *CodeLensParams) ([]CodeLens, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"CodeLensResolve", MethodCodeLensResolve,
			func(h *Handler) {
				h.CodeLensResolve = func(c *glsp.Context, p *CodeLens) (*CodeLens, error) { return p, nil }
			}, `{"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}`},
		{"TextDocumentDocumentLink", MethodTextDocumentDocumentLink,
			func(h *Handler) {
				h.TextDocumentDocumentLink = func(c *glsp.Context, p *DocumentLinkParams) ([]DocumentLink, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"DocumentLinkResolve", MethodDocumentLinkResolve,
			func(h *Handler) {
				h.DocumentLinkResolve = func(c *glsp.Context, p *DocumentLink) (*DocumentLink, error) { return p, nil }
			}, `{"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}`},
		{"TextDocumentColor", MethodTextDocumentColor,
			func(h *Handler) {
				h.TextDocumentColor = func(c *glsp.Context, p *DocumentColorParams) ([]ColorInformation, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"TextDocumentColorPresentation", MethodTextDocumentColorPresentation,
			func(h *Handler) {
				h.TextDocumentColorPresentation = func(c *glsp.Context, p *ColorPresentationParams) ([]ColorPresentation, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"color":{"red":0,"green":0,"blue":0,"alpha":1},"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}`},
		{"TextDocumentFormatting", MethodTextDocumentFormatting,
			func(h *Handler) {
				h.TextDocumentFormatting = func(c *glsp.Context, p *DocumentFormattingParams) ([]TextEdit, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"options":{}}`},
		{"TextDocumentRangeFormatting", MethodTextDocumentRangeFormatting,
			func(h *Handler) {
				h.TextDocumentRangeFormatting = func(c *glsp.Context, p *DocumentRangeFormattingParams) ([]TextEdit, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"options":{}}`},
		{"TextDocumentOnTypeFormatting", MethodTextDocumentOnTypeFormatting,
			func(h *Handler) {
				h.TextDocumentOnTypeFormatting = func(c *glsp.Context, p *DocumentOnTypeFormattingParams) ([]TextEdit, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0},"ch":"{","options":{}}`},
		{"TextDocumentRename", MethodTextDocumentRename,
			func(h *Handler) {
				h.TextDocumentRename = func(c *glsp.Context, p *RenameParams) (*WorkspaceEdit, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0},"newName":"foo"}`},
		{"TextDocumentPrepareRename", MethodTextDocumentPrepareRename,
			func(h *Handler) {
				h.TextDocumentPrepareRename = func(c *glsp.Context, p *PrepareRenameParams) (any, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentFoldingRange", MethodTextDocumentFoldingRange,
			func(h *Handler) {
				h.TextDocumentFoldingRange = func(c *glsp.Context, p *FoldingRangeParams) ([]FoldingRange, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"TextDocumentSelectionRange", MethodTextDocumentSelectionRange,
			func(h *Handler) {
				h.TextDocumentSelectionRange = func(c *glsp.Context, p *SelectionRangeParams) ([]SelectionRange, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"positions":[]}`},
		{"TextDocumentLinkedEditingRange", MethodTextDocumentLinkedEditingRange,
			func(h *Handler) {
				h.TextDocumentLinkedEditingRange = func(c *glsp.Context, p *LinkedEditingRangeParams) (*LinkedEditingRanges, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"TextDocumentPrepareCallHierarchy", MethodTextDocumentPrepareCallHierarchy,
			func(h *Handler) {
				h.TextDocumentPrepareCallHierarchy = func(c *glsp.Context, p *CallHierarchyPrepareParams) ([]CallHierarchyItem, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"CallHierarchyIncomingCalls", MethodCallHierarchyIncomingCalls,
			func(h *Handler) {
				h.CallHierarchyIncomingCalls = func(c *glsp.Context, p *CallHierarchyIncomingCallsParams) ([]CallHierarchyIncomingCall, error) {
					return nil, nil
				}
			}, `{"item":{"name":"A","kind":12,"uri":"file:///a.go","range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"selectionRange":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}}`},
		{"CallHierarchyOutgoingCalls", MethodCallHierarchyOutgoingCalls,
			func(h *Handler) {
				h.CallHierarchyOutgoingCalls = func(c *glsp.Context, p *CallHierarchyOutgoingCallsParams) ([]CallHierarchyOutgoingCall, error) {
					return nil, nil
				}
			}, `{"item":{"name":"A","kind":12,"uri":"file:///a.go","range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"selectionRange":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}}`},
		{"TextDocumentSemanticTokensFull", MethodTextDocumentSemanticTokensFull,
			func(h *Handler) {
				h.TextDocumentSemanticTokensFull = func(c *glsp.Context, p *SemanticTokensParams) (*SemanticTokens, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"}}`},
		{"TextDocumentSemanticTokensFullDelta", MethodTextDocumentSemanticTokensFullDelta,
			func(h *Handler) {
				h.TextDocumentSemanticTokensFullDelta = func(c *glsp.Context, p *SemanticTokensDeltaParams) (any, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"previousResultId":"prev"}`},
		{"TextDocumentSemanticTokensRange", MethodTextDocumentSemanticTokensRange,
			func(h *Handler) {
				h.TextDocumentSemanticTokensRange = func(c *glsp.Context, p *SemanticTokensRangeParams) (any, error) {
					return nil, nil
				}
			}, `{"textDocument":{"uri":"file:///a.go"},"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}}}`},
		{"TextDocumentMoniker", MethodTextDocumentMoniker,
			func(h *Handler) {
				h.TextDocumentMoniker = func(c *glsp.Context, p *MonikerParams) ([]Moniker, error) { return nil, nil }
			}, `{"textDocument":{"uri":"file:///a.go"},"position":{"line":0,"character":0}}`},
		{"WorkspaceWillCreateFiles", MethodWorkspaceWillCreateFiles,
			func(h *Handler) {
				h.WorkspaceWillCreateFiles = func(c *glsp.Context, p *CreateFilesParams) (*WorkspaceEdit, error) { return nil, nil }
			}, `{"files":[]}`},
		{"WorkspaceDidCreateFiles", MethodWorkspaceDidCreateFiles,
			func(h *Handler) {
				h.WorkspaceDidCreateFiles = func(c *glsp.Context, p *CreateFilesParams) error { return nil }
			}, `{"files":[]}`},
		{"WorkspaceWillRenameFiles", MethodWorkspaceWillRenameFiles,
			func(h *Handler) {
				h.WorkspaceWillRenameFiles = func(c *glsp.Context, p *RenameFilesParams) (*WorkspaceEdit, error) { return nil, nil }
			}, `{"files":[]}`},
		{"WorkspaceDidRenameFiles", MethodWorkspaceDidRenameFiles,
			func(h *Handler) {
				h.WorkspaceDidRenameFiles = func(c *glsp.Context, p *RenameFilesParams) error { return nil }
			}, `{"files":[]}`},
		{"WorkspaceWillDeleteFiles", MethodWorkspaceWillDeleteFiles,
			func(h *Handler) {
				h.WorkspaceWillDeleteFiles = func(c *glsp.Context, p *DeleteFilesParams) (*WorkspaceEdit, error) { return nil, nil }
			}, `{"files":[]}`},
		{"WorkspaceDidDeleteFiles", MethodWorkspaceDidDeleteFiles,
			func(h *Handler) {
				h.WorkspaceDidDeleteFiles = func(c *glsp.Context, p *DeleteFilesParams) error { return nil }
			}, `{"files":[]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/nil handler", func(t *testing.T) {
			h := &Handler{}
			h.SetInitialized(true)
			_, vm, _, _ := h.Handle(ctx(tt.method, tt.validJSON))
			if vm {
				t.Error("validMethod should be false when handler is nil")
			}
		})
		t.Run(tt.name+"/valid", func(t *testing.T) {
			h := &Handler{}
			h.SetInitialized(true)
			tt.setup(h)
			_, vm, vp, err := h.Handle(ctx(tt.method, tt.validJSON))
			if err != nil {
				t.Fatal(err)
			}
			if !vm || !vp {
				t.Errorf("validMethod=%v validParams=%v", vm, vp)
			}
		})
		t.Run(tt.name+"/invalid params", func(t *testing.T) {
			h := &Handler{}
			h.SetInitialized(true)
			tt.setup(h)
			_, vm, vp, _ := h.Handle(ctx(tt.method, `{invalid`))
			if !vm {
				t.Error("validMethod should be true")
			}
			if vp {
				t.Error("validParams should be false for invalid JSON")
			}
		})
	}
}
