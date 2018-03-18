package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (backend *KubernetesBackend) processPodsOnCurrentKubeNode() error {

	pods, err := backend.client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	// process pods
	for _, pod := range pods.Items {

		// skip pods on other nodes
		if pod.Spec.NodeName != backend.node.Name {
			continue
		}

		logger.Debugf("processing pod on current node: %v", pod.Name)
	}

	return nil
}
