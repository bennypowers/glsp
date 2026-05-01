package protocol

// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#positionEncodingKind

type PositionEncodingKind string

const (
	PositionEncodingKindUTF8  = PositionEncodingKind("utf-8")
	PositionEncodingKindUTF16 = PositionEncodingKind("utf-16")
	PositionEncodingKindUTF32 = PositionEncodingKind("utf-32")
)

type MarkdownClientCapabilities struct {
	Parser  string  `json:"parser"`
	Version *string `json:"version,omitempty"`

	/**
	 * A list of HTML tags that the client allows / supports in
	 * Markdown.
	 *
	 * @since 3.17.0
	 */
	AllowedTags []string `json:"allowedTags,omitempty"`
}
