#!/bin/bash

set -eo pipefail

WORKVOL=`pwd`/work-vol
mkdir -p ${WORKVOL}

if [ ! -e stage0.tar ]; then
	# setup build env
	docker build -f docker/Dockerfile.stage-1 -t wifi:raspbian-stage-1 docker/
	# setup binfmt_misc stuff
	docker run --privileged --rm -t -i wifi:raspbian-stage-1 "/stage-1.sh"
	# prefetch
	docker build -f docker/Dockerfile.stage0 -t wifi:raspbian-stage0 docker/

	echo removing ${WORKVOL}...
	sudo rm -rf "${WORKVOL}"
	docker run --privileged -v ${WORKVOL}:/work  --rm -t -i wifi:raspbian-stage0 /stage0.sh
	echo saving tar...
	pushd work-vol; sudo tar cvf ../stage0.tar *; popd
	sudo chown `whoami` stage0.tar
fi

# refresh workdir
if [ ! -e stage1.tar ]; then
	echo Unpack stage0.tar
	pushd "${WORKVOL}"; sudo tar xvf ../stage0.tar; popd
	docker build -f docker/Dockerfile.stage1 -t wifi:raspbian-stage1 docker/
	docker run --privileged -v ${WORKVOL}:/work  --rm -t -i wifi:raspbian-stage1 /stage1.sh
	echo saving tar...
	pushd work-vol; sudo tar cvf ../stage1.tar *; popd
	sudo chown `whoami` stage1.tar
fi

# refresh workdir
if [ ! -e stage2.tar ]; then
	echo Unpack stage1.tar
#	pushd "${WORKVOL}"; sudo tar xvf ../stage1.tar; popd
	docker build -f docker/Dockerfile.stage2 -t wifi:raspbian-stage2 docker/
	docker run --privileged -v ${WORKVOL}:/work  --rm -t -i wifi:raspbian-stage2 /stage2.sh
	echo saving tar...
	pushd work-vol; sudo tar cvf ../stage2.tar *; popd
	sudo chown `whoami` stage2.tar
fi

# copy over
# * mbrola
# * amixer volume stuff

# firmware for internal chipset
# XXX this isn't good
# cp ../firmware-nonfree/brcm80211/brcm/brcmfmac43430-sdio.* work-vol/now-raspbian.img/stage2/rootfs/lib/firmware/brcm/

# XXX get rid of the sudo stuff. yuck

#ROOTFS="work-vol/now-raspbian.img/stage2/rootfs"

ROOTIMG="work-vol/now-raspbian.img/export-image/now-raspbian.img-lite.img"
mkdir -p rootfs
echo abad
offsec=`( /sbin/fdisk ${ROOTIMG} 2>/dev/null  <<<"p" || true ) | grep img2 | awk ' { print $2 } '`
echo $offsec
offset=`expr $offsec \* 512`
echo $offset
sudo losetup -o "$offset"  /dev/loop4 "$ROOTIMG" 
echo ok losetup
/sbin/losetup
sudo mount /dev/loop4 rootfs/
echo mounted

# setup datacard mountpoint
FSTAB=rootfs/etc/fstab
sudo mkdir -p rootfs/mnt/datacard
grep datacard "$FSTAB" || (sudo bash -c "echo /dev/sda1 /mnt/datacard auto defaults,noatime 0 1 >>${FSTAB}")

cat rootfs/etc/fstab

# copy over daemon stuff
sudo mkdir -p rootfs/opt/owa/
sudo cp ../../daemon/* rootfs/opt/owa/
sudo ln -sf /opt/owa/owa.init rootfs/etc/init.d/owa

# add to run levels
for a in rootfs/etc/rc3.d/S05owa rootfs/etc/rc4.d/S05owa rootfs/etc/rc5.d/S05owa; do
	sudo ln -sf /etc/init.d/owa $a
done

sudo umount rootfs/
sudo losetup -d /dev/loop4
rmdir rootfs
