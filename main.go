package main

import (
	"context"
	"encoding/json"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	configMapWithSource := getConfigMapWithSourceObject(request["name"], request["namespace"], request["deps"], request["source"])
	_, err = clientset.CoreV1().ConfigMaps(request["namespace"]).Create(context.TODO(), configMapWithSource, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}

	configMapWithServer := getConfigMapWithServerObject(request["namespace"])
	_, err = clientset.CoreV1().ConfigMaps(request["namespace"]).Create(context.TODO(), configMapWithServer, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}

	pod := getPodObject(request["name"], request["namespace"])
	_, err = clientset.CoreV1().Pods(request["namespace"]).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Function created successfully...")
}
