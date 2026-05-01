package protocol

import "github.com/bennypowers/glsp"

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notebookDocument_synchronization

/**
 * A notebook document.
 *
 * @since 3.17.0
 */
type NotebookDocument struct {
	/**
	 * The notebook document's URI.
	 */
	URI DocumentUri `json:"uri"`

	/**
	 * The type of the notebook.
	 */
	NotebookType string `json:"notebookType"`

	/**
	 * The version number of this document (it will increase after each
	 * change, including undo/redo).
	 */
	Version Integer `json:"version"`

	/**
	 * Additional metadata stored with the notebook
	 * document.
	 */
	Metadata *any `json:"metadata,omitempty"`

	/**
	 * The cells of a notebook.
	 */
	Cells []NotebookCell `json:"cells"`
}

/**
 * A notebook cell.
 *
 * A cell's document URI must be unique across ALL notebook cells and can therefore
 * be used to uniquely identify a notebook cell or the cell's text document.
 *
 * @since 3.17.0
 */
type NotebookCell struct {
	/**
	 * The cell's kind
	 */
	Kind NotebookCellKind `json:"kind"`

	/**
	 * The URI of the cell's text document
	 * content.
	 */
	Document DocumentUri `json:"document"`

	/**
	 * Additional metadata stored with the cell.
	 */
	Metadata *any `json:"metadata,omitempty"`

	/**
	 * Additional execution summary information
	 * if supported by the client.
	 */
	ExecutionSummary *ExecutionSummary `json:"executionSummary,omitempty"`
}

/**
 * A notebook cell kind.
 *
 * @since 3.17.0
 */
type NotebookCellKind Integer

const (
	/**
	 * A markup-cell is formatted source that is used for display.
	 */
	NotebookCellKindMarkup = NotebookCellKind(1)

	/**
	 * A code-cell is source code.
	 */
	NotebookCellKindCode = NotebookCellKind(2)
)

type ExecutionSummary struct {
	/**
	 * The execution order number.
	 */
	ExecutionOrder UInteger `json:"executionOrder"`

	/**
	 * Whether the execution was successful or
	 * not if known by the client.
	 */
	Success *bool `json:"success,omitempty"`
}

/**
 * A versioned notebook document identifier.
 *
 * @since 3.17.0
 */
type VersionedNotebookDocumentIdentifier struct {
	/**
	 * The version number of this notebook document.
	 */
	Version Integer `json:"version"`

	/**
	 * The notebook document's URI.
	 */
	URI DocumentUri `json:"uri"`
}

/**
 * A change describing how to move a `NotebookCell`
 * array from state S to S'.
 *
 * @since 3.17.0
 */
