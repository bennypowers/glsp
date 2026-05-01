package protocol

import (
	"encoding/json"
	"testing"
)

func TestBoolOrStringUnmarshalJSON(t *testing.T) {
	t.Run("bool true", func(t *testing.T) {
		var bos BoolOrString
		if err := json.Unmarshal([]byte(`true`), &bos); err != nil {
			t.Fatal(err)
		}
		if bos.Value != true {
			t.Errorf("Value = %v, want true", bos.Value)
		}
	})

	t.Run("bool false", func(t *testing.T) {
		var bos BoolOrString
		if err := json.Unmarshal([]byte(`false`), &bos); err != nil {
			t.Fatal(err)
		}
		if bos.Value != false {
			t.Errorf("Value = %v, want false", bos.Value)
		}
	})

	t.Run("string", func(t *testing.T) {
		var bos BoolOrString
		if err := json.Unmarshal([]byte(`"hello"`), &bos); err != nil {
			t.Fatal(err)
		}
		if bos.Value != "hello" {
			t.Errorf("Value = %v, want hello", bos.Value)
		}
	})

	t.Run("in struct field", func(t *testing.T) {
		type Container struct {
			BOS BoolOrString `json:"bos"`
		}
		var c Container
		if err := json.Unmarshal([]byte(`{"bos":true}`), &c); err != nil {
			t.Fatal(err)
		}
		if c.BOS.Value != true {
			t.Errorf("BOS.Value = %v, want true", c.BOS.Value)
		}
	})
}

func TestBoolOrStringStringNilSafe(t *testing.T) {
	var bos BoolOrString
	got := bos.String()
	if got != "" {
		t.Errorf("String() = %q, want empty", got)
	}
}

func TestTextDocumentEditUnmarshalAnnotatedTextEdit(t *testing.T) {
	input := `{
		"textDocument": {"uri": "file:///a.go", "version": 1},
		"edits": [
			{"range": {"start": {"line":0,"character":0}, "end": {"line":0,"character":5}}, "newText": "plain"},
			{"range": {"start": {"line":1,"character":0}, "end": {"line":1,"character":5}}, "newText": "annotated", "annotationId": "ann1"}
		]
	}`
	var tde TextDocumentEdit
	if err := json.Unmarshal([]byte(input), &tde); err != nil {
		t.Fatal(err)
	}
	if len(tde.Edits) != 2 {
		t.Fatalf("len(Edits) = %d, want 2", len(tde.Edits))
	}
	if _, ok := tde.Edits[0].(TextEdit); !ok {
		t.Errorf("Edits[0] type = %T, want TextEdit", tde.Edits[0])
	}
	ate, ok := tde.Edits[1].(AnnotatedTextEdit)
	if !ok {
		t.Fatalf("Edits[1] type = %T, want AnnotatedTextEdit", tde.Edits[1])
	}
	if ate.AnnotationID != "ann1" {
		t.Errorf("AnnotationID = %q, want ann1", ate.AnnotationID)
	}
}

func TestWorkspaceEditUnmarshalDocumentChanges(t *testing.T) {
	input := `{
		"documentChanges": [
			{"textDocument": {"uri": "file:///a.go", "version": 1}, "edits": []},
			{"kind": "create", "uri": "file:///new.go"},
			{"kind": "rename", "oldUri": "file:///old.go", "newUri": "file:///new.go"},
			{"kind": "delete", "uri": "file:///del.go"}
		]
	}`
	var we WorkspaceEdit
	if err := json.Unmarshal([]byte(input), &we); err != nil {
		t.Fatal(err)
	}
	if len(we.DocumentChanges) != 4 {
		t.Fatalf("len(DocumentChanges) = %d, want 4", len(we.DocumentChanges))
	}
	if _, ok := we.DocumentChanges[0].(TextDocumentEdit); !ok {
		t.Errorf("[0] type = %T, want TextDocumentEdit", we.DocumentChanges[0])
	}
	if _, ok := we.DocumentChanges[1].(CreateFile); !ok {
		t.Errorf("[1] type = %T, want CreateFile", we.DocumentChanges[1])
	}
	if _, ok := we.DocumentChanges[2].(RenameFile); !ok {
		t.Errorf("[2] type = %T, want RenameFile", we.DocumentChanges[2])
	}
	if _, ok := we.DocumentChanges[3].(DeleteFile); !ok {
		t.Errorf("[3] type = %T, want DeleteFile", we.DocumentChanges[3])
	}
}

