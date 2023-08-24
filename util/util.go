package util

import (
	"fmt"
	"os"
)

func ExitOnError(err error) {
	if err != nil {
		fmt.Printf("Error encountered: %s \n", err)
		os.Exit(1)
	}
}
