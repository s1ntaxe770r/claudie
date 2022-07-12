package testingframework

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Berops/platform/proto/pb"
	"github.com/Berops/platform/services/kuber/server/kubectl"
	"github.com/rs/zerolog/log"
)

const (
	maxLonghornCheck = 240 // max allowed time for pods of longhorn-system to be ready [seconds]
	sleepSecPods     = 10  // seconds for one cycle of longhorn checks (the node and pod checks)
)

type KubectlOutputJSON struct {
	APIVersion string                   `json:"apiVersion"`
	Items      []map[string]interface{} `json:"items"`
	Kind       string                   `json:"kind"`
	Metadata   map[string]interface{}   `json:"metadata"`
}

// testLonghornDeployment function will perform actions needed to confirm that longhorn has been successfully deployed in the cluster
func testLonghornDeployment(config *pb.GetConfigFromDBResponse) error {
	//start longhorn testing
	clusters := config.Config.CurrentState.Clusters
	for _, cluster := range clusters {
		// check number of nodes in nodes.longhorn.io

		kubectl := kubectl.Kubectl{Kubeconfig: cluster.Kubeconfig}
		err := checkLonghornNodes(cluster, kubectl)
		if err != nil {
			return fmt.Errorf("error while checking the nodes.longhorn.io : %v", err)

		}
		// check if all pods from longhorn-system are ready
		err = checkLonghornPods(cluster.Kubeconfig, cluster.ClusterInfo.Name, kubectl)
		if err != nil {
			return fmt.Errorf("error while checking if all pods from longhorn-system are ready : %v", err)
		}
	}
	return nil
}

// checkLonghornNodes will check if the count of nodes.longhorn.io is same as number of schedulable nodes
func checkLonghornNodes(cluster *pb.K8Scluster, kubectl kubectl.Kubectl) error {
	allNodesFound := false
	readyCheck := 0
	workerCount := 0
	count := 0 //in order to save last value for error message, the var is defined here
	//count the worker nodes
	for _, nodepool := range cluster.ClusterInfo.NodePools {
		if !nodepool.IsControl {
			workerCount += int(nodepool.Count)
		}
	}
	// give them time of maxLonghornCheck seconds to be scheduled
	for readyCheck < maxLonghornCheck {
		out, err := kubectl.KubectlGet("nodes.longhorn.io -A -o json", "")
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("error while getting the nodes.longhorn.io in cluster %s : %v", cluster.ClusterInfo.Name, err))
		}
		allNodesFound, count, err = parseNodesOutput(out, workerCount)
		if err != nil {
			return fmt.Errorf(fmt.Sprintf("error while checking the kubectl output for  nodes.longhorn.io in cluster  %s : %v", cluster.ClusterInfo.Name, err))
		}
		// the number of worker nodes should be equal to number of scheduled nodes in longhorn
		// NOTE: by default, master nodes will not be used to schedule pods, however, if this changes the condition will be broken
		if allNodesFound {
			break
		}
		readyCheck += sleepSecPods
		time.Sleep(time.Duration(sleepSecPods) * time.Second)
		log.Info().Msgf("Waiting for nodes.longhorn.io to be initialized in cluster %s... [ %ds elapsed ]", cluster.ClusterInfo.Name, readyCheck)
	}
	if !allNodesFound {
		return fmt.Errorf(fmt.Sprintf("the count of schedulable nodes (%d) is not equal to nodes.longhorn.io (%d) in cluster %s", workerCount, count, cluster.ClusterInfo.Name))
	}
	return nil
}

// checkLonghornPods will check if the pods in longhorn-system namespace are in ready state
func checkLonghornPods(config, clusterName string, kubectl kubectl.Kubectl) error {
	readyCheck := 0
	allPodsReady := false
	// give them time of maxLonghornCheck seconds to be scheduled
	for readyCheck < maxLonghornCheck {
		out, err := kubectl.KubectlGet("pods -o json", "longhorn-system")
		if err != nil {
			return fmt.Errorf("error while getting the status of the pods in longhorn-system in cluster %s : %v", clusterName, err)
		}
		ready, err := parsePodsOutput(out)
		if err != nil {
			log.Error().Msgf("Error while checking longhorn pods : %v", err)
		}
		// if some are not ready, wait sleepSecPods seconds
		if !ready {
			readyCheck += sleepSecPods
		} else {
			allPodsReady = true
			break
		}
		time.Sleep(time.Duration(sleepSecPods) * time.Second)
		log.Info().Msgf("Waiting for pods from longhorn-system namespace in cluster %s to be in ready state... [ %ds elapsed ]", clusterName, readyCheck)
	}
	if !allPodsReady {
		return fmt.Errorf("pods in longhorn-system took too long to initialize in cluster %s", clusterName)
	}
	return nil
}

// function will parse kubectl json output regarding the longhorn nodes
// returns true if every pod is ready, false otherwise
func parseNodesOutput(out []byte, nodesExpected int) (bool, int, error) {
	// parse output
	var parsedJSON KubectlOutputJSON
	err := json.Unmarshal(out, &parsedJSON)
	if err != nil {
		return false, -1, fmt.Errorf("error while unmarshalling output data : %v", err)
	}
	// get number of nodes currently
	nodes := len(parsedJSON.Items)

	if nodes == nodesExpected {
		return true, nodes, nil
	} else {
		return false, nodes, nil
	}
}

// function will parse kubectl json output regarding the longhorn pods
// returns true if every pod is ready, false otherwise
func parsePodsOutput(out []byte) (bool, error) {
	// parse output
	var parsedJSON KubectlOutputJSON
	err := json.Unmarshal(out, &parsedJSON)
	if err != nil {
		return false, fmt.Errorf("error while unmarshalling output data : %v", err)
	}
	// iterate over all returned items
	for _, item := range parsedJSON.Items {
		// get status field
		status := item["status"].(map[string]interface{})
		// get container statuses
		containerStatuses := status["containerStatuses"].([]interface{})
		// check all container statuses if they are ready
		for _, conStat := range containerStatuses {
			ready := conStat.(map[string]interface{})["ready"].(bool)
			// if not ready, return false
			if !ready {
				log.Info().Msgf("Container %s is not ready yet...", conStat.(map[string]interface{})["name"].(string))
				return false, nil
			}
		}
	}
	return true, nil
}