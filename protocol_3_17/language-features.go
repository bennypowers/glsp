package protocol

import (
	"encoding/json"

	"github.com/bennypowers/glsp"
)

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_implementation

// Bug fix: 3.16 incorrectly embedded TypeDefinitionOptions instead of ImplementationOptions
type ImplementationRegistrationOptions struct {
	TextDocumentRegistrationOptions
	ImplementationOptions
	StaticRegistrationOptions
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_semanticTokens

// Bug fix: Requests.Range json tag changed from "Range" to "range"
type SemanticTokensClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`

	Requests struct {
		Range any `json:"range,omitempty"` // nil | bool | struct{}
		Full  any `json:"full,omitempty"`  // nil | bool | SemanticDelta
	} `json:"requests"`

	TokenTypes              []string      `json:"tokenTypes"`
	TokenModifiers          []string      `json:"tokenModifiers"`
	Formats                 []TokenFormat `json:"formats"`
	OverlappingTokenSupport *bool         `json:"overlappingTokenSupport,omitempty"`
	MultilineTokenSupport   *bool         `json:"multilineTokenSupport,omitempty"`
}

// ([json.Unmarshaler] interface)
func (self *SemanticTokensClientCapabilities) UnmarshalJSON(data []byte) error {
	var value struct {
		DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
		Requests            struct {
			Range json.RawMessage `json:"range,omitempty"` // nil | bool | struct{}
			Full  json.RawMessage `json:"full,omitempty"`  // nil | bool | SemanticDelta
		} `json:"requests"`
		TokenTypes              []string      `json:"tokenTypes"`
		TokenModifiers          []string      `json:"tokenModifiers"`
		Formats                 []TokenFormat `json:"formats"`
		OverlappingTokenSupport *bool         `json:"overlappingTokenSupport,omitempty"`
		MultilineTokenSupport   *bool         `json:"multilineTokenSupport,omitempty"`
	}

	if err := json.Unmarshal(data, &value); err == nil {
		self.DynamicRegistration = value.DynamicRegistration
		self.TokenTypes = value.TokenTypes
		self.TokenModifiers = value.TokenModifiers
		self.Formats = value.Formats
		self.OverlappingTokenSupport = value.OverlappingTokenSupport
		self.MultilineTokenSupport = value.MultilineTokenSupport

		if value.Requests.Range != nil {
			var value_ bool
			if err = json.Unmarshal(value.Requests.Range, &value_); err == nil {
				self.Requests.Range = value_
			} else {
				var value_ struct{}
				if err = json.Unmarshal(value.Requests.Range, &value_); err == nil {
					self.Requests.Range = value_
				} else {
					return err
				}
			}
		}

		if value.Requests.Full != nil {
			var value_ bool
			if err = json.Unmarshal(value.Requests.Full, &value_); err == nil {
				self.Requests.Full = value_
			} else {
				var value_ SemanticDelta
				if err = json.Unmarshal(value.Requests.Full, &value_); err == nil {
					self.Requests.Full = value_
				} else {
					return err
				}
			}
		}

		return nil
	} else {
		return err
	}
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_prepareTypeHierarchy

/**
 * @since 3.17.0
 */
type TypeHierarchyClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
}

type TypeHierarchyOptions struct {
	WorkDoneProgressOptions
}

type TypeHierarchyRegistrationOptions struct {
	TextDocumentRegistrationOptions
	TypeHierarchyOptions
	StaticRegistrationOptions
}

const MethodTextDocumentPrepareTypeHierarchy = Method("textDocument/prepareTypeHierarchy")

type TextDocumentPrepareTypeHierarchyFunc func(context *glsp.Context, params *TypeHierarchyPrepareParams) ([]TypeHierarchyItem, error)

type TypeHierarchyPrepareParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

/**
 * @since 3.17.0
 */
type TypeHierarchyItem struct {
	Name           string      `json:"name"`
	Kind           SymbolKind  `json:"kind"`
	Tags           []SymbolTag `json:"tags,omitempty"`
	Detail         *string     `json:"detail,omitempty"`
	URI            DocumentUri `json:"uri"`
	Range          Range       `json:"range"`
	SelectionRange Range       `json:"selectionRange"`
	Data           any         `json:"data,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#typeHierarchy_supertypes

const MethodTypeHierarchySupertypes = Method("typeHierarchy/supertypes")

type TypeHierarchySupertypesFunc func(context *glsp.Context, params *TypeHierarchySupertypesParams) ([]TypeHierarchyItem, error)

type TypeHierarchySupertypesParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item TypeHierarchyItem `json:"item"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#typeHierarchy_subtypes

const MethodTypeHierarchySubtypes = Method("typeHierarchy/subtypes")

type TypeHierarchySubtypesFunc func(context *glsp.Context, params *TypeHierarchySubtypesParams) ([]TypeHierarchyItem, error)

type TypeHierarchySubtypesParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item TypeHierarchyItem `json:"item"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_inlineValue

/**
 * @since 3.17.0
 */
type InlineValueClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
}

