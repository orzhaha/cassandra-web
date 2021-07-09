include version.mk

APP=cassandra-web
REPOSITORY ?=carrefourphx
IMG ?=$(REPOSITORY)/$(APP):$(BUILD_VERSION)
LATEST_IMG ?=$(REPOSITORY)/$(APP):latest

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

source-env:
	@gpg --import /gpg/gpg.public || :
	@gpg --allow-secret-key-import --import /gpg/gpg.private || :
	@echo `cat /gpg/gpg.ownertrust` | gpg --import-ownertrust
	@git-crypt unlock
	$(eval REPOSITORY := `grep 'REPOSITORY=' env.secret | sed 's/REPOSITORY=//'`)

# Build the docker image
docker-build: source-env
	docker build . -t ${IMG}
	docker tag ${IMG} ${LATEST_IMG}

# Push the docker image
docker-push: source-env
	docker push ${IMG}
	docker push ${LATEST_IMG}
