################################################################################
### FIOVB - AMD64                                                            ###
################################################################################

AMD64 = amd64
BIN_AMD64 = $(OUTPUT)/$(AMD64)/bin
PREFIX_AMD64 = CGO_ENABLED=1 GOOS=linux GOARCH=$(AMD64)

amd64: fiovb-tool-amd64
amd64-clean: fiovb-tool-clean-amd64

################################################################################
### FIOVB TOOL                                                               ###
################################################################################

fiovb-tool-amd64:
	@$(PREFIX_AMD64) $(GO) build -o $(BIN_AMD64)/fiovb github.com/OpenPixelSystems/go-fiovb/cmd/fiovb/
	@echo Compiled $@

fiovb-tool-clean-amd64:
	@rm -rf $(BIN_AMD64)/fiovb
	@echo Cleaned fiovb-tool-amd64