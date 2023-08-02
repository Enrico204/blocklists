package utils

import (
	"go.uber.org/zap"
	"strings"
	"testing"
)

func TestExtractNetworks(t *testing.T) {
	dirty := `# comment in one line
1.2.3.4
5.6.7.8 # comment in another line
# comment in one line
9.10.11.12/32

13.14.15.16
::2
`
	clean := `1.2.3.4/32
5.6.7.8/32
9.10.11.12/32
13.14.15.16/32
::2/128
`

	networks, err := ExtractNetworks(zap.NewNop().Sugar(), strings.NewReader(dirty))
	if err != nil {
		t.Fatal(err)
	} else if len(networks) != 5 {
		t.Fatal("len network wrong ", len(networks))
	}

	var buf strings.Builder
	for _, addr := range networks {
		buf.WriteString(addr.String())
		buf.WriteRune('\n')
	}
	if buf.String() != clean {
		t.Fatal("different outputs")
	}
}
