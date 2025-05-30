package cmd

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete messages meeting certain conditions from a queue",
	Args:  cobra.NoArgs,

	RunE: runDeleteCmd,
}

type messageDeleter struct {
	messageDecoder

	client *azqueue.QueueClient
	whatIf bool
}

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	script := cmd.Flag(FlagScript).Value.String()
	tmpl, err := getScriptTemplate(script)
	if nil != err {
		fmt.Fprintf(cmd.ErrOrStderr(), "Failed to parse script template: %v\n", err)
	}

	queueName := cmd.Flag(FlagQueue).Value.String()
	client, err := getQueueClient(queueName)
	if nil != err {
		return nil
	}

	deleter := newMessageDeleter(cmd, client)
	numMatched := 0

	for {
		msgs, cont, err := deleter.getMessages()
		if nil != err {
			// TODO: log error
			return err
		}

		if len(msgs) <= 0 {
			break
		}

		for _, msg := range msgs {
			match, err := execScriptTemplate(cmd, tmpl, msg)
			if nil != err {
				fmt.Fprintf(cmd.ErrOrStderr(), "Script template execution failed: %v\n", err)
				// TOOD: flag for "continue-on-error"
				return err
			}

			// TODO: flag for invert (delete all message that do *NOT* match)
			if len(match) > 0 {
				numMatched += 1
				if !deleter.whatIf {
					_, err := client.DeleteMessage(cmd.Context(), msg.MessageID, *msg.popReceipt, nil)
					if nil != err {
						fmt.Fprintf(cmd.ErrOrStderr(), "Failed to delete message %q: %v\n", msg.MessageID, err)
					}
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Message %q matched: %s.\n", msg.MessageID, match)
				}
			}
		}

		if !cont {
			break
		}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Matched %d messages from queue %q.\n", numMatched, queueName)

	return nil
}

func newMessageDeleter(cmd *cobra.Command, client *azqueue.QueueClient) messageDeleter {
	d := messageDeleter{
		messageDecoder: newMessageDecoder(cmd),
		client:         client,
	}

	if w := getBoolFlagValue(cmd, FlagWhatIf); nil != w {
		d.whatIf = *w
	}

	return d
}

func (d messageDeleter) getMessages() ([]*queueMessage, bool, error) {
	if d.whatIf {
		opts := &azqueue.PeekMessagesOptions{
			NumberOfMessages: to.Ptr(int32(32)),
		}

		if resp, err := d.client.PeekMessages(d.cmd.Context(), opts); nil != err {
			return nil, false, err
		} else {
			msgs := make([]*queueMessage, len(resp.Messages))
			for i, msg := range resp.Messages {
				msgs[i] = d.queueMessageForPeekedMessage(msg)
			}

			return msgs, false, nil
		}
	} else {
		opts := &azqueue.DequeueMessagesOptions{
			NumberOfMessages: to.Ptr(int32(32)),

			// TODO: make configurable - we don't want messages that were already viewed to show up again right away
			// VisibilityTimeout: to.Ptr(),
		}
		if resp, err := d.client.DequeueMessages(d.cmd.Context(), opts); nil != err {
			return nil, false, err
		} else {
			msgs := make([]*queueMessage, len(resp.Messages))
			for i, msg := range resp.Messages {
				msgs[i] = d.queueMessageForDequeuedMessage(msg)
			}

			return msgs, true, nil
		}
	}
}

func getScriptTemplate(script string) (*template.Template, error) {
	tmplFile := "script.template"
	tmpl := template.New(tmplFile)
	tmpl.Funcs(template.FuncMap{
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"iso": func(t time.Time) string {
			return t.UTC().Format(time.RFC3339)
		},
		"int": func(f float64) int {
			return int(f)
		},
	})
	return tmpl.Parse(script)
}

func execScriptTemplate(cmd *cobra.Command, tmpl *template.Template, msg *queueMessage) ([]byte, error) {
	buf := &bytes.Buffer{}
	// TODO: consider optional tee-ing of template output
	// w := io.MultiWriter(buf, cmd.ErrOrStderr())
	w := buf
	if err := tmpl.Execute(w, msg); nil != err {
		return nil, err
	}

	// TODO: when tee-ing:
	// fmt.Fprintln(cmd.ErrOrStderr())
	return buf.Bytes(), nil
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringP(FlagQueue, "q", "", "name of the queue")
	deleteCmd.MarkFlagRequired(FlagQueue)

	deleteCmd.Flags().StringP(FlagScript, "s", "", "message evaluation script; supports go templates (https://pkg.go.dev/text/template)")
	deleteCmd.MarkFlagRequired(FlagScript)

	deleteCmd.Flags().BoolP(FlagDecodeBase64, "b", false, "base64-decode the message text from base64")
	deleteCmd.Flags().BoolP(FlagDecodeJson, "j", false, "JSON-decode message text (after optional base64 decoding)")
	deleteCmd.Flags().Bool(FlagWhatIf, false, "only peek at the first 32 messages and evaluate; do NOT dequeue or delete")
}
