#!/bin/bash
mount -t binfmt_misc none /proc/sys/fs/binfmt_misc
update-binfmts --enable qemu-arm
ls -l /proc/sys/fs/binfmt_misc/
cat /proc/sys/fs/binfmt_misc/qemu-arm
