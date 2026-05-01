package protocol

import (
	"github.com/bennypowers/glsp"
	protocol316 "github.com/bennypowers/glsp/protocol_3_16"
)

/**
 * Client capabilities specific to diagnostic pull requests.
 *
 * @since 3.17.0
 */
type DiagnosticClientCapabilities struct {
	/**
	 * Whether implementation supports dynamic registration. If this is set to
	 * `true` the client supports the new
	 * `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
	 * return value for the corresponding server capability as well.
	 */
	DynamicRegistration bool `json:"dynamicRegistration"`

	/**
	 * Whether the clients supports related documents for document diagnostic
	 * pulls.
	 */
	RelatedDocumentSupport bool `json:"relatedDocumentSupport"`
}

/**
 * Diagnostic options.
 *
 * @since 3.17.0
 */
type DiagnosticOptions struct {
	protocol316.WorkDoneProgressOptions
	/**
	 * An optional identifier under which the diagnostics are
	 * managed by the client.
	 */
	Identifier *string `json:"identifier"`

	/**
	 * Whether the language has inter file dependencies meaning that
	 * editing code in one file can result in a different diagnostic
	 * set in another file. Inter file dependencies are common for
	 * most programming languages and typically uncommon for linters.
	 */
	InterFileDependencies bool `json:"interFileDependencies"`

	/**
	 * The server provides support for workspace diagnostics as well.
	 */
	WorkspaceDiagnostics bool `json:"workspaceDiagnostics"`
}

/**
 * Diagnostic registration options.
 *
 * @since 3.17.0
 */
type DiagnosticRegistrationOptions struct {
	protocol316.TextDocumentRegistrationOptions
	DiagnosticOptions
	protocol316.StaticRegistrationOptions
}

const MethodTextDocumentDiagnostic = protocol316.Method("textDocument/diagnostic")

type TextDocumentDiagnosticFunc func(context *glsp.Context, params *DocumentDiagnosticParams) (any, error)

/**
 * Parameters of the document diagnostic request.
 *
 * @since 3.17.0
 */
type DocumentDiagnosticParams struct {
	protocol316.WorkDoneProgressParams
	protocol316.PartialResultParams

	/**
	 * The text document.
	 */
	TextDocument protocol316.TextDocumentIdentifier `json:"textDocument"`

	/**
	 * The additional identifier  provided during registration.
	 */
	Identifier *string `json:"identifier,omitempty"`

	/**
	 * The result id of a previous response if provided.
	 */
	PreviousResultId *string `json:"previousResultId,omitempty"`
}

/**
 * The result of a document diagnostic pull request. A report can
 * either be a full report containing all diagnostics for the
 * requested document or a unchanged report indicating that nothing
 * has changed in terms of diagnostics in comparison to the last
 * pull request.
 *
 * @since 3.17.0
 */
type DocumentDiagnosticReport any // RelatedFullDocumentDiagnosticReport | RelatedUnchangedDocumentDiagnosticReport

/**
 * The document diagnostic report kinds.
 *
 * @since 3.17.0
 */
type DocumentDiagnosticReportKind string

const (
	/**
	 * A diagnostic report with a full
	 * set of problems.
	 */
	DocumentDiagnosticReportKindFull = DocumentDiagnosticReportKind("full")
	/**
	 * A report indicating that the last
	 * returned report is still accurate.
	 */
	DocumentDiagnosticReportKindUnchanged = DocumentDiagnosticReportKind("unchanged")
)

/**
 * A diagnostic report with a full set of problems.
 *
 * @since 3.17.0
 */
type FullDocumentDiagnosticReport struct {
	/**
	 * A full document diagnostic report.
	 */
	Kind string `json:"kind"`

	/**
	 * An optional result id. If provided it will
	 * be sent on the next diagnostic request for the
	 * same document.
	 */
	ResultID *string `json:"resultId,omitempty"`

	/**
	 * The actual items.
	 */
	Items []protocol316.Diagnostic `json:"items"`
}

/**
 * A diagnostic report indicating that the last returned
 * report is still accurate.
 *
 * @since 3.17.0
 */
type UnchangedDocumentDiagnosticReport struct {
	/**
	 * A document diagnostic report indicating
	 * no changes to the last result. A server can
	 * only return `unchanged` if result ids are
	 * provided.
	 */
	Kind string `json:"kind"`

	/**
	 * A result id which will be sent on the next
	 * diagnostic request for the same document.
	 */
	ResultID string `json:"resultId"`
}

/**
 * A full diagnostic report with a set of related documents.
 *
 * @since 3.17.0
 */
type RelatedFullDocumentDiagnosticReport struct {
	FullDocumentDiagnosticReport
	/**
	 * Diagnostics of related documents. This information is useful
	 * in programming languages where code in a file A can generate
	 * diagnostics in a file B which A depends on. An example of
	 * such a language is C/C++ where marco definitions in a file
	 * a.cpp and result in errors in a header file b.hpp.
	 *
	 * @since 3.17.0
	 */
	RelatedDocuments map[protocol316.DocumentUri]any `json:"relatedDocuments,omitempty"`
}

