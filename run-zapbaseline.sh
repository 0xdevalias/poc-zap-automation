#!/bin/sh

if [ -z "$1" ]; then
  echo "Usage: $0 http://example.com/"
  exit
fi

export SCANTARGET=$1
export SCANHOST=`sed -E -e 's_.*://([^/@]*@)?([^/:]+).*_\2_' <<< $SCANTARGET`

docker-compose -f docker-compose-run.yml run --rm zapbaseline
