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

func getConfigMapWithServerObject(namespace string) *core.ConfigMap {
	return &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "server",
			Namespace: namespace,
		},
		Data: map[string]string{
			"index.js": "    const express = require(\"express\");\n    const main = require(\"./handler\");\n    const app = express();\n    const port = 3000;\n\n    app.get(\"/\", (req, res) => {\n      const ret = main.main(res);\n      console.log(ret);\n      res.send(ret);\n    });\n\n    app.listen(port, () => {\n      console.log(`Example app listening at http://localhost:${port}`);\n    });",
		},
	}
}

func getPodObject(name string, namespace string) *core.Pod {
	i := int32(420)
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"run": name,
			},
		},
		Spec: core.PodSpec{
			Volumes: []core.Volume{
				{
					Name: "source-deps-" + name,
					VolumeSource: core.VolumeSource{
						ConfigMap: &core.ConfigMapVolumeSource{
							LocalObjectReference: core.LocalObjectReference{Name: "source-deps-" + name},
							DefaultMode:          &i,
						},
					},
				},
				{
					Name: "server",
					VolumeSource: core.VolumeSource{
						ConfigMap: &core.ConfigMapVolumeSource{
							LocalObjectReference: core.LocalObjectReference{Name: "server"},
							DefaultMode:          &i,
						},
					},
				},
				{
					Name: "emptydir",
					VolumeSource: core.VolumeSource{
						EmptyDir: &core.EmptyDirVolumeSource{
							Medium: core.StorageMediumDefault,
						},
					},
				},
			},
			DNSPolicy:     core.DNSClusterFirst,
			RestartPolicy: core.RestartPolicyAlways,
		},
	}
}
