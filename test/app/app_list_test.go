package app_test

import (
	"github.com/Rollmops/pctl/app"
	"github.com/kami-zh/go-capturer"
	"testing"
)

func TestAppListCommand(t *testing.T) {
	pctlApp := app.CreateCliApp()

	out := capturer.CaptureOutput(func() {
		err := pctlApp.Run([]string{"pctl", "list"})
		if err != nil {
			t.Fatal(err)
		}
	})

	expectedOutput := ""

	if out != expectedOutput {
		t.Fatalf("%s != %s", out, expectedOutput)
	}

}
