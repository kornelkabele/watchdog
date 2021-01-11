package utils

import (
	"fmt"
	"testing"
)

func TestWaitNetworkAvailable(t *testing.T) {
	if !WaitNetworkAvailable() {
		t.Error() // to indicate test failed
	}
}

func TestSplitCommand(t *testing.T) {
	result := SplitCommand("ffmpeg \"-rtsp_transport\" tcp -i \"rtsp:xxx yyy\" -frames:v 1")
	expected := []string{"ffmpeg", "-rtsp_transport", "tcp", "-i", "rtsp:xxx yyy", "-frames:v", "1"}
	if !strSliceEq(result, expected) {
		t.Error() // to indicate test failed
	}
}

func strSliceEq(a, b []string) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			fmt.Println(a[i], b[i])
			return false
		}
	}

	return true
}
