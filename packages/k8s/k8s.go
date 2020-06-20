/*
This package is a simple facade on top of k8s client library to
modularize kcad.go
*/
package k8s

import (
	"bytes"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
)

var config *rest.Config
var clientSet *kubernetes.Clientset

func Initialize() {
	// Authenticate
	config = InClusterAuth()
	// set global client set
	clientSet = InitClientSet()
}

func InClusterAuth() *rest.Config {
	// Authenticate with cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal("Couldn't load cluster configs", err)
	}
	return config
}

func InitClientSet() *kubernetes.Clientset {
	cSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Couldn't create client set to communicate with API", err)
	}
	return cSet
}

func QueryCadvisor() ([]string, []error) {
	cadvisorList := CAdvisorDiscovery()
	var dataHolder = make([]string, 0, 10)
	var errorHolder = make([]error, 0, 10)
	for _, url := range cadvisorList {
		data, err := clientSet.RESTClient().Get().AbsPath(url).Stream()
		if err != nil {
			errorHolder = append(errorHolder, err)
			continue
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(data)
		metricsStr := buf.String()
		dataHolder = append(dataHolder, metricsStr)
	}
	return dataHolder, errorHolder
}

func CAdvisorDiscovery() []string {
	cadvisorList := make([]string, 0, 10)
	nodes, err := clientSet.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		fmt.Println("Couldn't list nodes in k8s", err)
		return nil
	}

	for _, n := range nodes.Items {
		cadvisorList = append(cadvisorList, fmt.Sprintf("api/v1/nodes/%s/proxy/metrics/cadvisor", n.Name))
	}
	return cadvisorList
}
