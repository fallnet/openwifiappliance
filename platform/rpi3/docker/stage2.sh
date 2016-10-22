#!/bin/bash

set -eo pipefail

source /config

rm -rf /pi-gen/work
ln -s  /work /pi-gen/work
cat /dev/null >/pi-gen/export-image/01-set-sources/00-patches/0-sources.diff
cat /dev/null >/pi-gen/stage2/01-sys-tweaks/00-patches/02-swap.diff
cat /dev/null >/pi-gen/stage2/01-sys-tweaks/00-patches/03-console-setup.diff
cat /dev/null >/pi-gen/export-noobs/prerun.sh
cat /dev/null >/pi-gen/export-noobs/00-release/00-run.sh
cd /pi-gen
MAX_STAGE=2 ./build.sh
