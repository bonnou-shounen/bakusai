#!/bin/bash

cd `dirname $0`/..
mkdir -p tmp
tmpf=tmp/tmp.txt

dump="./bakusai dump-thread"

uri="$1"

if [[ -z "$uri" ]]; then
    exit
fi

while true; do
    $dump $uri > tmp/tmp.txt
    if [[ $? != 0 ]]; then
        exit
    fi

    title=`head -5 $tmpf | sed -n 's/^# T: //p'`
    name=`head -5 $tmpf | sed -n 's/^# C: //p' | sed 's/[\/ :]//g'`
    prev=`head -5 $tmpf | sed -n 's/^# P: //p'`

    mv $tmpf tmp/$name.txt

    echo "$name $title"

    if [[ -z "$prev" ]]; then
        break
    fi
    uri="$prev"

    sleep 1
done
