.PHONY:	all clean vendor api-gen
export CGO_ENABLED:=0

REPO=github.com/flatcar-linux/updateservicectl
ROLLER_URL ?= http://localhost:8000

all:
	go build -mod vendor -o bin/updateservicectl

vendor:
	GO111MODULE=on go mod vendor

clean:
	rm -rf bin

api-gen:
	# curl $(ROLLER_URL)/_ah/api/discovery/v1/apis/update/v1/rest > client/update/v1/update-api.json
	google-api-go-generator -api_json_file client/update/v1/update-api.json -api_pkg_base $(REPO)/client -gendir client
