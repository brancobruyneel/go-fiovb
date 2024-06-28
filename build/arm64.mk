################################################################################
### FIOVB - ARM64                                                            ###
################################################################################

ARM64 = arm64
BIN_ARM64 = $(OUTPUT)/$(ARM64)/bin
PREFIX_ARM64 = CGO_ENABLED=1 GOOS=linux GOARCH=$(ARM64)

arm64: fiovb-tool-arm64
arm64-clean: fiovb-tool-clean-arm64

################################################################################
### FIOVB TOOL                                                               ###
################################################################################

fiovb-tool-arm64:
	@$(PREFIX_ARM64) $(GO) build -o $(BIN_ARM64)/fiovb github.com/brancobruyneel/go-fiovb/cmd/fiovb/
	@echo Compiled $@

fiovb-tool-clean-arm64:
	@rm -rf $(BIN_ARM64)/fiovb
	@echo Cleaned fiovb-tool-arm64
