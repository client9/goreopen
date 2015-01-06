#!/bin/sh
#set -e

SERVER=127.0.0.1:8123
COUNT=300

killall example1
rm -f ./example1 /tmp/example.log*

echo "Starting server..."
go build .
./example1 &


echo "Starting requests...."
./curl.sh $SERVER $COUNT &
sleep 1

echo "Rotating..."
mv /tmp/example.log /tmp/example.log-old
killall -1 example1

# wait for curl to finish
wait %2

killall -TERM example1

# count results
wc -l /tmp/example.log*



# exit