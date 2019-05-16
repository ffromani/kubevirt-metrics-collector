all: binary

#vendor:
#	dep ensure

binary: #vendor
	./hack/build/build.sh ${VERSION}

release: binary
	mkdir -p _out
	cp cmd/kubevirt-metrics-collector/kubevirt-metrics-collector _out/kubevirt-metrics-collector-${VERSION}-linux-amd64
	hack/container/docker-push.sh ${VERSION}

clean:
	rm -f cmd/kubevirt-metrics-collector/kubevirt-metrics-collector
	rm -rf _out

unittests: binary
	go test -v ./...

.PHONY: all vendor binary release clean unittests