/**
 * An unchanged diagnostic report with a set of related documents.
 *
 * @since 3.17.0
 */
type RelatedUnchangedDocumentDiagnosticReport struct {
	UnchangedDocumentDiagnosticReport
	/**
	 * Diagnostics of related documents. This information is useful
	 * in programming languages where code in a file A can generate
	 * diagnostics in a file B which A depends on. An example of
	 * such a language is C/C++ where marco definitions in a file
	 * a.cpp and result in errors in a header file b.hpp.
	 *
	 * @since 3.17.0
	 */
	RelatedDocuments map[protocol316.DocumentUri]any `json:"relatedDocuments,omitempty"`
}

/**
 * A partial result for a document diagnostic report.
 *
 * @since 3.17.0
 */
type DocumentDiagnosticReportPartialResult struct {
	RelatedDocuments map[protocol316.DocumentUri]any `json:"relatedDocuments"`
}

/**
 * Cancellation data returned from a diagnostic request.
 *
 * @since 3.17.0
 */
type DiagnosticServerCancellationData struct {
	RetriggerRequest bool `json:"retriggerRequest"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#workspace_diagnostic

/**
 * Workspace client capabilities specific to diagnostic pull requests.
 *
 * @since 3.17.0
 */
type DiagnosticWorkspaceClientCapabilities struct {
	/**
	 * Whether the client implementation supports a refresh request sent from
	 * the server to the client.
	 *
	 * Note that this event is global and will force the client to re-calculate
	 * all pulled diagnostics for all open documents. It should be used with
	 * absolute care and is useful for situation where a server for example
	 * detects a project wide change that requires such a calculation.
	 */
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}

const MethodWorkspaceDiagnostic = protocol316.Method("workspace/diagnostic")

type WorkspaceDiagnosticFunc func(context *glsp.Context, params *WorkspaceDiagnosticParams) (*WorkspaceDiagnosticReport, error)

/**
 * Parameters of the workspace diagnostic request.
 *
 * @since 3.17.0
 */
type WorkspaceDiagnosticParams struct {
	protocol316.WorkDoneProgressParams
	protocol316.PartialResultParams

	/**
	 * The additional identifier provided during registration.
	 */
	Identifier *string `json:"identifier,omitempty"`

	/**
	 * The currently known diagnostic reports with their
	 * previous result ids.
	 */
	PreviousResultIds []PreviousResultId `json:"previousResultIds"`
}

/**
 * A previous result id in a workspace pull request.
 *
 * @since 3.17.0
 */
type PreviousResultId struct {
	/**
	 * The URI for which the client knows a result id.
	 */
	URI protocol316.DocumentUri `json:"uri"`

	/**
	 * The value of the previous result id.
	 */
	Value string `json:"value"`
}

/**
 * A workspace diagnostic report.
 *
 * @since 3.17.0
 */
type WorkspaceDiagnosticReport struct {
	Items []WorkspaceDocumentDiagnosticReport `json:"items"`
}

/**
 * A workspace diagnostic document report.
 *
 * @since 3.17.0
 */
type WorkspaceDocumentDiagnosticReport any // WorkspaceFullDocumentDiagnosticReport | WorkspaceUnchangedDocumentDiagnosticReport

/**
 * A full document diagnostic report for a workspace diagnostic result.
 *
 * @since 3.17.0
 */
type WorkspaceFullDocumentDiagnosticReport struct {
	FullDocumentDiagnosticReport

	/**
	 * The URI for which diagnostic information is reported.
	 */
	URI protocol316.DocumentUri `json:"uri"`

	/**
	 * The version number for which the diagnostics are reported.
	 * If the document is not marked as open `null` can be provided.
	 */
	Version *protocol316.Integer `json:"version"`
}

/**
 * An unchanged document diagnostic report for a workspace diagnostic result.
 *
 * @since 3.17.0
 */
type WorkspaceUnchangedDocumentDiagnosticReport struct {
	UnchangedDocumentDiagnosticReport

	/**
	 * The URI for which diagnostic information is reported.
	 */
	URI protocol316.DocumentUri `json:"uri"`

	/**
	 * The version number for which the diagnostics are reported.
	 * If the document is not marked as open `null` can be provided.
	 */
	Version *protocol316.Integer `json:"version"`
}

/**
 * A partial result for a workspace diagnostic report.
 *
 * @since 3.17.0
 */
type WorkspaceDiagnosticReportPartialResult struct {
	Items []WorkspaceDocumentDiagnosticReport `json:"items"`
}

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#diagnostic_refresh

const MethodWorkspaceDiagnosticRefresh = protocol316.Method("workspace/diagnostics/refresh")
