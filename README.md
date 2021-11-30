# LΛTENCY

A tool to measure the latency of messages published on the Ably Realtime Network. It utilises a structure of commands, arguments & flags. It supports Fully POSIX-compliant flags (including short & long versions)

## Requirements

Latency requires the following environment variables to be set: 

`ABLY_KEY` The Ably key used for authentication

----

## Commands

[**publish**] - Publishes a message on the Ably Realtime Network and prints the latency of the message

`latency publish --channel MAIN --event TEST`

[**subscribe**] - Subscribes to the Ably Realtime Network and prints the latency of the message

`latency subscribe --channel MAIN --event TEST`

[**help**] - Prints the help text for the command

---

### Basic Usage 


#### Build:
```
$ git clone https://github.com/tomweston/latency
$ cd latency
$ docker build -t latency .
```

### Run:
```
$ docker run -it latency
```

[**publish**]: https://github.com/tomweston/latency#commands
[**subscribe**]: https://github.com/tomweston/latency#commands
[**help**]: https://github.com/tomweston/latency#commands
