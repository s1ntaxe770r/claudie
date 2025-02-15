package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Berops/claudie/internal/utils"
	"github.com/Berops/claudie/proto/pb"
)

var cluster = &pb.K8Scluster{
	ClusterInfo: &pb.ClusterInfo{
		Name:       "TestName",
		PublicKey:  "public-key",
		PrivateKey: "private-key",
		NodePools: []*pb.NodePool{
			{
				Name: "Node-1",
				Nodes: []*pb.Node{
					{
						Name:     "server1",
						Public:   "2.2.2.2",
						Private:  "192.168.2.2",
						NodeType: pb.NodeType_master,
					},
					{
						Name:     "server2",
						Public:   "1.1.1.1",
						Private:  "192.168.2.1",
						NodeType: pb.NodeType_master,
					},
				},
			},
			{
				Name: "Node-2",
				Nodes: []*pb.Node{
					{
						Name:     "server3",
						Public:   "3.3.3.3",
						Private:  "192.168.2.3",
						NodeType: pb.NodeType_worker,
					},
					{
						Name:     "server4",
						Public:   "4.4.4.4",
						Private:  "192.168.2.4",
						NodeType: pb.NodeType_worker,
					},
				},
			},
		},
	},
	Kubernetes: "v1.19.0",
	Network:    "192.168.2.0/24",
}

func Test_createKeyFile(t *testing.T) {
	privateKeyFile := "private.pem"
	keyErr := utils.CreateKeyFile(cluster.ClusterInfo.GetPrivateKey(), ".", privateKeyFile)
	if keyErr != nil {
		t.Error("Error writing out .pem file doesn't exist")
	}

	if _, err := os.Stat(filepath.Join(".", privateKeyFile)); os.IsNotExist(err) {
		// path/to/whatever does not exist
		t.Errorf("%s file doesn't exist", privateKeyFile)
	}
}
