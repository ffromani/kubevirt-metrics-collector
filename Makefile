all: binary

docker: binary
	docker build .

dockertag: binary
	./hack/docker/tag.sh

dockerpush: binary
	./hack/docker/push.sh

vendor:
	dep ensure

binary: vendor
	./hack/build/build.sh

clean:
	rm -f cmd/kubevirt-metrics-collector/kubevirt-metrics-collector

.PHONY: all docker dockertag dockerpush binary clean

