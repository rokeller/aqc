package cmd

import (
	"reflect"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/spf13/cobra"
)

func Test_getQueueClientForCommandWithPrefix(t *testing.T) {
	createTestQueue(t, "test-client-creation-01")
	addTestQueueMesssages(t, "test-client-creation-01", []string{"test"})

	type args struct {
		cmd    *cobra.Command
		prefix string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Queue URL with SAS",
			args: args{
				cmd: func() *cobra.Command {
					c := &cobra.Command{}
					addQueueConnectionFlagsWithPrefix(c, "test-", "test")
					c.Flag("test-queue-url").Value.Set("http://127.0.0.1:10001/devstoreaccount1/test-client-creation-01?sv=2024-08-04&spr=https%2Chttp&st=2025-06-28T16%3A24%3A20Z&se=2100-01-01T00%3A00%3A00Z&sp=rp&sig=ZWsk%2B%2BUJjBlV0FX43gWq80P%2F6j6XveQVwBZmwDXHEIY%3D")

					return c
				}(),
				prefix: "test-",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getQueueClientForCommandWithPrefix(tt.args.cmd, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("getQueueClientForCommandWithPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			msgs := fetchTestQueueMesssagesFromQueue(t, got)
			if len(msgs) != 1 {
				t.Errorf("got %d message(s) in queue, want 1", len(msgs))
			}
			if *msgs[0].MessageText != "test" {
				t.Errorf("got %d in message, want 'test'", msgs[0].MessageText)
			}
		})
	}
}

func Test_getInt32FlagValue(t *testing.T) {
	type args struct {
		c        *cobra.Command
		flagName string
	}
	tests := []struct {
		name string
		args args
		want *int32
	}{
		{
			name: "Default value",
			args: args{
				c: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().Int32("test", 123, "")

					return cmd
				}(),
				flagName: "test",
			},
			want: to.Ptr(int32(123)),
		},
		{
			name: "Not found",
			args: args{
				c: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().Int32("test", 123, "")

					return cmd
				}(),
				flagName: "other",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInt32FlagValue(tt.args.c, tt.args.flagName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInt32FlagValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBoolFlagValue(t *testing.T) {
	type args struct {
		c        *cobra.Command
		flagName string
	}
	tests := []struct {
		name string
		args args
		want *bool
	}{
		{
			name: "Default value",
			args: args{
				c: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().Bool("test", true, "")

					return cmd
				}(),
				flagName: "test",
			},
			want: to.Ptr(true),
		},
		{
			name: "Not found",
			args: args{
				c: func() *cobra.Command {
					cmd := &cobra.Command{}
					cmd.Flags().Bool("test", false, "")

					return cmd
				}(),
				flagName: "other",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBoolFlagValue(tt.args.c, tt.args.flagName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getBoolFlagValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
