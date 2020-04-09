export UDIR= bin
export LDIR= .
export GOC = x86_64-xen-ethos-6g
export GOL = x86_64-xen-ethos-6l
export ETN2GO = etn2go
export EG2GO = eg2go
export ET2G = et2g

export GOARCH = amd64
export TARGET_ARCH = x86_64
export GOETHOSINCLUDE=/usr/lib64/go/pkg/ethos_$(GOARCH)
export GOLINUXINCLUDE=/usr/lib64/go/pkg/linux_$(GOARCH)

export SERVERROOT=server
export ETHOSROOT=$(SERVERROOT)/rootfs
export MINIMALTDROOT=$(SERVERROOT)/minimaltdfs

SRC_FILES = $(wildcard *.go)
BINARIES = $(patsubst %.go,%,$(SRC_FILES))
TYPES = StructA StructB StructResult
SUBMISSION_OUT := sangi.NaveenKumar.Assign2
TAR_OUT := $(SUBMISSION_OUT).tar

.PHONY: all install clean
all: ethosChat.go $(BINARIES)

ethosChat.go: EthosChat.t
	$(ETN2GO) . ethosChat main $^

clean:
	rm -fr $(UDIR)
	rm -fr $(SERVERROOT)
	rm -fr ethosChat ethosChatIndex ethosChat.go
