# Slack Emitter

## Idea #1

- connect to slack's RTM API (https://github.com/nlopes/slack/blob/master/examples/websocket/websocket.go)
- filter out channels that start with a given prefix ("#mt-", maybe make it configurable)
- initially just get some stats to the CLI
- next, optionally send it to an SQS queue as commands

... and that's a fail. We don't get messages other than to channels we're members of slash invited to.