func TestCompletionItemUnmarshalInsertReplaceEdit(t *testing.T) {
	input := `{
		"label": "test",
		"textEdit": {
			"newText": "test",
			"insert": {"start": {"line":0,"character":0}, "end": {"line":0,"character":3}},
			"replace": {"start": {"line":0,"character":0}, "end": {"line":0,"character":5}}
		}
	}`
	var ci CompletionItem
	if err := json.Unmarshal([]byte(input), &ci); err != nil {
		t.Fatal(err)
	}
	ire, ok := ci.TextEdit.(InsertReplaceEdit)
	if !ok {
		t.Fatalf("TextEdit type = %T, want InsertReplaceEdit", ci.TextEdit)
	}
	if ire.NewText != "test" {
		t.Errorf("NewText = %q, want test", ire.NewText)
	}
}

func TestCompletionItemUnmarshalRegularTextEdit(t *testing.T) {
	input := `{
		"label": "test",
		"textEdit": {
			"newText": "test",
			"range": {"start": {"line":0,"character":0}, "end": {"line":0,"character":3}}
		}
	}`
	var ci CompletionItem
	if err := json.Unmarshal([]byte(input), &ci); err != nil {
		t.Fatal(err)
	}
	te, ok := ci.TextEdit.(TextEdit)
	if !ok {
		t.Fatalf("TextEdit type = %T, want TextEdit", ci.TextEdit)
	}
	if te.NewText != "test" {
		t.Errorf("NewText = %q, want test", te.NewText)
	}
}

func TestMarkedStringUnmarshalJSON(t *testing.T) {
	t.Run("plain string roundtrip", func(t *testing.T) {
		var ms MarkedString
		if err := json.Unmarshal([]byte(`"hello"`), &ms); err != nil {
			t.Fatal(err)
		}
		data, err := json.Marshal(ms)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != `"hello"` {
			t.Errorf("roundtrip = %s, want \"hello\"", data)
		}
	})

	t.Run("language block roundtrip", func(t *testing.T) {
		var ms MarkedString
		if err := json.Unmarshal([]byte(`{"language":"go","value":"func main()"}`), &ms); err != nil {
			t.Fatal(err)
		}
		data, err := json.Marshal(ms)
		if err != nil {
			t.Fatal(err)
		}
		var mss MarkedStringStruct
		if err := json.Unmarshal(data, &mss); err != nil {
			t.Fatal(err)
		}
		if mss.Language != "go" {
			t.Errorf("Language = %q, want go", mss.Language)
		}
	})

	t.Run("in struct field", func(t *testing.T) {
		type Container struct {
			MS MarkedString `json:"ms"`
		}
		var c Container
		if err := json.Unmarshal([]byte(`{"ms":"hello"}`), &c); err != nil {
			t.Fatal(err)
		}
		data, err := json.Marshal(c.MS)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != `"hello"` {
			t.Errorf("roundtrip = %s, want \"hello\"", data)
		}
	})
}

func TestSemanticTokensOptionsUnmarshalWorkDoneProgress(t *testing.T) {
	input := `{
		"workDoneProgress": true,
		"legend": {"tokenTypes": ["type"], "tokenModifiers": ["mod"]},
		"full": true
	}`
	var opts SemanticTokensOptions
	if err := json.Unmarshal([]byte(input), &opts); err != nil {
		t.Fatal(err)
	}
	if opts.WorkDoneProgress == nil || !*opts.WorkDoneProgress {
		t.Error("WorkDoneProgress should be true")
	}
}

func TestTextDocumentSyncOptionsUnmarshalNullSave(t *testing.T) {
	input := `{"save": null}`
	var opts TextDocumentSyncOptions
	if err := json.Unmarshal([]byte(input), &opts); err != nil {
		t.Fatal(err)
	}
	if opts.Save != nil {
		t.Errorf("Save = %v, want nil", opts.Save)
	}
}

func TestTextDocumentSyncOptionsUnmarshalBoolSave(t *testing.T) {
	input := `{"save": true}`
	var opts TextDocumentSyncOptions
	if err := json.Unmarshal([]byte(input), &opts); err != nil {
		t.Fatal(err)
	}
	if opts.Save != true {
		t.Errorf("Save = %v, want true", opts.Save)
	}
}

func TestTextDocumentSyncOptionsUnmarshalOptionsSave(t *testing.T) {
	input := `{"save": {"includeText": true}}`
	var opts TextDocumentSyncOptions
	if err := json.Unmarshal([]byte(input), &opts); err != nil {
		t.Fatal(err)
	}
	so, ok := opts.Save.(SaveOptions)
	if !ok {
		t.Fatalf("Save type = %T, want SaveOptions", opts.Save)
	}
	if so.IncludeText == nil || !*so.IncludeText {
		t.Error("IncludeText should be true")
	}
}
