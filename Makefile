IMG="platform/rpi3/work-vol/now-raspbian.img/export-image/now-raspbian.img-lite.img"
ROOTFS="platform/rpi3/work-vol/now-raspbian.img/stage2/rootfs"

ifndef OWA_DEV
 $(error OWA_DEV not defined)
endif
ifndef OWA_MNT
 $(error OWA_MNT not defined)
endif
#OWA_DEV="/dev/sdi"
#OWA_MNT="/mnt/owa"

install-sync: $(IMG)
	mount ${OWA_DEV}2 /mnt/owa
	rsync -aHAXx ${ROOTFS_DIR}/ /mnt/owa/

install-image: $(IMG)
	dd if=$(IMG) of=$(OWA_DEV)

$(IMG):
	cd platform/rpi3/ && ./host_run.sh


# Not needed yet. rpi wifi is no good.
#
#platform/firmware-nonfree:
#	cd platform && git clone https://github.com/RPi-Distro/firmware-nonfree
#
#platform/linux-firmware:
#	cd platform && git clone https://git.kernel.org/pub/scm/linux/kernel/git/firmware/linux-firmware.git
	
