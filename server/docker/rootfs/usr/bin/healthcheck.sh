#!/bin/sh

for service in hbbr hbbs api; do
  [ "$(/command/s6-svstat -o up /run/s6-rc/servicedirs/$service)" = "true" ] || exit 1
done
wget -q -T 4 -O /dev/null http://127.0.0.1:21114/api/ || exit 1
