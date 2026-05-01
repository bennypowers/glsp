package protocol

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bennypowers/glsp"
)

func TestGetSetTraceValue(t *testing.T) {
	for _, tv := range []TraceValue{TraceValueOff, TraceValueMessage, TraceValueVerbose} {
		SetTraceValue(tv)
		if got := GetTraceValue(); got != tv {
			t.Errorf("SetTraceValue(%s): GetTraceValue() = %s", tv, got)
		}
	}
}

func TestSetTraceValueMessages(t *testing.T) {
	SetTraceValue("messages")
	if got := GetTraceValue(); got != TraceValueMessage {
		t.Errorf("SetTraceValue(\"messages\"): GetTraceValue() = %s, want %s", got, TraceValueMessage)
	}
}

func TestHasTraceLevel(t *testing.T) {
	tests := []struct {
		current TraceValue
		query   TraceValue
		want    bool
	}{
		{TraceValueOff, TraceValueOff, false},
		{TraceValueOff, TraceValueMessage, false},
		{TraceValueOff, TraceValueVerbose, false},
		{TraceValueMessage, TraceValueOff, false},
		{TraceValueMessage, TraceValueMessage, true},
		{TraceValueMessage, TraceValueVerbose, false},
		{TraceValueVerbose, TraceValueOff, true},
		{TraceValueVerbose, TraceValueMessage, true},
		{TraceValueVerbose, TraceValueVerbose, true},
	}
	for _, tt := range tests {
		SetTraceValue(tt.current)
		got := HasTraceLevel(tt.query)
		if got != tt.want {
			t.Errorf("HasTraceLevel(%s) with current=%s: got %v, want %v", tt.query, tt.current, got, tt.want)
		}
	}
}

func TestHasTraceMessageType(t *testing.T) {
	tests := []struct {
		current TraceValue
		msgType MessageType
		want    bool
	}{
		{TraceValueOff, MessageTypeError, false},
		{TraceValueOff, MessageTypeWarning, false},
		{TraceValueOff, MessageTypeInfo, false},
		{TraceValueOff, MessageTypeLog, false},
		{TraceValueMessage, MessageTypeError, true},
		{TraceValueMessage, MessageTypeWarning, true},
		{TraceValueMessage, MessageTypeInfo, true},
		{TraceValueMessage, MessageTypeLog, false},
		{TraceValueVerbose, MessageTypeError, true},
		{TraceValueVerbose, MessageTypeWarning, true},
		{TraceValueVerbose, MessageTypeInfo, true},
		{TraceValueVerbose, MessageTypeLog, true},
	}
	for _, tt := range tests {
		SetTraceValue(tt.current)
		got := HasTraceMessageType(tt.msgType)
		if got != tt.want {
			t.Errorf("HasTraceMessageType(%d) with current=%s: got %v, want %v", tt.msgType, tt.current, got, tt.want)
		}
	}
}

func TestHasTraceLevelPanicsOnInvalidLevel(t *testing.T) {
	SetTraceValue("bogus")
	defer func() {
		if r := recover(); r == nil {
			t.Error("HasTraceLevel did not panic on invalid trace level")
		}
	}()
	HasTraceLevel(TraceValueMessage)
}

func TestHasTraceMessageTypePanicsOnInvalidType(t *testing.T) {
	SetTraceValue(TraceValueVerbose)
	defer func() {
		if r := recover(); r == nil {
			t.Error("HasTraceMessageType did not panic on invalid message type")
		}
	}()
	HasTraceMessageType(MessageType(99))
}

func TestTrace(t *testing.T) {
	t.Run("fires when level matches", func(t *testing.T) {
		SetTraceValue(TraceValueVerbose)
		var called atomic.Bool
		var mu sync.Mutex
		var gotMethod string
		ctx := &glsp.Context{
			Notify: func(method string, params any) {
				mu.Lock()
				defer mu.Unlock()
				called.Store(true)
				gotMethod = method
			},
		}
		err := Trace(ctx, MessageTypeError, "test message")
		if err != nil {
			t.Fatalf("Trace returned error: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
		if !called.Load() {
			t.Error("Trace did not fire notify")
		}
		mu.Lock()
		if gotMethod != ServerWindowLogMessage {
			t.Errorf("Trace notified method %s, want %s", gotMethod, ServerWindowLogMessage)
		}
		mu.Unlock()
	})

	t.Run("does not fire when level does not match", func(t *testing.T) {
		SetTraceValue(TraceValueOff)
		var called atomic.Bool
		ctx := &glsp.Context{
			Notify: func(method string, params any) {
				called.Store(true)
			},
		}
		err := Trace(ctx, MessageTypeError, "should not fire")
		if err != nil {
			t.Fatalf("Trace returned error: %v", err)
		}
		time.Sleep(50 * time.Millisecond)
		if called.Load() {
			t.Error("Trace fired notify when it should not have")
		}
	})
}
