#
# Makefile
# josephzeng, 2018-08-02 23:13
#

SHELL=/bin/bash
export GOPATH=$(shell pwd):/home/josephzeng/Env/go
GOBIN=go
GLIDE=glide
GOBUILD=$(GOBIN) build
GOCLEAN=$(GOBIN) clean
GOTEST=$(GOBIN) test
GOGET=$(GOBIN) get
DEMOBIN=./bin/demo
DEMONSRC=./src/demo
ACCESSBIN=./bin/access
ACCESSSRC=./src/access


all: build_gateway build_demo

test: test_demo  test_gateway
build: build_gateway build_demo

test_demo:
	$(GOTEST) -v $(DEMONSRC)

build_demo:
	rm -rf $(DEMOBIN)
	$(GOBUILD) -o $(DEMOBIN) $(DEMONSRC)

run_demo: build_demo
	$(DEMOBIN)

test_access:
	rm -rf ./bin/test_gateway
	$(GOBUILD) -o ./bin/test_access ./src/test
	./bin/test_access

build_access:
	rm -rf $(ACCESSBIN)
	$(GOBUILD) -o $(ACCESSBIN) $(ACCESSSRC)

run_access: build_access
	#ps -ef| grep $(ACCESSBIN) | grep -v grep | awk '{print $$2}' | xargs kill -9
	$(ACCESSBIN)

deps:
	cd ./src                                               ;\
	$(GLIDE) get gopkg.in/mgo.v2                           ;\
	$(GLIDE) get github.com/rifflock/lfshook#v2.3          ;\
	$(GLIDE) get github.com/sirupsen/logrus#v1.0.5         ;\
	$(GLIDE) get github.com/go-ini/ini#v1.38.1             ;\
	$(GLIDE) get github.com/urfave/cli#v1.20.0             ;\
	$(GLIDE) get github.com/gogo/protobuf/proto            ;\
	$(GLIDE) get github.com/gogo/protobuf/jsonpb           ;\
	$(GLIDE) get github.com/gogo/protobuf/protoc-gen-gogo  ;\
	$(GLIDE) get github.com/gogo/protobuf/gogoproto        ;\
	$(GLIDE) install

clean:
	rm -rf ./bin/*
	rm -rf ./logs/*

.PHONY:test build clean deps test_demo build_demo run_demo test_access build_access run_access
