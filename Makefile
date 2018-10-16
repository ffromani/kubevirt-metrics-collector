all: binary

docker: binary
	docker build .

vendor:
	dep ensure

binary: vendor
	cd cmd/kube-metrics-collector && go build -v .

clean:
	rm -f cmd/kube-metrics-collector/kube-metrics-collector

.PHONY: all docker binary clean

