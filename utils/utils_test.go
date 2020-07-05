package utils

import (
	"fmt"
	"testing"
)


func TestUtils(t *testing.T) {
	if Md5sum("webhooks") == "C10F40999B74C408263F790B30E70EFE" {
		fmt.Println("Md5sum method is mistake.")
	}
}
