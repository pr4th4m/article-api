package endpoint

import (
	"testing"
)

func TestVOne(t *testing.T) {

	result := VOne("test")
	if result != "/api/v1/test" {
		t.Errorf("Endpoint not prefixed")
	}

}
