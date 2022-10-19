# Cloudwatch Logs Utility

This repo has some basic utilities for working with with Cloudwatch Logs.

## getlogstream

`getlogstream` is the first utility and it will get the entire log stream and send it to STDOUT.
From there you can of course pipe it to a file etc...

### installation

To install it you can run

```bash
go install github.com/Reisender/cloudwatchlogs/cmd/getlogstream@latest
```

Assuming your `go/bin` dir is in your path, you can just run it like so

```bash
getlogstream --help
```
