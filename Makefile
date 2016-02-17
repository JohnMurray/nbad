BUILD_OPTS=-p 4
BIN_NAME=nbad

default: compile

compile: setup
	go build $(BUILD_OPTS) -o $(BIN_NAME)
	go vet -x
	golint .


# meta-task for performing all setup tasks
setup: get-deps

get-deps:
	go get -u github.com/Syncbak-Git/nsca
	go get -u github.com/golang/lint/golint

clean-compile: BUILD_OPTS += -a
clean-compile: compile
