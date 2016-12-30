LIBDIR ?= usr/lib
SBINDIR ?= usr/sbin

all: # nothing to do

install:
	mkdir -p $(DESTDIR)/$(SBINDIR)
	cp openshift-pre-backup $(DESTDIR)/$(SBINDIR)
	chmod 755 $(DESTDIR)/$(SBINDIR)/openshift-pre-backup
