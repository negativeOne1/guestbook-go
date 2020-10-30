CMD_DIR           = cmd
BIN_DIR           = bin
CMDS              = $(subst $(CMD_DIR)/,,$(wildcard $(CMD_DIR)/*))
LDFLAGS           = -ldflags="-s -w"

GOCMD             = go
GOBUILD           = $(GOCMD) build
ECHO              = @echo

ifneq ($(V),1)
	Q  = @
endif

.phony: all
all: $(CMDS)


clean: ## Remove previous build
	$(ECHO) "  CLEAN"
	$(Q)rm -f $(CMD_DIR)/*

$(CMDS):
	$(Q)$(ECHO) "  GO" $@
	$(Q)$(GOBUILD) -o $(BIN_DIR)/$@ $(CMD_DIR)/$@/main.go
