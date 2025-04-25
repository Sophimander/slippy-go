package slippygo

import (
	"fmt"
	"testing"
)

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
