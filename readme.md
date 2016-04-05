# Nagios Buffering Agent daemon (NBAd)

[![Build Status](https://travis-ci.org/JohnMurray/nbad.svg?branch=master)](https://travis-ci.org/JohnMurray/nbad)
[![Go Report](https://goreportcard.com/badge/github.com/johnmurray/nbad)](https://goreportcard.com/report/github.com/johnmurray/nbad)

_This project is an early-stage work-in-progress! Do not use for production use._

A small application that emulates the Nagios NSCA agent locally and acts as
a buffering agent to perform "smarter" actions on behalf of the application
raising alerts.

## Problem

+ Passive checks that don’t clear themselves. This causes a lot of
  burden on pager to manually resolve these issue. Also, this results
  in alerts that cause additional noise.
+ Transient alerts (alert and quickly clear)
+ Flooding of passive checks when things go wrong can cause additional
  stress on nagios server to process input.


## Solution

The solution to this is to _ideally_ have smarter applications that take care of OK'ing
alerts when a problem is gone or throttling the number of alerts that it sends when things
are actually going wrong. This could be perhaps solved by better tooling or framework support
for whatever the application is written in. However there are a lot of languages out there
and even more frameworks and libraries. By solving this from an _external_ approach, we
get a solution that everyone can utilize.

How NBA _attmpts_ to solve this problem

+ Automatically report "OK" status on alerts that become "stale"
  that previously reported an error. (likely resolved)
+ Buffer duplicate alerts to reduce noise / spam to monitoring server


__Possible Future Additions__

+ Define a threshold that must be met before an error condition is propagated up-stream
+ Use a config to define default and per-service behaviors.



## Building / Running

I currently do not have the application pre-built for you or packaged in any distribution-specific
format. You need to use the old-fashion style of checking out the code and building from source. Luckily
this project is built fairly easily.

Before running the commands below, make sure you have Go 1.6 installed.

```bash
git clone https://github.com/JohnMurray/nbad.git
cd nbad

make setup && make
```

This will compile the program into a single file `nbad` in the root directory of the project. You
can then run via

```bash
./nbad
```

Or you can review the full set of command line options by passing in the `-h` flag.

```bash
./nbad -h
NAME:
   nbad - NSCA Buffering Agent (daemon) - Emulates NSCA interface as local buffer/proxy

USAGE:
   nbad [global options] command [command options] [arguments...]

VERSION:
   1.0

COMMANDS:
GLOBAL OPTIONS:
   --config, -c "/etc/nbad/conf.json"	Location of config file on disk
   --trace, -t				Turn on trace-logging [$NBAD_TRACE]
   --help, -h				show help
   --version, -v			print the version
```

## Configuration

NBAd uses a simple JSON configuration file that is required to run. By default NBAd looks for this configuration file at `/etc/nbad/conf.json`. You can override this with the `-c` (`--config`) flag.

The configuration options are

Value|Type|Description
-----|----|-----------
gateway_message_buffer_size|unsigned int|The number of messages to buffer in memory for the gateway
message_cache_ttl_in_seconds|unsigned int|The time before a message expires (possibly causing upstream state changes)
message_init_buffer_ttl_in_seconds|unsigned int|The amount of time a message is buffered before actioned upon


## Testing / Debugging

There is a small shell script in the repository `send_nsca.sh` that mimics the regular
`send_nsca` command in a very small way. It uses `echo` and `nc` to send messages in the
nsca v3 format to `localhost:5667` (default nbad port). To see the script options you can
simply run with the `--help` flag.

```
λ ./send_nsca.sh --help
./send_nsca.sh [-e E -h H -s S -m M|--help]

  --help              Print this error message
  -e, --error-code    Error code value (0, 1, 2)
  -h, --host          Host that the check is originating from
  -s, --service       Service that the check is attached to
  -m, --message       Check description message
```


Some examples of how to use the script (also it's a very simple script so you can always
just crack open the source).

```
# send OK
./send_nsca.sh  -e 0 -h my-service-host -s my-service -m "everything is cool"

# send WARNING
./send_nsca.sh  -e 1 -h my-service-host -s my-service -m "things are warming up"

# send CRITICAL
./send_nsca.sh  -e 2 -h my-service-host -s my-service -m "oh god! everything is on fire!!!!"
```

The script will take care of doing some basic validations on the input data in terms of length
and what not as well as properly padding and forming the message. I'm not certain how this would
work with non-ascii data, so... on you're own there.



## TODO

[ ] Enable debug and trace logging from the command line
  [ ] Review logging to ensure completeness and proper log levels
[ ] Mimic nsca server better / more
  [ ] Implement CRC
  [ ] return initialization packet with IV and Timestamp in initial server connect
[ ] Flush messages on some cache interval (to avoid transient error conditions)
[ ] HTTP / RESTful interface
