# Adapted from https://github.com/rust-embedded/rust-raspi3-OS-tutorials/blob/6a03ddea11b740323d01a478dcca3ed726cbbba2/docker/raspi3-openocd/Dockerfile
# The above file is MIT licensed.

FROM debian:stretch

LABEL maintainer="Ashley Davis <ashley_davis10419@hotmail.com>"

RUN set -ex ; \
	tempPkgs='          \
		automake        \
		build-essential \
		ca-certificates \
		git             \
		libtool         \
		pkg-config      \
		'; \
	apt-get update ; \
	apt-get install -q -y --no-install-recommends $tempPkgs libusb-1.0.0-dev ; \
	git clone --depth 1 https://github.com/riscv/riscv-openocd openocd; \
	cd openocd ; ./bootstrap ; ./configure --enable-ftdi ; make -j4 ; \
	make install ; \
	apt-get purge -y --auto-remove $tempPkgs ; \
	apt-get autoremove -q -y ; \
	apt-get clean -q -y ; \
	rm -rf /var/lib/apt/lists/*

COPY hifive1.cfg /openocd/

ENTRYPOINT ["openocd"]
CMD ["-f", "/openocd/hifive1.cfg"]
