package k8s_info

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
)

var (
	deploymentData      []map[string]interface{}
	deploymentDataMutex sync.RWMutex
)

// UpdateDeploymentData 获取k8s deployment资源生成json数据
log.printf("获取k8s deloyment信息")
func UpdateDeploymentData(namespaces ...string) {
	cmd := exec.Command("kubectl", "--kubeconfig=/root/.kube/config.prd", "--output", "json", "get", "deployment", "-A")
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing kubectl command: %s\n", err)
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		log.Printf("Error unmarshalling JSON data: %s\n", err)
		return
	}

	items := data["items"].([]interface{})
	var newDeploymentData []map[string]interface{}

	for _, item := range items {
		app := item.(map[string]interface{})
		metadata := app["metadata"].(map[string]interface{})
		name := metadata["name"].(string)
		namespace := metadata["namespace"].(string)

		for _, ns := range namespaces {
			if strings.Contains(namespace, ns) {
				newDeploymentData = append(newDeploymentData, map[string]interface{}{
					"name":  fmt.Sprintf("%s/%s", namespace, name),
					"value": fmt.Sprintf("export APP_NAME=%s;export NAME_SPACE=%s;", name, namespace),
					"ns":    namespace,
				})
			}
		}
	}

	deploymentDataMutex.Lock()
	deploymentData = newDeploymentData
	deploymentDataMutex.Unlock()
}
