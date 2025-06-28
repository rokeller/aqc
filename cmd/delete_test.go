package cmd

import (
	"encoding/json"
	"testing"
)

type user struct {
	Username string `json:"username"`
	Id       int64  `json:"id"`
}

func makeUserJson(t *testing.T, username string, id int) string {
	obj := user{
		Username: username,
		Id:       int64(id),
	}
	json, err := json.Marshal(obj)
	if nil != err {
		t.Fatalf("failed to marshal to JSON: %v", err)
	}

	return string(json)
}

func Test_runDeleteCmd(t *testing.T) {
	createTestQueue(t, "test-delete-01")
	addTestQueueMesssages(t, "test-delete-01", []string{
		makeUserJson(t, "john.doe@user.com", 1),
		makeUserJson(t, "alice@foo.com", 2),
		makeUserJson(t, "bob@foo.com", 3),
		makeUserJson(t, "eve@foo.com", 4),
	})

	createTestQueue(t, "test-delete-02")
	addTestQueueMesssages(t, "test-delete-02", []string{"message-01"})

	createTestQueue(t, "test-delete-03")
	addTestQueueMesssages(t, "test-delete-03", []string{
		makeUserJson(t, "one@user.com", 1),
		makeUserJson(t, "ten@user.com", 10),
	})

	tc := []testCase{
		{
			name: "Delete users with a*",
			args: []string{"delete", "--use-storage-emulator", "-q=test-delete-01", "--decode-json",
				"--script", `{{ $u := .MessageJson.username | lower }}{{ if and (ge $u "a") (lt $u "b") }}delete{{ end }}`},
			stdOut: "Matched 1 message(s).",
			verify: func(t *testing.T) {
				c := fetchTestQueueMessageCount(t, "test-delete-01")
				if c != 3 {
					t.Errorf("got %d messages in queue, want 3", c)
				}
			},
		},
		{
			name: "Delete messages inserted after 2025-01-01",
			args: []string{"delete", "--use-storage-emulator", "-q=test-delete-02",
				"--script", `{{ if gt (iso .InsertionTime) "2025-01-01T" }}delete{{ end }}`},
			stdOut: "Matched 1 message(s).",
			verify: func(t *testing.T) {
				c := fetchTestQueueMessageCount(t, "test-delete-02")
				if c != 0 {
					t.Errorf("got %d messages in queue, want 0", c)
				}
			},
		},
		{
			name: "Delete users with Id > 5",
			args: []string{"delete", "--use-storage-emulator", "-q=test-delete-03", "--decode-json",
				"--script", `{{ $uid := .MessageJson.id | int }}{{ if gt $uid 5 }}delete{{ end }}`},
			stdOut: "Matched 1 message(s).",
			verify: func(t *testing.T) {
				c := fetchTestQueueMessageCount(t, "test-delete-03")
				if c != 1 {
					t.Errorf("got %d messages in queue, want 1", c)
				}
			},
		},
	}

	executeTestCases(t, tc)
}
