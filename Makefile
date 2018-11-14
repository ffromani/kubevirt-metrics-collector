VERSIONFILE := internal/pkg/version/version.go

all: binary

docker: binary
	docker build .

dockertag: binary
	docker build -t fromanirh/kube-metrics-collector:devel .

vendor:
	dep ensure

binary: vendor gensrc
	cd cmd/kube-metrics-collector && go build -v .

clean:
	rm -f cmd/kube-metrics-collector/kube-metrics-collector

gensrc:
	@./genver.sh > $(VERSIONFILE)

.PHONY: all docker binary clean

