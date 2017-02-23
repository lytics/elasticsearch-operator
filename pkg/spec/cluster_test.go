package spec

import (
	"encoding/json"
	"testing"
)

var csRaw = []byte(`{
  "apiVersion": "enterprises.upmc.com/v1",
  "kind": "ElasticsearchCluster",
  "metadata": {
    "name": "es-cluster"
  },
  "spec": {
    "cluster-name": "primary",
    "zones": ["us-central1-a", "us-central1-f"],
    "node-specs":{
        "master":{
            "replicas": 3
        },
        "client": {
            "replicas": 1
        },
        "data": {
            "replicas": 3,
            "cpu-req": "8000m",
            "mem-req": "30Gi",
            "heap-max": 28000,
            "heap-min": 28000
        },
        "ingest": {}
    },
    "data-volume-size": "10Gi",
    "storage": {
      "type": "pd-ssd",
      "storage-class-provisioner": "kubernetes.io/gce-pd"
    }
  }
}`)

func TestUnmarshalTPR(t *testing.T) {
	var esc ElasticSearchCluster
	err := json.Unmarshal(csRaw, &esc)
	if err != nil {
		t.Errorf("error unmarshaling: %v", err)
	}

	if esc.Kind != "ElasticsearchCluster" {
		t.Errorf("bad kind value: %s", esc.Kind)
	}

	if esc.Spec.NodeSpecs.Data.CpuReq != "8000m" {
		t.Errorf("error unmarshalling data.cpu-req: %v", esc.Spec.NodeSpecs.Data.CpuReq)
	}
}
