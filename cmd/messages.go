package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/spf13/cobra"
)

type messageDecoder struct {
	cmd    *cobra.Command
	base64 bool
	json   bool
}

type queueMessage struct {
	// The number of times the message has been dequeued.
	DequeueCount int64
	// The time that the Message will expire and be automatically deleted.
	ExpirationTime time.Time
	// The time the Message was inserted into the Queue.
	InsertionTime time.Time
	// The Id of the Message.
	MessageID string
	// This value is required to delete the Message. If deletion fails using
	// this popreceipt then the message has been dequeued by another client.
	popReceipt *string
	// The content of the Message.
	MessageText string

	// The JSON decoded message
	MessageJson any
}

func newMessageDecoder(cmd *cobra.Command) messageDecoder {
	d := messageDecoder{
		cmd: cmd,
	}

	if b64 := getBoolFlagValue(cmd, FlagDecodeBase64); nil != b64 {
		d.base64 = *b64
	}

	if json := getBoolFlagValue(cmd, FlagDecodeJson); nil != json {
		d.json = *json
	}

	return d
}
func (d messageDecoder) queueMessageForDequeuedMessage(msg *azqueue.DequeuedMessage) *queueMessage {
	res := &queueMessage{
		DequeueCount:   *msg.DequeueCount,
		ExpirationTime: *msg.ExpirationTime,
		InsertionTime:  *msg.InsertionTime,
		MessageID:      *msg.MessageID,
		popReceipt:     msg.PopReceipt,
		MessageText:    *msg.MessageText,
	}
	return d.decode(res)
}

func (d messageDecoder) queueMessageForPeekedMessage(msg *azqueue.PeekedMessage) *queueMessage {
	res := &queueMessage{
		DequeueCount:   *msg.DequeueCount,
		ExpirationTime: *msg.ExpirationTime,
		InsertionTime:  *msg.InsertionTime,
		MessageID:      *msg.MessageID,
		MessageText:    *msg.MessageText,
	}
	return d.decode(res)
}

func (d messageDecoder) decode(msg *queueMessage) *queueMessage {
	text := msg.MessageText
	if "" != text {
		if d.base64 {
			if decoded, err := base64.StdEncoding.DecodeString(text); nil != err {
				fmt.Fprintf(d.cmd.ErrOrStderr(), "Failed to base64-decode message: %v\n", err)
			} else {
				text = string(decoded)
			}
		}

		if d.json {
			var val jsonValue
			if err := json.Unmarshal([]byte(text), &val); nil != err {
				fmt.Fprintf(d.cmd.ErrOrStderr(), "Failed to JSON-decode message %q: %v\n", text, err)
			} else {
				msg.MessageJson = val.Get()
			}
		}
	}
	msg.MessageText = text

	return msg
}
