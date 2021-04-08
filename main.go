package main

import (
	"encoding/json"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"net/http"
	"path/filepath"
)

var clientset *kubernetes.Clientset

func main() {
	var err error
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/function", createFunction)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func createFunction(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var request map[string]string
	err := decoder.Decode(&request)
	if err != nil {
		panic(err)
	}

	fmt.Println("Function created successfully...")
}
