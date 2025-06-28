package main

import (
	"testing"

	"github.com/rokeller/aqc/test_helpers"
)

func TestMain(t *testing.T) {
	test_helpers.ExecuteWithExit(t, "TestMain", func(t *testing.T) {
		main()
	}, 0)
}
