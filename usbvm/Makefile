ISO=alpine-virt-3.4.6-x86_64.iso
LINUX=linux-4.8.9
#LINUX=wireless-drivers-next
KERNEL=$(LINUX)/arch/x86/boot/bzImage
QEMUCMD=qemu-system-x86_64 -enable-kvm -nographic -net user,vlan=0 -net nic

all: usbvm.qcow2


$(ISO):
	wget -c https://fr.alpinelinux.org/alpine/v3.4/releases/x86_64/$(ISO)

$(LINUX)/Makefile: $(LINUX).tar.xz
	tar Jxvf $^
	cp config-usbvm $(LINUX)/.config

$(LINUX).tar.xz:
	wget -c https://cdn.kernel.org/pub/linux/kernel/v4.x/$(LINUX).tar.xz

$(KERNEL): $(LINUX)/Makefile
	cd $(LINUX) && make -j9 && make -j9 modules
	touch $(KERNEL)

PWD=$(shell pwd)

usbvm_shared/bzImage: $(KERNEL)
	mkdir -p usbvm_shared
	cd $(LINUX) && INSTALL_MOD_PATH=$(PWD)/usbvm_shared make -j9 modules_install
	cp $(KERNEL) usbvm_shared/

# image that contains the setup from iso
usbvm.setup.qcow2: $(ISO) install.txt
	qemu-img create -f qcow2 $@ 4G
	./applyscript install.txt | $(QEMUCMD) -hda $@ -cdrom $<

usbvm.dev.qcow2: usbvm.setup.qcow2 dev.txt
	qemu-img create -b $< -f qcow2 $@
	./applyscript dev.txt | $(QEMUCMD) -hda $@	\
	      -fsdev local,id=usbvm_dev,path=`pwd`/usbvm_shared,security_model=none	\
	      -device virtio-9p-pci,fsdev=usbvm_dev,mount_tag=usbvm_shared

usbvm_shared/diod: usbvm.dev.qcow2 diod.txt
	qemu-img create -b $< -f qcow2 usbvm.diod.qcow2
	./applyscript diod.txt | $(QEMUCMD) -hda usbvm.diod.qcow2 \
	      -fsdev local,id=usbvm_dev,path=`pwd`/usbvm_shared,security_model=none	\
	      -device virtio-9p-pci,fsdev=usbvm_dev,mount_tag=usbvm_shared
	rm usbvm.diod.qcow2

# image that contains all the installed packages
usbvm.config.qcow2: usbvm.setup.qcow2 configure.txt
	qemu-img create -b $< -f qcow2 $@
	./applyscript configure.txt | $(QEMUCMD) -hda $@

# image that contains the kernel
usbvm.kernel.qcow2: usbvm.config.qcow2 usbvm_shared/bzImage kernel.txt wlan-firmware
	cp extlinux.conf usbvm_shared/
	cp ffserver.conf usbvm_shared/
	cp -r wlan-firmware/firmware  usbvm_shared/
	qemu-img create -b $< -f qcow2 $@
	./applyscript kernel.txt | $(QEMUCMD) -hda $@ \
	      -fsdev local,id=usbvm_dev,path=`pwd`/usbvm_shared,security_model=none \
	      -device virtio-9p-pci,fsdev=usbvm_dev,mount_tag=usbvm_shared

.PHONY: run
run: usbvm.kernel.qcow2
	qemu-img create -b $^ -f qcow2 usbvm.run.qcow2
	$(QEMUCMD) -hda usbvm.run.qcow2			\
		-smp 8					\
		-usb					\
		-device usb-ehci,id=ehci			\
		-device usb-host,vendorid=0x13d3,productid=0x3327,bus=ehci.0	\
		-device usb-host,vendorid=0x1908,productid=0x2311,bus=ehci.0	\
		-device usb-host,vendorid=0x148f,productid=0x7601,bus=ehci.0	\
		-device usb-host,vendorid=0x148f,productid=0x760b,bus=ehci.0

.PHONY: run-net
run-net: usbvm.kernel.qcow2
	qemu-img create -b $^ -f qcow2 usbvm.run.qcow2
	$(QEMUCMD) -hda usbvm.run.qcow2	-smp 8 -usb -device usb-ehci,id=ehci \
	 $(shell lsusb | egrep "(ASIX|Ralink|Realtek|IMC)" | cut -f1 -d: | awk '{ print "-device usb-host,hostbus="$$2",hostaddr="$$4",bus=ehci.0" } ' | sed 's/=0*/=/g')

.PHONY: run-devs
run-devs: usbvm.kernel.qcow2
	qemu-img create -b $^ -f qcow2 usbvm.run.qcow2
	$(QEMUCMD) -hda usbvm.run.qcow2	-smp 8 -usb -device usb-ehci,id=ehci $(DEVS)


# creates an asf stream from usb webcam
.PHONY: cam
cam: usbvm.kernel.qcow2
	qemu-img create -b $^ -f qcow2 usbvm.cam.qcow2
	./applyscript cam.txt | qemu-system-x86_64 -enable-kvm -nographic \
		-net user,id=camnet,hostfwd=tcp::8090-:8090,vlan=0	\
		-net nic					\
		-hda usbvm.cam.qcow2			\
		-usb					\
		-device usb-ehci,id=ehci			\
		-device usb-host,vendorid=0x1908,productid=0x2311,bus=ehci.0



# creates 9p server exporting ar9721 spectral data on port 5640 at /mnt/dev/ath9k
# mount -t 9p 127.0.0.1 /mnt/usbvm \
#     -oversion=9p2000.L,uname=root,access=any,aname=/mnt/dev/ath9k,trans=tcp,port=5640,dfltuid=0,dfltgid=0
.PHONY: spectral
spectral: usbvm.kernel.qcow2
	qemu-img create -b $^ -f qcow2 usbvm.spectral.qcow2
	# XXX needs usb pass-through stuff
	./applyscript spectral.txt | qemu-system-x86_64 -enable-kvm -nographic	\
		-hda usbvm.spectral.qcow2 					\
		-net user,id=vm9p,hostfwd=tcp::5640-:564,vlan=0			\
		-net nic							\
		-usb								\
		-device usb-ehci,id=ehci					\
		-device usb-host,vendorid=0x13d3,productid=0x3327,bus=ehci.0

wlan-firmware:
	git clone https://github.com/OpenELEC/wlan-firmware.git

.PHONY: clean
clean:
	rm -f *.qcow2

