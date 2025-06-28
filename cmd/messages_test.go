package cmd

import (
	"reflect"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/spf13/cobra"
)

func Test_messageDecoder_queueMessageForPeekedMessage(t *testing.T) {
	timestamp1 := time.Now()
	timestamp2 := time.Now()

	type fields struct {
		base64 bool
		json   bool
	}
	type args struct {
		msg *azqueue.PeekedMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *queueMessage
	}{
		{
			name: "Simple base64 encoded text",
			fields: fields{
				base64: true,
			},
			args: args{
				msg: &azqueue.PeekedMessage{
					DequeueCount:   to.Ptr(int64(123)),
					ExpirationTime: to.Ptr(timestamp1),
					InsertionTime:  to.Ptr(timestamp2),
					MessageID:      to.Ptr("msg-id"),
					MessageText:    to.Ptr("SGVsbG8sIFdvcmxkIQ=="),
				},
			},
			want: &queueMessage{
				DequeueCount:   123,
				ExpirationTime: timestamp1,
				InsertionTime:  timestamp2,
				MessageID:      "msg-id",
				MessageText:    "Hello, World!",
			},
		},
		{
			name: "JSON encoded string",
			fields: fields{
				json: true,
			},
			args: args{
				msg: &azqueue.PeekedMessage{
					DequeueCount:   to.Ptr(int64(234)),
					ExpirationTime: to.Ptr(timestamp1),
					InsertionTime:  to.Ptr(timestamp2),
					MessageID:      to.Ptr("second-msg-id"),
					MessageText:    to.Ptr(`"Hello, World!"`),
				},
			},
			want: &queueMessage{
				DequeueCount:   234,
				ExpirationTime: timestamp1,
				InsertionTime:  timestamp2,
				MessageID:      "second-msg-id",
				MessageText:    `"Hello, World!"`,
				MessageJson:    "Hello, World!",
			},
		},
		{
			name: "Non-base64-encoded text",
			fields: fields{
				base64: true,
			},
			args: args{
				msg: &azqueue.PeekedMessage{
					DequeueCount:   to.Ptr(int64(345)),
					ExpirationTime: to.Ptr(timestamp1),
					InsertionTime:  to.Ptr(timestamp2),
					MessageID:      to.Ptr("third-msg-id"),
					MessageText:    to.Ptr("Hello, World!"),
				},
			},
			want: &queueMessage{
				DequeueCount:   345,
				ExpirationTime: timestamp1,
				InsertionTime:  timestamp2,
				MessageID:      "third-msg-id",
				MessageText:    "Hello, World!",
			},
		},
		{
			name: "Non-JSON-encoded text",
			fields: fields{
				json: true,
			},
			args: args{
				msg: &azqueue.PeekedMessage{
					DequeueCount:   to.Ptr(int64(456)),
					ExpirationTime: to.Ptr(timestamp1),
					InsertionTime:  to.Ptr(timestamp2),
					MessageID:      to.Ptr("fourth-msg-id"),
					MessageText:    to.Ptr("Hello, World!"),
				},
			},
			want: &queueMessage{
				DequeueCount:   456,
				ExpirationTime: timestamp1,
				InsertionTime:  timestamp2,
				MessageID:      "fourth-msg-id",
				MessageText:    "Hello, World!",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			d := messageDecoder{
				cmd:    cmd,
				base64: tt.fields.base64,
				json:   tt.fields.json,
			}
			if got := d.queueMessageForPeekedMessage(tt.args.msg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("messageDecoder.queueMessageForPeekedMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
