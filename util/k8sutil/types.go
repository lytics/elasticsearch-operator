/*
Copyright (c) 2017, UPMC Enterprises
All rights reserved.
Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
    * Redistributions of source code must retain the above copyright
      notice, this list of conditions and the following disclaimer.
    * Redistributions in binary form must reproduce the above copyright
      notice, this list of conditions and the following disclaimer in the
      documentation and/or other materials provided with the distribution.
    * Neither the name UPMC Enterprises nor the
      names of its contributors may be used to endorse or promote products
      derived from this software without specific prior written permission.
THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL UPMC ENTERPRISES BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
*/

package k8sutil

import (
	"fmt"

	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func CreateDataNodeStatefulSetSpec(baseImage, clusterName, statefulSetName, storageClass string, mlow, mmax int, volumeSize resource.Quantity) *apps.StatefulSetSpec {
	replicas := int32(1)
	memOpts := fmt.Sprintf("-Xms%dm -Xmx%dm", mlow, mmax)
	spec := &apps.StatefulSetSpec{
		Replicas:    &replicas,
		ServiceName: "es-data-svc",
		Template: v1.PodTemplateSpec{
			ObjectMeta: v1.ObjectMeta{
				Labels: map[string]string{
					"component": "elasticsearch",
					"role":      "data",
					"name":      statefulSetName,
				},
				Annotations: map[string]string{
					"pod.beta.kubernetes.io/init-containers": "[ { \"name\": \"sysctl\", \"image\": \"busybox\", \"imagePullPolicy\": \"IfNotPresent\", \"command\": [\"sysctl\", \"-w\", \"vm.max_map_count=262144\"], \"securityContext\": { \"privileged\": true } }]",
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					v1.Container{
						Name: statefulSetName,
						SecurityContext: &v1.SecurityContext{
							Privileged: &[]bool{true}[0],
							Capabilities: &v1.Capabilities{
								Add: []v1.Capability{
									"IPC_LOCK",
								},
							},
						},
						Image:           baseImage,
						ImagePullPolicy: "Always",
						Env: []v1.EnvVar{
							v1.EnvVar{
								Name: "NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							v1.EnvVar{
								Name:  "CLUSTER_NAME",
								Value: clusterName,
							},
							v1.EnvVar{
								Name:  "NODE_MASTER",
								Value: "false",
							},
							v1.EnvVar{
								Name:  "HTTP_ENABLE",
								Value: "false",
							},
							v1.EnvVar{
								Name:  "ES_JAVA_OPTS",
								Value: memOpts,
							},
						},
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								Name:          "transport",
								ContainerPort: 9300,
								Protocol:      v1.ProtocolTCP,
							},
						},
						VolumeMounts: []v1.VolumeMount{
							v1.VolumeMount{
								Name:      "es-data",
								MountPath: "/data",
							},
							v1.VolumeMount{
								Name:      "es-certs",
								MountPath: "/elasticsearch/config/certs",
							},
						},
					},
				},
				Volumes: []v1.Volume{
					v1.Volume{
						Name: "es-certs",
						VolumeSource: v1.VolumeSource{
							Secret: &v1.SecretVolumeSource{
								SecretName: "es-certs",
							},
						},
					},
				},
			},
		},
		VolumeClaimTemplates: []v1.PersistentVolumeClaim{
			v1.PersistentVolumeClaim{
				ObjectMeta: v1.ObjectMeta{
					Name: "es-data",
					Annotations: map[string]string{
						"volume.beta.kubernetes.io/storage-class": storageClass,
					},
				},
				Spec: v1.PersistentVolumeClaimSpec{
					AccessModes: []v1.PersistentVolumeAccessMode{
						v1.ReadWriteOnce,
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceStorage: volumeSize,
						},
					},
				},
			},
		},
	}
	return spec
}

func CreateNodeDeployment(baseImage, deploymentName, role, isNodeMaster, httpEnable, clusterName string, replicas *int32) *v1beta1.DeploymentSpec {
	return &v1beta1.DeploymentSpec{
		Replicas: replicas,
		Template: v1.PodTemplateSpec{
			ObjectMeta: v1.ObjectMeta{
				Labels: map[string]string{
					"component": "elasticsearch",
					"role":      role,
					"name":      deploymentName,
				},
				Annotations: map[string]string{
					"pod.beta.kubernetes.io/init-containers": "[ { \"name\": \"sysctl\", \"image\": \"busybox\", \"imagePullPolicy\": \"IfNotPresent\", \"command\": [\"sysctl\", \"-w\", \"vm.max_map_count=262144\"], \"securityContext\": { \"privileged\": true } }]",
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					v1.Container{
						Name: deploymentName,
						SecurityContext: &v1.SecurityContext{
							Privileged: &[]bool{true}[0],
							Capabilities: &v1.Capabilities{
								Add: []v1.Capability{
									"IPC_LOCK",
								},
							},
						},
						Image:           baseImage,
						ImagePullPolicy: "Always",
						Env: []v1.EnvVar{
							v1.EnvVar{
								Name: "NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							v1.EnvVar{
								Name:  "CLUSTER_NAME",
								Value: clusterName,
							},
							v1.EnvVar{
								Name:  "NODE_MASTER",
								Value: isNodeMaster,
							},
							v1.EnvVar{
								Name:  "NODE_DATA",
								Value: "false",
							},
							v1.EnvVar{
								Name:  "HTTP_ENABLE",
								Value: httpEnable,
							},
							v1.EnvVar{
								Name:  "ES_JAVA_OPTS",
								Value: "-Xms1024m -Xmx1024m",
							},
						},
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								Name:          "transport",
								ContainerPort: 9300,
								Protocol:      v1.ProtocolTCP,
							},
							v1.ContainerPort{
								Name:          "http",
								ContainerPort: 9200,
								Protocol:      v1.ProtocolTCP,
							},
						},
						VolumeMounts: []v1.VolumeMount{
							v1.VolumeMount{
								Name:      "storage",
								MountPath: "/data",
							},
							v1.VolumeMount{
								Name:      "es-certs",
								MountPath: "/elasticsearch/config/certs",
							},
						},
					},
				},
				Volumes: []v1.Volume{
					v1.Volume{
						Name: "storage",
						VolumeSource: v1.VolumeSource{
							EmptyDir: &v1.EmptyDirVolumeSource{},
						},
					},
					v1.Volume{
						Name: "es-certs",
						VolumeSource: v1.VolumeSource{
							Secret: &v1.SecretVolumeSource{
								SecretName: "es-certs",
							},
						},
					},
				},
			},
		},
	}
}
