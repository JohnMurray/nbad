BUILD_OPTS=-p 4
BIN_NAME=nbad

default: compile

compile: setup
	go build $(BUILD_OPTS) -o $(BIN_NAME)

setup:
	go get -u github.com/Syncbak-Git/nsca

clean-compile: BUILD_OPTS += -a
clean-compile: compile
