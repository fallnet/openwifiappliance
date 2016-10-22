#!/bin/bash

set -eo pipefail

source /config

if [ -e /pi-gen/work ]; then
	cp -r /pi-gen/work/* /work/
	rm -rf /pi-gen/work
fi
ln -s  /work /pi-gen/work
cd /pi-gen
RUN_STAGE=0 ./build.sh
