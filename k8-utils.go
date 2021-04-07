package main

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

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
			InitContainers: []core.Container{
				{
					Name:  "init",
					Image: "busybox",
					Command: []string{
						"sh",
						"-c",
						"cp /tmp/server/index.js /sources",
					},
					VolumeMounts: []core.VolumeMount{
						{
							Name:      "emptydir",
							MountPath: "/sources",
						},
						{
							Name:      "source-deps-" + name,
							MountPath: "/tmp/src/handler.js",
							SubPath:   "source",
						},
						{
							Name:      "source-deps-" + name,
							MountPath: "/tmp/src/package.json",
							SubPath:   "dependencies",
						},
						{
							Name:      "server",
							MountPath: "/tmp/server/index.js",
							SubPath:   "index.js",
						},
					},
				},
				{
					Name:  "init2",
					Image: "busybox",
					Command: []string{
						"sh",
						"-c",
						"cp /tmp/src/* /sources",
					},
					VolumeMounts: []core.VolumeMount{
						{
							Name:      "emptydir",
							MountPath: "/sources",
						},
						{
							Name:      "source-deps-" + name,
							MountPath: "/tmp/src/handler.js",
							SubPath:   "source",
						},
						{
							Name:      "source-deps-" + name,
							MountPath: "/tmp/src/package.json",
							SubPath:   "dependencies",
						},
						{
							Name:      "server",
							MountPath: "/tmp/server/index.js",
							SubPath:   "index.js",
						},
					},
				},
				{
					Name:  "init3",
					Image: "node:alpine",
					Command: []string{
						"sh",
						"-c",
						"cd /sources && npm install",
					},
					VolumeMounts: []core.VolumeMount{
						{
							Name:      "emptydir",
							MountPath: "/sources",
						},
					},
				},
			},
			Containers: []core.Container{
				{
					Name:  name,
					Image: "node:alpine",
					Command: []string{
						"node",
						"/sources/index.js",
					},
					VolumeMounts: []core.VolumeMount{
						{
							Name:      "emptydir",
							MountPath: "/sources",
						},
					},
				},
			},
			DNSPolicy:     core.DNSClusterFirst,
			RestartPolicy: core.RestartPolicyAlways,
		},
	}
}

func getConfigMapWithSourceObject(t map[string]string) *core.ConfigMap {
	return &core.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "source-deps-" + t["name"],
			Namespace: t["namespace"],
		},
		Data: map[string]string{
			"dependencies": t["deps"],
			"source":       t["source"],
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

func getServiceObject(name string, namespace string) *core.Service {
	return &core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels: map[string]string{
				"run": name,
			},
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Port:     3000,
					Protocol: "TCP",
					TargetPort: intstr.IntOrString{
						IntVal: 3000,
					},
				},
			},
			Selector: map[string]string{
				"run": name,
			},
		},
	}
}