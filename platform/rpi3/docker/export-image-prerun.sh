#!/bin/bash -e
IMG_FILE="${STAGE_WORK_DIR}/${IMG_DATE}-${IMG_NAME}${IMG_SUFFIX}.img"

unmount_image ${IMG_FILE}

rm -f ${IMG_FILE}

rm -rf ${ROOTFS_DIR}
mkdir -p ${ROOTFS_DIR}

BOOT_SIZE=$(du -sh ${EXPORT_ROOTFS_DIR}/boot -B M | cut -f 1 | tr -d M)
TOTAL_SIZE=$(du -sh ${EXPORT_ROOTFS_DIR} -B M | cut -f 1 | tr -d M)
IMG_MB=$(expr $BOOT_SIZE \* 2 \+ $TOTAL_SIZE \+ 512)
IMG_SIZE="${IMG_MB}"M

fallocate -l ${IMG_SIZE} ${IMG_FILE} || dd if=/dev/zero of=$IMG_FILE bs=1M count=${IMG_MB}


fdisk ${IMG_FILE} > /dev/null 2>&1 <<EOF
o
n


8192
+`expr $BOOT_SIZE \* 3`M
p
t
c
n


8192


p
w
EOF

LOOP_DEV=`kpartx -asv ${IMG_FILE} | grep -E -o -m1 'loop[[:digit:]]+' | head -n 1`
BOOT_DEV=/dev/mapper/${LOOP_DEV}p1
ROOT_DEV=/dev/mapper/${LOOP_DEV}p2

mkdosfs -n boot -S 512 -s 16 -v $BOOT_DEV > /dev/null
mkfs.ext4 -O ^huge_file $ROOT_DEV > /dev/null
tune2fs -c 0 $ROOT_DEV

mount -v $ROOT_DEV ${ROOTFS_DIR} -t ext4
mkdir -p ${ROOTFS_DIR}/boot
mount -v $BOOT_DEV ${ROOTFS_DIR}/boot -t vfat

rsync -aHAX --one-file-system --progress ${EXPORT_ROOTFS_DIR}/ ${ROOTFS_DIR}/
