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

package spec

import "github.com/upmc-enterprises/elasticsearch-operator/pkg/snapshot"

// ElasticSearchCluster defines the cluster
type ElasticSearchCluster struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   map[string]string `json:"metadata"`
	Spec       ClusterSpec       `json:"spec"`
}

// ClusterSpec defines cluster options
type ClusterSpec struct {
	// ClusterName is the elasticsearch cluster name
	ClusterName string `json:"cluster-name"`

	// NodeSelector specifies a map of key-value pairs. For the pod to be eligible
	// to run on a node, the node must have each of the indicated key-value pairs as
	// labels.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Zones specifies a map of key-value pairs. Defines which zones
	// to deploy persistent volumes for data nodes
	Zones []string `json:"zones,omitempty"`

	// NodeSpecs is a map of each node type's configuration
	// settings.
	NodeSpecs *NodeSpecs `json:"node-specs,omitempty"`

	// DataDiskSize specifies how large the persistent volume should be attached
	// to the data nodes in the ES cluster
	DataDiskSize string `json:"data-volume-size"`

	// ElasticSearchImage specifies the docker image to use (optional)
	ElasticSearchImage string `json:"elastic-search-image"`

	// Snapshot defines how snapshots are scheduled
	Snapshot Snapshot `json:"snapshot"`

	// Storage defines how volumes are provisioned
	Storage Storage `json:"storage"`

	Scheduler *snapshot.Scheduler
}

// NodeSpecs marshals ThirdPartyResource data for specifying each separate node
// type which forms the elasticsearch cluster.
type NodeSpecs struct {
	Master *NodeTypeSettings `json:"master,omitempty"`
	Client *NodeTypeSettings `json:"client,omitempty"`
	Data   *NodeTypeSettings `json:"data,omitempty"`
	Ingest *NodeTypeSettings `json:"ingest,omitempty"`
}

// NodeTypeSettings marshals Elasticsearch settings for each node type.
// The data is used to configure settings for k8s Deployments or StatefulSets.
type NodeTypeSettings struct {
	Replicas int32  `json:"replicas,omitempty"`
	CpuReq   string `json:"cpu-req,omitempty"`
	MemReq   string `json:"mem-req,omitempty"`
	CpuLimit string `json:"cpu-limit,omitempty"`
	MemLimit string `json:"mem-limit,omitempty"`

	HeapMax int `json:"heap-max,omitempty"`
	HeapMin int `json:"heap-min,omitempty"`
}

// Snapshot defines all params to create / store snapshots
type Snapshot struct {
	// Enabled determines if snapshots are enabled
	SchedulerEnabled bool `json:"scheduler-enabled"`

	// BucketName defines the AWS S3 bucket to store snapshots
	BucketName string `json:"bucket-name"`

	// CronSchedule defines how to run the snapshots
	// SEE: https://godoc.org/github.com/robfig/cron
	CronSchedule string `json:"cron-schedule"`
}

// Storage defines how dynamic volumes are created
// https://kubernetes.io/docs/user-guide/persistent-volumes/
type Storage struct {
	// StorageType is the type of storage to create
	StorageType string `json:"type"`

	// StorageClassProvisoner is the storage provisioner type
	StorageClassProvisoner string `json:"storage-class-provisioner"`
}
