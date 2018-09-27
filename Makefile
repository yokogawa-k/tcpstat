PKG = github.com/yokogawa-k/tcpstat
COMMIT = $$(git describe --tags --always)
BUILD_LDFLAGS = -X $(PKG).main.commit=$(COMMIT)
RELEASE_BUILD_LDFLAGS = -s -w $(BUILD_LDFLAGS)

.PHONY: build
build:
	go build -ldflags="$(BUILD_LDFLAGS)"

test:
	go test -test.v -cover -coverprofile=.profile.cov -covermode=count -timeout=30m -parallel=4 $(MAKEFILE_DIR)
	go tool cover -func=.profile.cov

lint: fmt
	gometalinter --vendor --skip=vendor/ --cyclo-over=16 --disable=gas --disable=maligned --disable=gosec --deadline=2m .

lint-fast: fmt
	gometalinter --fast --vendor --skip=vendor/ --cyclo-over=16 --disable=gas --disable=maligned --disable=gosec --deadline=2m .

tools:
	go get -v github.com/alecthomas/gometalinter
	gometalinter --install

.PHONY: all fmt test lint lint-fast cover tools
