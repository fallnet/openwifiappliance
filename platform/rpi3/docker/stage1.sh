#!/bin/bash

set -eo pipefail

source /config

rm -rf /pi-gen/work
ln -s  /work /pi-gen/work
cd /pi-gen
cat /dev/null >/pi-gen/stage1/02-net-tweaks/00-patches/*persist*
MAX_STAGE=1 ./build.sh
