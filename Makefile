VERSIONDIR := internal/pkg/version
VERSIONFILE := $(VERSIONDIR)/version.go

all: binary

docker: binary
	docker build .

dockertag: binary
	./hack/dockertag.sh

dockerpush: binary
	./hack/dockerpush.sh

vendor:
	dep ensure

binary: vendor gensrc
	cd cmd/kubevirt-metrics-collector && go build -v .

clean:
	rm -f cmd/kubevirt-metrics-collector/kubevirt-metrics-collector

gensrc:
	@mkdir -p $(VERSIONDIR) && ./hack/genver.sh > $(VERSIONFILE)

.PHONY: all docker dockertag dockerpush binary clean

