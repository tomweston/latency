# LΛTENCY

## Overview

Latency is a tool to measure the latency of messages published on the Ably Realtime Network. It utilises a structure of commands, arguments & flags. It supports Fully POSIX-compliant flags (including short & long versions)

## Requirements

Latency requires the following environment variables to be set: 

`ABLY_KEY` The Ably key used for authentication

`ABLY_CHANNEL` The Ably Realtime channel used for communication

`ABLY_EVENT` The Ably Event used for the service

----

## Commands

[**publish**] - Publishes a message on the Ably Realtime Network and prints the latency of the message

`latency publish`

[**subscribe**] - Subscribes to the Ably Realtime Network and prints the latency of the message

`latency subscribe`

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
