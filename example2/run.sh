#!/bin/sh
#set -e

SERVER=127.0.0.1:8123
COUNT=300

killall example2
rm -f ./example2 /tmp/example.log*

echo "Starting server..."
go build .
./example2 &


echo "Starting requests...."
./curl.sh $SERVER $COUNT &
sleep 1

echo "Rotating..."
mv /tmp/example.log /tmp/example.log-old
killall -1 example2

# wait for curl to finish
wait %2

killall -TERM example2

# count results
wc -l /tmp/example.log*



# exit