type InlineValueOptions struct {
	WorkDoneProgressOptions
}

type InlineValueRegistrationOptions struct {
	InlineValueOptions
	TextDocumentRegistrationOptions
	StaticRegistrationOptions
}

const MethodTextDocumentInlineValue = Method("textDocument/inlineValue")

type TextDocumentInlineValueFunc func(context *glsp.Context, params *InlineValueParams) ([]InlineValue, error)

type InlineValueParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      InlineValueContext     `json:"context"`
}

type InlineValueContext struct {
	FrameId         Integer `json:"frameId"`
	StoppedLocation Range   `json:"stoppedLocation"`
}

type InlineValue any // InlineValueText | InlineValueVariableLookup | InlineValueEvaluatableExpression

type InlineValueText struct {
	Range Range  `json:"range"`
	Text  string `json:"text"`
}

type InlineValueVariableLookup struct {
	Range               Range   `json:"range"`
	VariableName        *string `json:"variableName,omitempty"`
	CaseSensitiveLookup bool    `json:"caseSensitiveLookup"`
}

type InlineValueEvaluatableExpression struct {
	Range      Range   `json:"range"`
	Expression *string `json:"expression,omitempty"`
}

const ServerWorkspaceInlineValueRefresh = Method("workspace/inlineValue/refresh")

type InlineValueWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocument_inlayHint

type InlayHintClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	ResolveSupport      *struct {
		Properties []string `json:"properties"`
	} `json:"resolveSupport,omitempty"`
}

type InlayHintOptions struct {
	WorkDoneProgressOptions
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

type InlayHintRegistrationOptions struct {
	InlayHintOptions
	TextDocumentRegistrationOptions
	StaticRegistrationOptions
}

const MethodTextDocumentInlayHint = Method("textDocument/inlayHint")

type TextDocumentInlayHintFunc func(context *glsp.Context, params *InlayHintParams) ([]InlayHint, error)

type InlayHintParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
}

type InlayHint struct {
	Position     Position       `json:"position"`
	Label        any            `json:"label"` // string | InlayHintLabelPart[]
	Kind         *InlayHintKind `json:"kind,omitempty"`
	TextEdits    []TextEdit     `json:"textEdits,omitempty"`
	Tooltip      any            `json:"tooltip,omitempty"` // string | MarkupContent
	PaddingLeft  *bool          `json:"paddingLeft,omitempty"`
	PaddingRight *bool          `json:"paddingRight,omitempty"`
	Data         any            `json:"data,omitempty"`
}

type InlayHintLabelPart struct {
	Value    string    `json:"value"`
	Tooltip  any       `json:"tooltip,omitempty"` // string | MarkupContent
	Location *Location `json:"location,omitempty"`
	Command  *Command  `json:"command,omitempty"`
}

type InlayHintKind Integer

const (
	InlayHintKindType      = InlayHintKind(1)
	InlayHintKindParameter = InlayHintKind(2)
)

const MethodInlayHintResolve = Method("inlayHint/resolve")

type InlayHintResolveFunc func(context *glsp.Context, params *InlayHint) (*InlayHint, error)

const ServerWorkspaceInlayHintRefresh = Method("workspace/inlayHint/refresh")

type InlayHintWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
