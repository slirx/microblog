package api

import (
	"encoding/json"
	"testing"
)

func TestCorrectResponse(t *testing.T) {
	r := Response{RequestID: "msg-id-1"}
	data, err := json.Marshal(&r)
	if err != nil {
		t.Fatalf("got: %s, want: nil", err)
	}

	in := string(data)
	want := `{"request_id":"msg-id-1"}`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}
}

func TestEmptyResponse(t *testing.T) {
	r := Response{RequestID: ""}
	data, err := json.Marshal(&r)
	if err != nil {
		t.Fatalf("got: %v, want: nil", err)
	}

	in := string(data)
	want := `{"request_id":""}`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}
}

func TestNewResponse(t *testing.T) {
	r := NewResponse("message-1")
	data, err := json.Marshal(&r)
	if err != nil {
		t.Fatalf("got: %v, want: nil", err)
	}

	in := string(data)
	want := `{"request_id":"message-1"}`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}
}

func TestCorrectMessageResponse(t *testing.T) {
	r := MessageResponse{
		RequestID: "msg2",
		Type:      MessageTypeSuccess,
		Message:   "test",
	}
	data, err := json.Marshal(&r)
	if err != nil {
		t.Fatalf("got: %v, want: nil", err)
	}

	in := string(data)
	want := `{"request_id":"msg2","type":"success","message":"test"}`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}
}

func TestNewMessageResponse(t *testing.T) {
	r := NewMessageResponse("message-id", MessageTypeError, "test msg")
	data, err := json.Marshal(&r)
	if err != nil {
		t.Fatalf("got: %v, want: nil", err)
	}

	in := string(data)
	want := `{"request_id":"message-id","type":"error","message":"test msg"}`

	if in != want {
		t.Fatalf("got: %s, want: %s", in, want)
	}
}
