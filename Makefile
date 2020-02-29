GO = go
ENV = export GOPATH=$(CURDIR)/lib/src
GOGET = export GOPATH=$(CURDIR)/lib/src; $(GO) get
BUILD = $(ENV); $(GO) build

tst:
	@echo "$(CURDIR)"

dep:

	$(GOGET) github.com/BurntSushi/toml
	$(GOGET) github.com/gorilla/websocket
	$(GOGET) gopkg.in/doug-martin/goqu.v3
	$(GOGET) gopkg.in/doug-martin/goqu.v3/adapters/mysql
	$(GOGET) github.com/go-sql-driver/mysql
	$(GOGET) github.com/google/gops
	$(GOGET) gopkg.in/iconv.v1

tc-dep: dep
	$(GOGET) github.com/golang/protobuf/proto
	$(GOGET) github.com/streadway/amqp
	$(GOGET) github.com/assembla/cony

mtapitest:
	$(BUILD) -o bin/mtapitest.exe apps/apitest.go

mtdealerapi: dep
	$(BUILD) -o bin/mtdealerapi.exe apps/dealerapi.go

mtreportapi: dep
	$(BUILD) -o bin/mtreportapi.exe apps/reportapi.go
	$(BUILD) -o bin/mtreportapi apps/reportapi.go

tcapi: dep
	$(BUILD) -o bin/tcapi apps/tcapi.go

tcdealer: dep
	$(BUILD) -o bin/tcdealer apps/tcdealer.go

gops: dep
	$(BUILD) -o bin/gops github.com/google/gops

tc-dealer: tc-dep

clean:
	rm -rf ./bin/*
