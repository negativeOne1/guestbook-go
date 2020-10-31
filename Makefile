CMD_DIR           = cmd
BIN_DIR           = bin
CMDS              = $(subst $(CMD_DIR)/,,$(wildcard $(CMD_DIR)/*))
LDFLAGS           = -ldflags="-s -w"

GOCMD             = go
GOBUILD           = $(GOCMD) build
ECHO              = @echo
CGO               = CGO_ENABLED=0 GOOS=linux

ifneq ($(V),1)
	Q = @
endif


.phony: all
all: $(CMDS)

clean: ## Remove previous build
	$(ECHO) "  CLEAN"
	$(Q)rm -f $(BIN_DIR)/*

$(CMDS):
	$(Q)$(ECHO) "  GO" $@
	$(Q) $(CGO) $(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$@ $(CMD_DIR)/$@/main.go

build: guestbook
	docker build -t guestbook:latest .
