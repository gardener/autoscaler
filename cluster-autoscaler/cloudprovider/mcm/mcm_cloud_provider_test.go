package mcm

import "testing"

var (
	stop = make(chan struct{})
)

func TestDeleteNodes(t *testing.T) {
	createMcmManager(stop, testNamespace)
}
