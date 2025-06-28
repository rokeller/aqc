package cmd

import (
	"testing"
)

func Test_getQueueClientForQueueURLWithSASToken(t *testing.T) {
	createTestQueue(t, "test-client-creation-sas")
	addTestQueueMesssages(t, "test-client-creation-sas", []string{"one", "two", "three"})

	type args struct {
		queueURL string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Can fetch messages",
			args: args{
				queueURL: "http://127.0.0.1:10001/devstoreaccount1/test-client-creation-sas?sv=2024-08-04&spr=https%2Chttp&st=2025-06-28T15%3A15%3A11Z&se=2100-01-01T00%3A00%3A00Z&sp=rp&sig=H6vFIUxzqT2H0tb6dR6LpNgoTQW3rqOlG3nAe6tTOvw%3D",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getQueueClientForQueueURLWithSASToken(tt.args.queueURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("getQueueClientForQueueURLWithSASToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			msgs := fetchTestQueueMesssagesFromQueue(t, got)
			if len(msgs) != 3 {
				t.Errorf("got %d messages in queue, want 3", len(msgs))
			}
			if *msgs[0].MessageText != "one" {
				t.Errorf("got %q in first message, want 'one'", *msgs[0].MessageText)
			}
			if *msgs[1].MessageText != "two" {
				t.Errorf("got %q in second message, want 'two'", *msgs[1].MessageText)
			}
			if *msgs[2].MessageText != "three" {
				t.Errorf("got %q in third message, want 'three'", *msgs[2].MessageText)
			}
		})
	}
}
