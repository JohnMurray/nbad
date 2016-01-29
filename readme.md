# Nagios Buffering Agent daemon (NBAd)

_This project is an early-stage work-in-progress! Do not use for production use._

A small application that emulates the Nagios NSCA agent locally and acts as
a buffering agent to perform "smarter" actions on behalf of the application
raising alerts.

## Problem

+ Passive checks that don’t clear themselves. This causes a lot of
  burden on pager to manually resolve these issue. Also, this results
  in alerts that cause additional noise.
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
- Buffer duplicate alerts to reduce noise / spam to monitoring server
- Use a config to define default and per-service behaviors.




## Testing

There is a small shell script in the repository `send_nsca.sh` that mimics the regular
`send_nsca` command in a very small way. It uses `echo` and `nc` to send messages in the
nsca v3 format to `localhost:5667` (default nbad port). You can use the script like so:

```
# send OK
./send_nsca.sh  -e 0 -h my-service-host -s my-service -m "everything is cool"

# send WARNING
./send_nsca.sh  -e 1 -h my-service-host -s my-service -m "things are warming up"

# send CRITICAL
./send_nsca.sh  -e 2 -h my-service-host -s my-service -m "oh god! everything is on fire!!!!"
```

The script will take care of doing some basic validations on the input data in terms of lenght
and what not as well as properly padding and forming the message. I'm not certain how this would
work with non-ascii data, so... on you're own there.
