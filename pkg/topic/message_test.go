package topic

import (
	"reflect"
	"testing"
)

func TestMessagesAreRoundTrippedCorrectly(t *testing.T) {

	message := NewMessage([]byte("hello"))

	if !(message.String() == "hello") {
		t.Error("Messages aren't roundtripping content correctly.")
	}

	content := message.Bytes()

	if !reflect.DeepEqual(content, []byte("hello")) {
		t.Error("Messages aren't roundtripping content correctly.")
	}

}
