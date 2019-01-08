package k8sutils

import (
	"k8s.io/client-go/rest"

	"io/ioutil"
	"log"
	"os"
)

type TLSInfo struct {
	CertFilePath string
	KeyFilePath  string
}

func (ti *TLSInfo) IsEnabled() bool {
	return ti.CertFilePath != "" && ti.KeyFilePath != ""
}

func (ti *TLSInfo) UpdateFromK8S() error {
	if _, err := rest.InClusterConfig(); err != nil {
		// is not a real error, rather a supported case. So, let's swallow the error
		log.Printf("running outside a K8S cluster")
		return nil
	}
	certsDirectory, err := ioutil.TempDir("", "certsdir")
	if err != nil {
		return err
	}
	defer os.RemoveAll(certsDirectory)
	namespace, err := GetNamespace()
	if err != nil {
		log.Printf("Error searching for namespace: %v", err)
		return err
	}
	certStore, err := GenerateSelfSignedCert(certsDirectory, "kubevirt-metrics-collector", namespace)
	if err != nil {
		log.Printf("unable to generate certificates: %v", err)
		return err
	}
	ti.CertFilePath = certStore.CurrentPath()
	ti.KeyFilePath = certStore.CurrentPath()
	log.Printf("running in a K8S cluster: with configuration %#v", *ti)
	return nil
}
