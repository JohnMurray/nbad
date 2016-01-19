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

There isn't much in the way of testing at the moment. Latest test command I've
been using (disable byte-length check in message.go)

```bash
echo -n -e '\x00\x03\x00\x00\x11\x12\x13\x14\x11\x12\x13\x14\x00\x01hello\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00service-name\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00something has went horribly wrong!!' | nc localhost 5667
```