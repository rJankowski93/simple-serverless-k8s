package main

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getConfigMapWithSourceObject(name string, namespace string, deps string, source string) *core.ConfigMap {
	return &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "source-deps-" + name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"dependencies": deps,
			"source":       source,
		},
	}
}
