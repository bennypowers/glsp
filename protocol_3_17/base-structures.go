package protocol

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	protocol316 "github.com/tliron/glsp/protocol_3_16"
)

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#uri

type DocumentUri = string

type URI = string

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#regExp

/**
 * Client capabilities specific to regular expressions.
 */
type RegularExpressionsClientCapabilities struct {
	/**
	 * The engine's name.
	 */
	Engine string `json:"engine"`

	/**
	 * The engine's version.
	 */
	Version *string `json:"version,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocuments

var EOL = []string{"\n", "\r\n", "\r"}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#positionEncodingKind

/**
 * A set of predefined position encoding kinds.
 *
 * @since 3.17.0
 */
type PositionEncodingKind string

const (
	/**
	 * Character offsets count UTF-8 code units.
	 */
	PositionEncodingKindUTF8 = PositionEncodingKind("utf-8")

	/**
	 * Character offsets count UTF-16 code units.
	 *
	 * This is the default and must always be supported
	 * by servers
	 */
	PositionEncodingKindUTF16 = PositionEncodingKind("utf-16")

	/**
	 * Character offsets count UTF-32 code units.
	 *
	 * Implementation note: these are the same as Unicode code points,
	 * so this `PositionEncodingKind` may also be used for an
	 * encoding-agnostic representation of character offsets.
	 */
	PositionEncodingKindUTF32 = PositionEncodingKind("utf-32")
)

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#position

type Position struct {
	/**
	 * Line position in a document (zero-based).
	 */
	Line UInteger `json:"line"`

	/**
	 * Character offset on a line in a document (zero-based). The meaning of this
	 * offset is determined by the negotiated `PositionEncodingKind`.
	 *
	 * If the character value is greater than the line length it defaults back
	 * to the line length.
	 */
	Character UInteger `json:"character"`
}

func (self Position) IndexIn(content string) int {
	// This code is modified from the gopls implementation found:
	// https://cs.opensource.google/go/x/tools/+/refs/tags/v0.1.5:internal/span/utf16.go;l=70

	// In accordance with the LSP Spec:
	// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocuments
	// self.Character represents utf-16 code units by default, not bytes and so we need to
	// convert utf-16 code units to a byte offset.

	// Find the byte offset for the line
	index := 0
	for row := UInteger(0); row < self.Line; row++ {
		content_ := content[index:]
		if next := strings.Index(content_, "\n"); next != -1 {
			index += next + 1
		} else {
			return 0
		}
	}

	// The index represents the byte offset from the beginning of the line
	// count self.Character utf-16 code units from the index byte offset.

	byteOffset := index
	remains := content[index:]
	chr := int(self.Character)

	for count := 1; count <= chr; count++ {

		if len(remains) <= 0 {
			// char goes past content
			// this a error
			return 0
		}

		r, w := utf8.DecodeRuneInString(remains)
		if r == '\n' {
			// Per the LSP spec:
			//
			// > If the character value is greater than the line length it
			// > defaults back to the line length.
			break
		}

		remains = remains[w:]
		if r >= 0x10000 {
			// a two point rune
			count++
			// if we finished in a two point rune, do not advance past the first
			if count > chr {
				break
			}
		}
		byteOffset += w

	}

	return byteOffset
}

func (self Position) EndOfLineIn(content string) Position {
	index := self.IndexIn(content)
	content_ := content[index:]
	if eol := strings.Index(content_, "\n"); eol != -1 {
		return Position{
			Line:      self.Line,
			Character: self.Character + UInteger(eol),
		}
	} else {
		return self
	}
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#range

type Range struct {
	/**
	 * The range's start position.
	 */
	Start Position `json:"start"`

	/**
	 * The range's end position.
	 */
	End Position `json:"end"`
}

func (self Range) IndexesIn(content string) (int, int) {
	return self.Start.IndexIn(content), self.End.IndexIn(content)
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#location

type Location struct {
	URI   DocumentUri `json:"uri"`
	Range Range       `json:"range"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#locationLink

type LocationLink struct {
	/**
	 * Span of the origin of this link.
	 *
	 * Used as the underlined span for mouse interaction. Defaults to the word
	 * range at the mouse position.
	 */
	OriginSelectionRange *Range `json:"originSelectionRange,omitempty"`

	/**
	 * The target resource identifier of this link.
	 */
	TargetURI DocumentUri `json:"targetUri"`

	/**
	 * The full target range of this link. If the target for example is a symbol
	 * then target range is the range enclosing this symbol not including
	 * leading/trailing whitespace but everything else like comments. This
	 * information is typically used to highlight the range in the editor.
	 */
	TargetRange Range `json:"targetRange"`

	/**
	 * The range that should be selected and revealed when this link is being
	 * followed, e.g the name of a function. Must be contained by the the
	 * `targetRange`. See also `DocumentSymbol#range`
	 */
	TargetSelectionRange Range `json:"targetSelectionRange"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#diagnostic

type Diagnostic struct {
	/**
	 * The range at which the message applies.
	 */
	Range Range `json:"range"`

	/**
	 * The diagnostic's severity. Can be omitted. If omitted it is up to the
	 * client to interpret diagnostics as error, warning, info or hint.
	 */
	Severity *DiagnosticSeverity `json:"severity,omitempty"`

	/**
	 * The diagnostic's code, which might appear in the user interface.
	 */
	Code *IntegerOrString `json:"code,omitempty"`

	/**
	 * An optional property to describe the error code.
	 *
	 * @since 3.16.0
	 */
	CodeDescription *CodeDescription `json:"codeDescription,omitempty"`

	/**
	 * A human-readable string describing the source of this
	 * diagnostic, e.g. 'typescript' or 'super lint'.
	 */
	Source *string `json:"source,omitempty"`

	/**
	 * The diagnostic's message.
	 */
	Message string `json:"message"`

	/**
	 * Additional metadata about the diagnostic.
	 *
	 * @since 3.15.0
	 */
	Tags []DiagnosticTag `json:"tags,omitempty"`

	/**
	 * An array of related diagnostic information, e.g. when symbol-names within
	 * a scope collide all definitions can be marked via this property.
	 */
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`

	/**
	 * A data entry field that is preserved between a
	 * `textDocument/publishDiagnostics` notification and
	 * `textDocument/codeAction` request.
	 *
	 * @since 3.16.0
	 */
	Data any `json:"data,omitempty"`
}

type DiagnosticSeverity Integer

const (
	/**
	 * Reports an error.
	 */
	DiagnosticSeverityError = DiagnosticSeverity(1)

	/**
	 * Reports a warning.
	 */
	DiagnosticSeverityWarning = DiagnosticSeverity(2)

	/**
	 * Reports an information.
	 */
	DiagnosticSeverityInformation = DiagnosticSeverity(3)

	/**
	 * Reports a hint.
	 */
	DiagnosticSeverityHint = DiagnosticSeverity(4)
)

/**
 * The diagnostic tags.
 *
 * @since 3.15.0
 */
type DiagnosticTag Integer

const (
	/**
	 * Unused or unnecessary code.
	 *
	 * Clients are allowed to render diagnostics with this tag faded out
	 * instead of having an error squiggle.
	 */
	DiagnosticTagUnnecessary = DiagnosticTag(1)

	/**
	 * Deprecated or obsolete code.
	 *
	 * Clients are allowed to rendered diagnostics with this tag strike through.
	 */
	DiagnosticTagDeprecated = DiagnosticTag(2)
)

/**
 * Represents a related message and source code location for a diagnostic.
 * This should be used to point to code locations that cause or are related to
 * a diagnostics, e.g when duplicating a symbol in a scope.
 */
type DiagnosticRelatedInformation struct {
	/**
	 * The location of this related diagnostic information.
	 */
	Location Location `json:"location"`

	/**
	 * The message of this related diagnostic information.
	 */
	Message string `json:"message"`
}

/**
 * Structure to capture a description for an error code.
 *
 * @since 3.16.0
 */
type CodeDescription struct {
	/**
	 * An URI to open with more information about the diagnostic error.
	 */
	HRef URI `json:"href"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#command

type Command struct {
	/**
	 * Title of the command, like `save`.
	 */
	Title string `json:"title"`

	/**
	 * The identifier of the actual command handler.
	 */
	Command string `json:"command"`

	/**
	 * Arguments that the command handler should be
	 * invoked with.
	 */
	Arguments []any `json:"arguments,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textEdit

type TextEdit struct {
	/**
	 * The range of the text document to be manipulated. To insert
	 * text into a document create a range where start === end.
	 */
	Range Range `json:"range"`

	/**
	 * The string to be inserted. For delete operations use an
	 * empty string.
	 */
	NewText string `json:"newText"`
}

/**
 * Additional information that describes document changes.
 *
 * @since 3.16.0
 */
type ChangeAnnotation struct {
	/**
	 * A human-readable string describing the actual change. The string
	 * is rendered prominent in the user interface.
	 */
	Label string `json:"label"`

	/**
	 * A flag which indicates that user confirmation is needed
	 * before applying the change.
	 */
	NeedsConfirmation *bool `json:"needsConfirmation,omitempty"`

	/**
	 * A human-readable string which is rendered less prominent in
	 * the user interface.
	 */
	Description *string `json:"description,omitempty"`
}

/**
 * An identifier referring to a change annotation managed by a workspace
 * edit.
 *
 * @since 3.16.0
 */
type ChangeAnnotationIdentifier = string

/**
 * A special text edit with an additional change annotation.
 *
 * @since 3.16.0
 */
type AnnotatedTextEdit struct {
	TextEdit

	/**
	 * The actual annotation identifier.
	 */
	AnnotationID ChangeAnnotationIdentifier `json:"annotationId"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocumentEdit

type TextDocumentEdit struct {
	/**
	 * The text document to change.
	 */
	TextDocument OptionalVersionedTextDocumentIdentifier `json:"textDocument"`

	/**
	 * The edits to be applied.
	 *
	 * @since 3.16.0 - support for AnnotatedTextEdit. This is guarded by the
	 * client capability `workspace.workspaceEdit.changeAnnotationSupport`
	 */
	Edits []any `json:"edits"` // TextEdit | AnnotatedTextEdit
}

// ([json.Unmarshaler] interface)
func (self *TextDocumentEdit) UnmarshalJSON(data []byte) error {
	var value struct {
		TextDocument OptionalVersionedTextDocumentIdentifier `json:"textDocument"`
		Edits        []json.RawMessage                       `json:"edits"` // TextEdit | AnnotatedTextEdit
	}

	if err := json.Unmarshal(data, &value); err == nil {
		self.TextDocument = value.TextDocument

		for _, edit := range value.Edits {
			var value TextEdit
			if err = json.Unmarshal(edit, &value); err == nil {
				self.Edits = append(self.Edits, value)
			} else {
				var value AnnotatedTextEdit
				if err = json.Unmarshal(edit, &value); err == nil {
					self.Edits = append(self.Edits, value)
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

// Re-use file operations from 3.16 as they are unchanged
type CreateFileOptions = protocol316.CreateFileOptions
type CreateFile = protocol316.CreateFile
type RenameFileOptions = protocol316.RenameFileOptions
type RenameFile = protocol316.RenameFile
type DeleteFileOptions = protocol316.DeleteFileOptions
type DeleteFile = protocol316.DeleteFile

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#workspaceEdit

type WorkspaceEdit struct {
	/**
	 * Holds changes to existing resources.
	 */
	Changes map[DocumentUri][]TextEdit `json:"changes,omitempty"`

	/**
	 * Depending on the client capability
	 * `workspace.workspaceEdit.resourceOperations` document changes are either
	 * an array of `TextDocumentEdit`s to express changes to n different text
	 * documents where each text document edit addresses a specific version of
	 * a text document. Or it can contain above `TextDocumentEdit`s mixed with
	 * create, rename and delete file / folder operations.
	 *
	 * Whether a client supports versioned document edits is expressed via
	 * `workspace.workspaceEdit.documentChanges` client capability.
	 *
	 * If a client neither supports `documentChanges` nor
	 * `workspace.workspaceEdit.resourceOperations` then only plain `TextEdit`s
	 * using the `changes` property are supported.
	 */
	DocumentChanges []any `json:"documentChanges,omitempty"` // TextDocumentEdit | CreateFile | RenameFile | DeleteFile

	/**
	 * A map of change annotations that can be referenced in
	 * `AnnotatedTextEdit`s or create, rename and delete file / folder
	 * operations.
	 *
	 * Whether clients honor this property depends on the client capability
	 * `workspace.changeAnnotationSupport`.
	 *
	 * @since 3.16.0
	 */
	ChangeAnnotations map[ChangeAnnotationIdentifier]ChangeAnnotation `json:"changeAnnotations,omitempty"`
}

// ([json.Unmarshaler] interface)
func (self *WorkspaceEdit) UnmarshalJSON(data []byte) error {
	var value struct {
		Changes           map[DocumentUri][]TextEdit                      `json:"changes"`
		DocumentChanges   []json.RawMessage                               `json:"documentChanges"` // TextDocumentEdit | CreateFile | RenameFile | DeleteFile
		ChangeAnnotations map[ChangeAnnotationIdentifier]ChangeAnnotation `json:"changeAnnotations"`
	}

	if err := json.Unmarshal(data, &value); err == nil {
		self.Changes = value.Changes
		self.ChangeAnnotations = value.ChangeAnnotations

		for _, documentChange := range value.DocumentChanges {
			var value TextDocumentEdit
			if err = json.Unmarshal(documentChange, &value); err == nil {
				self.DocumentChanges = append(self.DocumentChanges, value)
			} else {
				var value CreateFile
				if err = json.Unmarshal(documentChange, &value); err == nil {
					self.DocumentChanges = append(self.DocumentChanges, value)
				} else {
					var value RenameFile
					if err = json.Unmarshal(documentChange, &value); err == nil {
						self.DocumentChanges = append(self.DocumentChanges, value)
					} else {
						var value DeleteFile
						if err = json.Unmarshal(documentChange, &value); err == nil {
							self.DocumentChanges = append(self.DocumentChanges, value)
						} else {
							return err
						}
					}
				}
			}
		}

		return nil
	} else {
		return err
	}
}

// Re-use capabilities from 3.16 where applicable
type WorkspaceEditClientCapabilities = protocol316.WorkspaceEditClientCapabilities
type ResourceOperationKind = protocol316.ResourceOperationKind
type FailureHandlingKind = protocol316.FailureHandlingKind

const (
	ResourceOperationKindCreate = protocol316.ResourceOperationKindCreate
	ResourceOperationKindRename = protocol316.ResourceOperationKindRename
	ResourceOperationKindDelete = protocol316.ResourceOperationKindDelete
)

const (
	FailureHandlingKindAbort                 = protocol316.FailureHandlingKindAbort
	FailureHandlingKindTransactional         = protocol316.FailureHandlingKindTransactional
	FailureHandlingKindTextOnlyTransactional = protocol316.FailureHandlingKindTextOnlyTransactional
	FailureHandlingKindUndo                  = protocol316.FailureHandlingKindUndo
)

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocumentIdentifier

type TextDocumentIdentifier struct {
	/**
	 * The text document's URI.
	 */
	URI DocumentUri `json:"uri"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocumentItem

type TextDocumentItem struct {
	/**
	 * The text document's URI.
	 */
	URI DocumentUri `json:"uri"`

	/**
	 * The text document's language identifier.
	 */
	LanguageID string `json:"languageId"`

	/**
	 * The version number of this document (it will increase after each
	 * change, including undo/redo).
	 */
	Version Integer `json:"version"`

	/**
	 * The content of the opened text document.
	 */
	Text string `json:"text"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#versionedTextDocumentIdentifier

type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier

	/**
	 * The version number of this document.
	 *
	 * The version number of a document will increase after each change,
	 * including undo/redo. The number doesn't need to be consecutive.
	 */
	Version Integer `json:"version"`
}

type OptionalVersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier

	/**
	 * The version number of this document. If an optional versioned text document
	 * identifier is sent from the server to the client and the file is not
	 * open in the editor (the server has not received an open notification
	 * before) the server can send `null` to indicate that the version is
	 * known and the content on disk is the master (as specified with document
	 * content ownership).
	 *
	 * The version number of a document will increase after each change,
	 * including undo/redo. The number doesn't need to be consecutive.
	 */
	Version *Integer `json:"version"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocumentPositionParams

type TextDocumentPositionParams struct {
	/**
	 * The text document.
	 */
	TextDocument TextDocumentIdentifier `json:"textDocument"`

	/**
	 * The position inside the text document.
	 */
	Position Position `json:"position"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#documentFilter

type DocumentFilter struct {
	/**
	 * A language id, like `typescript`.
	 */
	Language *string `json:"language,omitempty"`

	/**
	 * A Uri [scheme](#Uri.scheme), like `file` or `untitled`.
	 */
	Scheme *string `json:"scheme,omitempty"`

	/**
	 * A glob pattern, like `*.{ts,js}`.
	 *
	 * Glob patterns can have the following syntax:
	 * - `*` to match one or more characters in a path segment
	 * - `?` to match on one character in a path segment
	 * - `**` to match any number of path segments, including none
	 * - `{}` to group conditions (e.g. `**​/*.{ts,js}` matches all TypeScript
	 *   and JavaScript files)
	 * - `[]` to declare a range of characters to match in a path segment
	 *   (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
	 * - `[!...]` to negate a range of characters to match in a path segment
	 *   (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but
	 *   not `example.0`)
	 */
	Pattern *string `json:"pattern,omitempty"`
}

type DocumentSelector []DocumentFilter

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#staticRegistrationOptions

/**
 * Static registration options to be returned in the initialize request.
 */
type StaticRegistrationOptions struct {
	/**
	 * The id used to register the request. The id can be used to deregister
	 * the request again. See also Registration#id.
	 */
	ID *string `json:"id,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#textDocumentRegistrationOptions

/**
 * General text document registration options.
 */
type TextDocumentRegistrationOptions struct {
	/**
	 * A document selector to identify the scope of the registration. If set to
	 * null the document selector provided on the client side will be used.
	 */
	DocumentSelector *DocumentSelector `json:"documentSelector"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#markupContent

/**
 * Describes the content type that a client supports in various
 * result literals like `Hover`, `ParameterInfo` or `CompletionItem`.
 *
 * Please note that `MarkupKinds` must not start with a `$`. This kinds
 * are reserved for internal usage.
 */
type MarkupKind string

const (
	/**
	 * Plain text is supported as a content format
	 */
	MarkupKindPlainText = MarkupKind("plaintext")

	/**
	 * Markdown is supported as a content format
	 */
	MarkupKindMarkdown = MarkupKind("markdown")
)

/**
 * A `MarkupContent` literal represents a string value which content is
 * interpreted base on its kind flag. Currently the protocol supports
 * `plaintext` and `markdown` as markup kinds.
 *
 * If the kind is `markdown` then the value can contain fenced code blocks like
 * in GitHub issues.
 *
 * Here is an example how such a string can be constructed using
 * JavaScript / TypeScript:
 * ```typescript
 * let markdown: MarkdownContent = {
 *  kind: MarkupKind.Markdown,
 *	value: [
 *		'# Header',
 *		'Some text',
 *		'```typescript',
 *		'someCode();',
 *		'```'
 *	].join('\n')
 * };
 * ```
 *
 * *Please Note* that clients might sanitize the return markdown. A client could
 * decide to remove HTML from the markdown to avoid script execution.
 */
type MarkupContent struct {
	/**
	 * The type of the Markup
	 */
	Kind MarkupKind `json:"kind"`

	/**
	 * The content itself
	 */
	Value string `json:"value"`
}

/**
 * Client capabilities specific to the used markdown parser.
 *
 * @since 3.16.0
 */
type MarkdownClientCapabilities struct {
	/**
	 * The name of the parser.
	 */
	Parser string `json:"parser"`

	/**
	 * The version of the parser.
	 */
	Version *string `json:"version,omitempty"`

	/**
	 * A list of HTML tags that the client allows / supports in
	 * Markdown.
	 *
	 * @since 3.17.0
	 */
	AllowedTags []string `json:"allowedTags,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#workDoneProgress

type WorkDoneProgressBegin struct {
	Kind string `json:"kind"` // == "begin"

	/**
	 * Mandatory title of the progress operation. Used to briefly inform about
	 * the kind of operation being performed.
	 *
	 * Examples: "Indexing" or "Linking dependencies".
	 */
	Title string `json:"title"`

	/**
	 * Controls if a cancel button should show to allow the user to cancel the
	 * long running operation. Clients that don't support cancellation are
	 * allowed to ignore the setting.
	 */
	Cancellable *bool `json:"cancellable,omitempty"`

	/**
	 * Optional, more detailed associated progress message. Contains
	 * complementary information to the `title`.
	 *
	 * Examples: "3/25 files", "project/src/module2", "node_modules/some_dep".
	 * If unset, the previous progress message (if any) is still valid.
	 */
	Message *string `json:"message,omitempty"`

	/**
	 * Optional progress percentage to display (value 100 is considered 100%).
	 * If not provided infinite progress is assumed and clients are allowed
	 * to ignore the `percentage` value in subsequent in report notifications.
	 *
	 * The value should be steadily rising. Clients are free to ignore values
	 * that are not following this rule. The value range is [0, 100]
	 */
	Percentage *UInteger `json:"percentage,omitempty"`
}

type WorkDoneProgressReport struct {
	Kind string `json:"kind"` // == "report"

	/**
	 * Controls enablement state of a cancel button. This property is only valid
	 *  if a cancel button got requested in the `WorkDoneProgressStart` payload.
	 *
	 * Clients that don't support cancellation or don't support control the
	 * button's enablement state are allowed to ignore the setting.
	 */
	Cancellable *bool `json:"cancellable,omitempty"`

	/**
	 * Optional, more detailed associated progress message. Contains
	 * complementary information to the `title`.
	 *
	 * Examples: "3/25 files", "project/src/module2", "node_modules/some_dep".
	 * If unset, the previous progress message (if any) is still valid.
	 */
	Message *string `json:"message,omitempty"`

	/**
	 * Optional progress percentage to display (value 100 is considered 100%).
	 * If not provided infinite progress is assumed and clients are allowed
	 * to ignore the `percentage` value in subsequent in report notifications.
	 *
	 * The value should be steadily rising. Clients are free to ignore values
	 * that are not following this rule. The value range is [0, 100]
	 */
	Percentage *UInteger `json:"percentage,omitempty"`
}

type WorkDoneProgressEnd struct {
	Kind string `json:"kind"` // == "end"

	/**
	 * Optional, a final message indicating to for example indicate the outcome
	 * of the operation.
	 */
	Message *string `json:"message,omitempty"`
}

type WorkDoneProgressParams struct {
	/**
	 * An optional token that a server can use to report work done progress.
	 */
	WorkDoneToken *ProgressToken `json:"workDoneToken,omitempty"`
}

type WorkDoneProgressOptions struct {
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#partialResults

type PartialResultParams struct {
	/**
	 * An optional token that a server can use to report partial results (e.g.
	 * streaming) to the client.
	 */
	PartialResultToken *ProgressToken `json:"partialResultToken,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#traceValue

type TraceValue string

const (
	TraceValueOff     = TraceValue("off")
	TraceValueMessage = TraceValue("message") // The spec clearly says "message", but some implementations use "messages" instead
	TraceValueVerbose = TraceValue("verbose")
)
