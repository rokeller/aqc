# Azure Queue Commands (aqc)

A simple CLI to mess with Azure storage queues.

## Usage

`aqc` offers the following commands:

* `add` - add messages to a queue.
* `clear` - clear all messages from a queue.
* `delete` - selectively delete messages from a queue. See more [below](#delete-messages).
* `move` - move messages from one queue to another.
* `peek` - peek messages in a queue.

You can use the `-h` on any command to get more information:

```bash
aqc -h
# or
aqc add -h
# etc.
```

## Delete messages

Deleting messages from a queue currently is the only command that supports
_scripting_ using [go templates](https://pkg.go.dev/text/template). The way this
works is kind of simple. The user provides a _script_ (a go template) that is
evaluated for every message, and if the script produces any non-empty output,
that is interpreted as the decision to delete the message.

An object with the following (case-sensitive) properties is made available to
the _script_ go template for every message:

* `DequeueCount` (int64) - the number of times the message has already been
  dequeued.
* `InsertionTime` (timestamp) - the date/time when the message was added to the
  queue.
* `ExpirationTime` (timestamp) - the date/time when the message will expire and
  therefore will be removed from the queue.
* `MessageID` (string) - the unique ID (UUID) of the message in the queue.
* `MessageText` (string) - the text of the message in the queue. If the `-b` /
  `--decode-base64` flag is set to true and the message is base64 encoded, this
  will be the base64 _decoded_ value.
* `MessageJson` (JSON value) - the decoded JSON value if and only if the `-j` /
  `--decode-json` flag is set to true, otherwise (or if JSON decoding fails) this
  will be null (`nil`).

### Additional functions

Aside from the [functions available with go templates out-of-the-box](https://pkg.go.dev/text/template#hdr-Functions),
`aqc` makes the following functions available for use in _scripts_.

* `lower(string)` - convert string input to lower-case.
* `upper(string)` - convert string input to upper-case.
* `iso(timestamp)` - convert timestamp input to its ISO/RFC3339 representation
  (i.e. like `2025-06-28T12:34:02Z`). This makes timestamps more suitable for
  comparison to check before/after kind of conditions.
* `int(float)` - convert a float input to an integer to allow integer comparison.
  This can be very useful for numeric JSON property values, which due to lack of
  schema details will always be intepreted as float values.

### Example

Let's assume for a moment that all messages in the queue should be JSON encoded
objects representing users with a `username` property (string) and an `id`
property (int64) that is used to add users to a store/remote system/whatever.
That is, you'll see JSON encoded objects like the following in the queue:

```json
{
    "username": "foo@bar.com",
    "id": 123456
}
```

If now you want to remove all messages for usernames starting with `a`, you could
use `aqc` as follows. Note that we do _not_ assume that messages are base64
encoded too.

```bash
aqc delete -j --queue-url https://myaccount.queue.core.windows.net/add-users \
  --script '{{ $u := .MessageJson.username | lower }}{{ if and (ge $u "a") (lt $u "b") }}delete{{ end }}'
```

That is, the script produces the output `delete` (which is a non-empty but
otherwise irrelevant string) if and only if the `username` property of the JSON
object in the message, converted to lower-case, is greater than or equal to `a`
and is less than `b`. Thus, all messages with `username` values that start with
`a` or `A` would match and be deleted from the queue.

## Authentication

Authentication is attempted automatically according to the options listed in
[Azure authentication with the Azure Identity module for Go
](https://learn.microsoft.com/en-us/azure/developer/go/azure-sdk-authentication?tabs=bash#2-authenticate-with-azure).

The only exception to this is when a queue URL is passed with a SAS token in it.
In that case, the URL is used as-is and no other authentication mechanisms are
tried.

Running `aqc` is most easy on Azure workloads that have a managed identity
assigned to them (see options 2 and 3 on the page linked above). Using a service
principal with a secret is almost as easy, using the `AZURE_CLIENT_ID`,
`AZURE_TENANT_ID` and `AZURE_CLIENT_SECRET` environment variables (option 1 on
the page linked above). Using `aqc` while being logged in to `AzCLI` also works
without the need to provide any further credentials.
