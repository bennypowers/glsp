package protocol

import (
	"github.com/tliron/glsp"
	protocol316 "github.com/tliron/glsp/protocol_3_16"
)

const MethodTextDocumentInlayHint = protocol316.Method("textDocument/inlayHint")

type TextDocumentInlayHintFunc func(context *glsp.Context, params *InlayHintParams) ([]InlayHint, error)

/**
 * Inlay hint options used during static registration.
 *
 * @since 3.17.0
 */
type InlayHintOptions struct {
	protocol316.WorkDoneProgressOptions
}

/**
 * Inlay hint options used during static or dynamic registration.
 *
 * @since 3.17.0
 */
type InlayHintRegistrationOptions struct {
	InlayHintOptions
	protocol316.TextDocumentRegistrationOptions
	protocol316.StaticRegistrationOptions
}

/**
 * A parameter literal used in inlay hint requests.
 *
 * @since 3.17.0
 */
type InlayHintParams struct {
	/**
	 * The text document.
	 */
	TextDocument protocol316.TextDocumentIdentifier `json:"textDocument"`

	/**
	 * The document range for which inlay hints should be computed.
	 */
	Range protocol316.Range `json:"range"`
}

/**
 * Inlay hint kinds.
 *
 * @since 3.17.0
 */
type InlayHintKind int

const (
	/**
	 * An inlay hint that for a type annotation.
	 */
	InlayHintKindType InlayHintKind = 1

	/**
	 * An inlay hint that is for a parameter.
	 */
	InlayHintKindParameter InlayHintKind = 2
)

/**
 * Inlay hint information.
 *
 * @since 3.17.0
 */
type InlayHint struct {
	/**
	 * The position of this hint.
	 */
	Position protocol316.Position `json:"position"`

	/**
	 * The label of this hint. A human readable string or an array of
	 * InlayHintLabelPart label parts.
	 *
	 * *Note* that neither the string nor the label part can be empty.
	 */
	Label string `json:"label"`

	/**
	 * The kind of this hint. Can be omitted in which case the client
	 * should fall back to a reasonable default.
	 */
	Kind *InlayHintKind `json:"kind,omitempty"`

	/**
	 * Optional text edits that are performed when accepting this inlay hint.
	 *
	 * *Note* that edits are expected to change the document so that the inlay
	 * hint (or its nearest variant) is now part of the document and the inlay
	 * hint itself is now obsolete.
	 */
	Tooltip any `json:"tooltip,omitempty"` // string | MarkupContent

	/**
	 * Render padding before the hint.
	 *
	 * Note: Padding should use the editor's background color, not the
	 * background color of the hint itself. That means padding can be used
	 * to visually align/separate an inlay hint.
	 */
	PaddingLeft *bool `json:"paddingLeft,omitempty"`

	/**
	 * Render padding after the hint.
	 *
	 * Note: Padding should use the editor's background color, not the
	 * background color of the hint itself. That means padding can be used
	 * to visually align/separate an inlay hint.
	 */
	PaddingRight *bool `json:"paddingRight,omitempty"`
}
