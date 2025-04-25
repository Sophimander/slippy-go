package slippygo

import (
	"fmt"
	"testing"
)

var TimeData = []string{
	"MORS#762",
	"XATU#0",
	"RON#404",
	"ANZ#139",
	"SHOOP#0",
	"POLY#832",
	"GOOPY#1",
	"BEL#306",
	"JEO#807",
	"SO#0",
	"NAT#4713",
}

func TestMain(t *testing.T) {
	fmt.Println("Main Test")

	client := NewClient()
	for _, cc := range TimeData {
		sUser, err := client.Run(cc)
		if err != nil {
			t.Error(err)
			return
		}

		fmt.Println(sUser)
	}
}
