package commands

import (
	"os"
	"testing"

	"github.com/mitchellh/cli"
)

func TestSignupCommand(t *testing.T) {
	singup := SignupCommand{
		UI: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	_ = singup
}
