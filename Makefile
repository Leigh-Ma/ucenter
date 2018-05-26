PROGS = oops

LDFLAGS =

.PHONY: all
all: $(PROGS)

$(PROGS):
	go install ${LDFLAGS} ucenter


