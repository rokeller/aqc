package cmd

import (
	"testing"
)

func TestExecute(t *testing.T) {
	type args struct {
		errHandlerFactory func(*testing.T) func(error)
	}
	tests := []struct {
		name    string
		cmdArgs []string
		args    args
	}{
		{
			name:    "Missing command",
			cmdArgs: []string{},
			args: args{
				errHandlerFactory: func(t *testing.T) func(error) {
					return func(err error) {
						t.Error("Expected no error")
					}
				},
			},
		},
		// {
		// 	name:    "Add with missing required flags",
		// 	cmdArgs: []string{"add"},
		// 	args: args{
		// 		errHandlerFactory: func(t *testing.T) func(error) {
		// 			return func(got error) {
		// 				want := errors.New("requires at least 1 arg(s), only received 0")
		// 				if got.Error() != want.Error() {
		// 					t.Errorf("got %q, want %q", got, want)
		// 				}
		// 			}
		// 		},
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.cmdArgs)
			Execute(tt.args.errHandlerFactory(t))
		})
	}
}
