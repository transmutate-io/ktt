# The name of the executable (default is current directory name)
TARGET_BIN := build/ktt
BASE_PKG := transmutate.io/pkg/ktt
TARGET_PKG := $(BASE_PKG)/cmd/ktt


VERSION := `git describe --always --tags HEAD`
COMMIT := `git rev-parse HEAD`

# Use linker flags to provide version/commit settings to the target
LDFLAGS=-ldflags "-X=$(BASE_PKG)/cmd/ktt/cmd.Version=$(VERSION) -X=$(BASE_PKG)/cmd/ktt/cmd.Commit=$(COMMIT)"

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all

all: build

build: $(TARGET_BIN)

$(TARGET_BIN): $(SRC)
	@go build -o $(TARGET_BIN) $(LDFLAGS) $(TARGET_PKG)

clean:
	@rm -rf $(TARGET_BIN)

install:
	@go install $(LDFLAGS) $(TARGET_PKG)

test:
	@go test ./...