package k8s_info

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	DeploymentData      []map[string]interface{} // 将变量名首字母大写以导出
	DeploymentDataMutex sync.RWMutex
)

// UpdateDeploymentData 获取k8s deployment资源生成json数据
func UpdateDeploymentData(namespaces ...string) {
	log.Println("执行kubectl命令")
	cmd := exec.Command("kubectl", "--kubeconfig=/root/.kube/config", "--output", "json", "get", "deployment", "-A")
	stdout, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing kubectl command: %s\n", err)
		return
	}

	log.Println("JSON解码")
	var data map[string]interface{}
	if err := json.Unmarshal(stdout, &data); err != nil {
		log.Printf("Error unmarshalling JSON data: %s\n", err)
		return
	}

	log.Println("正在生成新的Deployment数据")
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

	log.Println("Deployment数据更新完成")
	DeploymentDataMutex.Lock()
	DeploymentData = newDeploymentData
	DeploymentDataMutex.Unlock()
}

// SetLogger 设置k8s_info包的日志输出
func SetLogger(logger *lumberjack.Logger) {
	log.SetOutput(logger)
}
