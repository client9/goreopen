#!/bin/sh
server=$1
count=$2
a=1
while [ $a -le $count ]
do
   echo $a
   curl --silent "$server/$a" > /dev/null
   a=`expr $a + 1`
done