type NotebookCellArrayChange struct {
	/**
	 * The start offset of the cell that changed.
	 */
	Start UInteger `json:"start"`

	/**
	 * The deleted cells
	 */
	DeleteCount UInteger `json:"deleteCount"`

	/**
	 * The new cells, if any
	 */
	Cells []NotebookCell `json:"cells,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notebookDocument_didOpen

const MethodNotebookDocumentDidOpen = Method("notebookDocument/didOpen")

type NotebookDocumentDidOpenFunc func(context *glsp.Context, params *DidOpenNotebookDocumentParams) error

/**
 * The params sent in a open notebook document notification.
 *
 * @since 3.17.0
 */
type DidOpenNotebookDocumentParams struct {
	/**
	 * The notebook document that got opened.
	 */
	NotebookDocument NotebookDocument `json:"notebookDocument"`

	/**
	 * The text documents that represent the content
	 * of a notebook cell.
	 */
	CellTextDocuments []TextDocumentItem `json:"cellTextDocuments"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notebookDocument_didChange

const MethodNotebookDocumentDidChange = Method("notebookDocument/didChange")

type NotebookDocumentDidChangeFunc func(context *glsp.Context, params *DidChangeNotebookDocumentParams) error

/**
 * The params sent in a change notebook document notification.
 *
 * @since 3.17.0
 */
type DidChangeNotebookDocumentParams struct {
	/**
	 * The notebook document that did change. The version number points
	 * to the version after all provided changes have been applied.
	 */
	NotebookDocument VersionedNotebookDocumentIdentifier `json:"notebookDocument"`

	/**
	 * The actual changes to the notebook document.
	 *
	 * The change describes single state change to the notebook document.
	 * So it moves a notebook document, its cells and its cell text document
	 * contents from state S to S'.
	 *
	 * To mirror the content of a notebook using change events use the
	 * following approach:
	 * - start with the same initial content
	 * - apply the 'notebookDocument/didChange' notifications in the order
	 *   you receive them.
	 */
	Change NotebookDocumentChangeEvent `json:"change"`
}

/**
 * A change event for a notebook document.
 *
 * @since 3.17.0
 */
type NotebookDocumentChangeEvent struct {
	/**
	 * The changed meta data if any.
	 */
	Metadata *any `json:"metadata,omitempty"`

	/**
	 * Changes to cells
	 */
	Cells *struct {
		/**
		 * Changes to the cell structure to add or
		 * remove cells.
		 */
		Structure *struct {
			/**
			 * The change to the cell array.
			 */
			Array NotebookCellArrayChange `json:"array"`

			/**
			 * Additional opened cell text documents.
			 */
			DidOpen []TextDocumentItem `json:"didOpen,omitempty"`

			/**
			 * Additional closed cell text documents.
			 */
			DidClose []TextDocumentIdentifier `json:"didClose,omitempty"`
		} `json:"structure,omitempty"`

		/**
		 * Changes to notebook cells properties like its
		 * kind, execution summary or metadata.
		 */
		Data []NotebookCell `json:"data,omitempty"`

		/**
		 * Changes to the text content of notebook cells.
		 */
		TextContent []struct {
			Document VersionedTextDocumentIdentifier  `json:"document"`
			Changes  []TextDocumentContentChangeEvent `json:"changes"`
		} `json:"textContent,omitempty"`
	} `json:"cells,omitempty"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notebookDocument_didSave

const MethodNotebookDocumentDidSave = Method("notebookDocument/didSave")

type NotebookDocumentDidSaveFunc func(context *glsp.Context, params *DidSaveNotebookDocumentParams) error

/**
 * The params sent in a save notebook document notification.
 *
 * @since 3.17.0
 */
type DidSaveNotebookDocumentParams struct {
	/**
	 * The notebook document that got saved.
	 */
	NotebookDocument NotebookDocumentIdentifier `json:"notebookDocument"`
}

/**
 * A literal to identify a notebook document in the client.
 *
 * @since 3.17.0
 */
type NotebookDocumentIdentifier struct {
	/**
	 * The notebook document's URI.
	 */
	URI DocumentUri `json:"uri"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notebookDocument_didClose

const MethodNotebookDocumentDidClose = Method("notebookDocument/didClose")

type NotebookDocumentDidCloseFunc func(context *glsp.Context, params *DidCloseNotebookDocumentParams) error

/**
 * The params sent in a close notebook document notification.
 *
 * @since 3.17.0
 */
type DidCloseNotebookDocumentParams struct {
	/**
	 * The notebook document that got closed.
	 */
	NotebookDocument NotebookDocumentIdentifier `json:"notebookDocument"`

	/**
	 * The text documents that represent the content
	 * of a notebook cell that got closed.
	 */
	CellTextDocuments []TextDocumentIdentifier `json:"cellTextDocuments"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#notebookDocument_synchronization

/**
 * Notebook specific client capabilities.
 *
 * @since 3.17.0
 */
type NotebookDocumentClientCapabilities struct {
	/**
	 * Capabilities specific to notebook document synchronization
	 *
	 * @since 3.17.0
	 */
	Synchronization NotebookDocumentSyncClientCapabilities `json:"synchronization"`
}

/**
 * Notebook specific synchronization client capabilities.
 *
 * @since 3.17.0
 */
type NotebookDocumentSyncClientCapabilities struct {
	/**
	 * Whether implementation supports dynamic registration. If this is
	 * set to `true` the client supports the new
	 * `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
	 * return value for the corresponding server capability as well.
	 */
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`

	/**
	 * The client supports sending execution summary data per cell.
	 */
	ExecutionSummarySupport *bool `json:"executionSummarySupport,omitempty"`
}

/**
 * Options specific to a notebook plus its associated text documents.
 *
 * @since 3.17.0
 */
type NotebookDocumentSyncOptions struct {
	/**
	 * The notebooks to be synced
	 */
	NotebookSelector []NotebookSelector `json:"notebookSelector"`

	/**
	 * Whether save notification should be forwarded to
	 * the server. Will only be honored if mode === `notebook`.
	 */
	Save *bool `json:"save,omitempty"`
}

/**
 * Registration options specific to a notebook.
 *
 * @since 3.17.0
 */
type NotebookDocumentSyncRegistrationOptions struct {
	NotebookDocumentSyncOptions
	StaticRegistrationOptions
}

/**
 * A notebook selector is the combination of one or two filters:
 * - `notebook`: a filter that applies to notebook documents
 * - `cells`: a filter that applies to the cells of matching notebook documents.
 *
 * @since 3.17.0
 */
type NotebookSelector struct {
	/**
	 * The notebook to be synced If a string
	 * value is provided it matches against the
	 * notebook type. '*' matches every notebook.
	 */
	Notebook any `json:"notebook,omitempty"` // string | NotebookDocumentFilter

	/**
	 * The cells of the matching notebook to be synced.
	 */
	Cells []NotebookCellSelector `json:"cells,omitempty"`
}

/**
 * A notebook document filter represents a filter on notebook documents.
 *
 * @since 3.17.0
 */
type NotebookDocumentFilter struct {
	/**
	 * The type of the enclosing notebook.
	 */
	NotebookType *string `json:"notebookType,omitempty"`

	/**
	 * A Uri [scheme](#Uri.scheme), like `file` or `untitled`.
	 */
	Scheme *string `json:"scheme,omitempty"`

	/**
	 * A glob pattern.
	 */
	Pattern *string `json:"pattern,omitempty"`
}

/**
 * A notebook cell selector.
 *
 * @since 3.17.0
 */
type NotebookCellSelector struct {
	Language string `json:"language"`
}
