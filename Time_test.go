package connect

import (
	"encoding/json"
	"testing"
)

func TestTimeUnmarshalJSON(t *testing.T) {
	var t0 Time

	input := []byte(`"2019-01-12T11:45:23.0"`)

	err := json.Unmarshal(input, &t0)
	if err != nil {
		t.Fatalf("Error parsing %s: %s", string(input), err.Error())
	}

	if t0.String() != "2019-01-12 11:45:23 +0000 UTC" {
		t.Errorf("Failed to parse `%s` correct, got %s", string(input), t0.String())
	}
}
