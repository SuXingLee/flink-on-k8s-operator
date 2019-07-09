/*
Copyright 2019 Google LLC.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	flinkoperatorv1alpha1 "github.com/googlecloudplatform/flink-operator/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Converter which converts the FlinkSessionCluster spec to the desired
// underlying Kubernetes resource specs.

// Gets the desired JobManager deployment spec from the FlinkSessionCluster spec.
func getDesiredJobManagerDeployment(flinkSessionCluster *flinkoperatorv1alpha1.FlinkSessionCluster) appsv1.Deployment {
	var clusterNamespace = flinkSessionCluster.ObjectMeta.Namespace
	var clusterName = flinkSessionCluster.ObjectMeta.Name
	var imageSpec = flinkSessionCluster.Spec.ImageSpec
	var jobManagerSpec = flinkSessionCluster.Spec.JobManagerSpec
	var rpcPort = corev1.ContainerPort{Name: "rpc", ContainerPort: *jobManagerSpec.Ports.RPC}
	var blobPort = corev1.ContainerPort{Name: "blob", ContainerPort: *jobManagerSpec.Ports.Blob}
	var queryPort = corev1.ContainerPort{Name: "query", ContainerPort: *jobManagerSpec.Ports.Query}
	var uiPort = corev1.ContainerPort{Name: "ui", ContainerPort: *jobManagerSpec.Ports.UI}
	var jobManagerDeploymentName = clusterName + "-jobmanager"
	var labels = map[string]string{
		"cluster":   clusterName,
		"app":       "flink",
		"component": "jobmanager",
	}
	var jobManagerDeployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       clusterNamespace,
			Name:            jobManagerDeploymentName,
			OwnerReferences: []metav1.OwnerReference{toOwnerReference(flinkSessionCluster)},
			Labels:          labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: jobManagerSpec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "jobmanager",
							Image: *imageSpec.URI,
							Args:  []string{"jobmanager"},
							Ports: []corev1.ContainerPort{rpcPort, blobPort, queryPort, uiPort},
							Env:   []corev1.EnvVar{corev1.EnvVar{Name: "JOB_MANAGER_RPC_ADDRESS", Value: jobManagerDeploymentName}},
						},
					},
				},
			},
		},
	}
	return jobManagerDeployment
}

// Gets the desired JobManager service spec from the FlinkSessionCluster spec.
func getDesiredJobManagerService(flinkSessionCluster *flinkoperatorv1alpha1.FlinkSessionCluster) corev1.Service {
	var clusterNamespace = flinkSessionCluster.ObjectMeta.Namespace
	var clusterName = flinkSessionCluster.ObjectMeta.Name
	var jobManagerSpec = flinkSessionCluster.Spec.JobManagerSpec
	var rpcPort = corev1.ServicePort{
		Name:       "rpc",
		Port:       *jobManagerSpec.Ports.RPC,
		TargetPort: intstr.FromString("rpc")}
	var blobPort = corev1.ServicePort{
		Name:       "blob",
		Port:       *jobManagerSpec.Ports.Blob,
		TargetPort: intstr.FromString("blob")}
	var queryPort = corev1.ServicePort{
		Name:       "query",
		Port:       *jobManagerSpec.Ports.Query,
		TargetPort: intstr.FromString("query")}
	var uiPort = corev1.ServicePort{
		Name:       "ui",
		Port:       *jobManagerSpec.Ports.UI,
		TargetPort: intstr.FromString("ui")}
	var jobManagerServiceName = clusterName + "-jobmanager"
	var labels = map[string]string{
		"cluster":   clusterName,
		"app":       "flink",
		"component": "jobmanager",
	}
	var jobManagerService = corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       clusterNamespace,
			Name:            jobManagerServiceName,
			OwnerReferences: []metav1.OwnerReference{toOwnerReference(flinkSessionCluster)},
			Labels:          labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    []corev1.ServicePort{rpcPort, blobPort, queryPort, uiPort},
		},
	}
	return jobManagerService
}

// Gets the desired TaskManager deployment spec from the FlinkSessionCluster spec.
func getDesiredTaskManagerDeployment(flinkSessionCluster *flinkoperatorv1alpha1.FlinkSessionCluster) appsv1.Deployment {
	var clusterNamespace = flinkSessionCluster.ObjectMeta.Namespace
	var clusterName = flinkSessionCluster.ObjectMeta.Name
	var imageSpec = flinkSessionCluster.Spec.ImageSpec
	var taskManagerSpec = flinkSessionCluster.Spec.TaskManagerSpec
	var dataPort = corev1.ContainerPort{Name: "data", ContainerPort: *taskManagerSpec.Ports.Data}
	var rpcPort = corev1.ContainerPort{Name: "rpc", ContainerPort: *taskManagerSpec.Ports.RPC}
	var queryPort = corev1.ContainerPort{Name: "query", ContainerPort: *taskManagerSpec.Ports.Query}
	var taskManagerDeploymentName = clusterName + "-taskmanager"
	var jobManagerDeploymentName = clusterName + "-jobmanager"
	var labels = map[string]string{
		"cluster":   clusterName,
		"app":       "flink",
		"component": "taskmanager",
	}
	var taskManagerDeployment = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:       clusterNamespace,
			Name:            taskManagerDeploymentName,
			OwnerReferences: []metav1.OwnerReference{toOwnerReference(flinkSessionCluster)},
			Labels:          labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: taskManagerSpec.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						corev1.Container{
							Name:  "taskmanager",
							Image: *imageSpec.URI,
							Args:  []string{"taskmanager"},
							Ports: []corev1.ContainerPort{dataPort, rpcPort, queryPort},
							Env:   []corev1.EnvVar{corev1.EnvVar{Name: "JOB_MANAGER_RPC_ADDRESS", Value: jobManagerDeploymentName}},
						},
					},
				},
			},
		},
	}
	return taskManagerDeployment
}

// Converts the FlinkSessionCluster as owner reference for its child resources.
func toOwnerReference(flinkSessionCluster *flinkoperatorv1alpha1.FlinkSessionCluster) metav1.OwnerReference {
	return metav1.OwnerReference{
		APIVersion:         flinkSessionCluster.APIVersion,
		Kind:               flinkSessionCluster.Kind,
		Name:               flinkSessionCluster.Name,
		UID:                flinkSessionCluster.UID,
		Controller:         &[]bool{true}[0],
		BlockOwnerDeletion: &[]bool{false}[0],
	}
}