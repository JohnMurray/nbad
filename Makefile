BUILD_OPTS=-p 4 -race
BIN_NAME=nbad

default: compile

compile:
	go build $(BUILD_OPTS) -o $(BIN_NAME)
	go vet
	golint .
	@gotags -tag-relative=true -R=true -sort=true -f="tags" -fields=+l .


# meta-task for performing all setup tasks
setup: get-deps
	@cp etc/pre-push-git-hook .git/hooks/pre-push

get-deps:
	go get -u github.com/Syncbak-Git/nsca
	go get -u github.com/golang/lint/golint
	go get -u github.com/codegangsta/cli
	go get -u github.com/jstemmer/gotags

clean-compile: BUILD_OPTS += -a
clean-compile: compile
