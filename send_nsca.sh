#!/bin/bash

# Simple script to emulate send_nsca in a very simple way by generating some
# (hopefully) sane messages and sending them to the local nbad service.
#
# Ex:
#    go build && nohup ./nbad &
#    ./send_nsca.sh -e 3 -h my-server -s critical-service -m "things be bad!"


# Parse command line arguments
while [[ $# > 1 ]] ; do
    key="$1"

    case $key in
        -e|--error-code)
        ERROR_CODE="$2"
        shift # past argument
        ;;
        -h|--host)
        HOST="$2"
        shift # past argument
        ;;
        -s|--service)
        SERVICE="$2"
        shift # past argument
        ;;
        -m|--message)
        MESSAGE="$2"
        shift # past argument
        ;;
        *)
                # unknown option
        ;;
    esac
    shift # past argument or value
done

# Some basic validations of input
host_len=${#HOST}
if [[ $host_len -gt 64 ]] ; then
    echo "Host cannot be more than 64 characters long"
    exit 1
fi

srv_len=${#SERVICE}
if [[ $srv_len -gt 128 ]] ; then
    echo "Service description cannot be more than 128 characters long"
    exit 1
fi

msg_len=${#MESSAGE}
if [[ $msg_len -gt 512 ]] ; then
    echo "Error message can not be more than 512 characters long"
    exit 1
fi


# Define some constants for working with the message format
VERSION='\x00\x03'                 # v3 Nagios message
PADDING='\x00\x00'                 # Part of binary protocol
CRC='\x11\x12\x13\x14'             # Fake values, not currently enforced
TIMESTAMP='\x11\x12\x13\x14'       # Fake values, not currently enforced
ERROR='\x00\x0'$ERROR_CODE

i=0
to=$((64 - $host_len))
while [ $i -lt $to ] ; do
    HOST=$HOST'\x00'
    let i=$i+1
done

i=0
to=$((128 - $srv_len))
while [ $i -lt $to ] ; do
    SERVICE=$SERVICE'\x00'
    let i=$i+1
done

i=0
to=$((512 - $msg_len))
while [ $i -lt $to ] ; do
    MESSAGE=$MESSAGE'\x00'
    let i=$i+1
done

echo -n -e $VERSION$PADDING$CRC$TIMESTAMP$ERROR$HOST$SERVICE$MESSAGE$PADDING | nc localhost 5667
