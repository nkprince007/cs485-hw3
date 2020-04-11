export UDIR= .
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
SUBMISSION_OUT := sangi.NaveenKumar.Assign3
TAR_OUT := $(SUBMISSION_OUT).tar

.PHONY: all install clean
all: $(BINARIES)

ethosChat.go: EthosChat.t
	$(ETN2GO) . ethosChat main $^

ethosChatClient: ethosChatClient.go ethosChat.go
	ethosGo $^

ethosChatService: ethosChatService.go ethosChat.go
	ethosGo $^

install: $(BINARIES)
	rm -fr $(SERVERROOT)
	ethosParams $(SERVERROOT)
	(cd $(SERVERROOT) && ethosMinimaltdBuilder)
	ethosTypeInstall ethosChat
	ethosDirectoryInstall user/nobody/chatrooms $(ETHOSROOT)/types/spec/ethosChat/ChatRoom all
	ethosDirCreate $(ETHOSROOT)/services/ethosChat $(ETHOSROOT)/types/spec/ethosChat/ChatRpc all
	install -D ethosChatClient ethosChatService $(ETHOSROOT)/programs

clean:
	rm -fr ethosChatService ethosChatClient
	rm -fr *.goo.ethos
	rm -fr $(SERVERROOT)
	rm -fr ethosChat ethosChatIndex ethosChat.go
	rm -fr $(TAR_OUT)

archive: clean
	mkdir -p $(SUBMISSION_OUT)
	cp *.t $(SUBMISSION_OUT)
	cp *.go $(SUBMISSION_OUT)
	cp Makefile $(SUBMISSION_OUT)
	tar -cvf $(TAR_OUT) $(SUBMISSION_OUT)
	rm -fr $(SUBMISSION_OUT)
