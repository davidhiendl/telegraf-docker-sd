package kubernetes

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"errors"
	"strconv"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
)

func (backend *KubernetesBackend) GetRunMode() int {
	return backend.runMode
}

func (backend *KubernetesBackend) createKubeClient() (*kubernetes.Clientset, error) {

	// attempt to connect inside the cluster if possible
	config, errIn := rest.InClusterConfig()
	if errIn == nil {
		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		// set run mode
		backend.runMode = RUN_MODE_INTERNAL
		return client, nil
	}

	// attempt to connect using external client
	client, err := backend.getKubeExternalConfig()
	if err != nil {
		return nil, err
	}

	// set run mode
	backend.runMode = RUN_MODE_EXTERNAL
	return client, nil
}

func (backend *KubernetesBackend) getKubeExternalConfig() (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	rest.InClusterConfig()
	config, err := clientcmd.BuildConfigFromFlags("", backend.config.KubeConfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	return kubernetes.NewForConfig(config)
}

// find current hostname
func (backend *KubernetesBackend) getAgentHostname() (string, error) {
	// check if override is defined
	if len(backend.config.NodeNameOverride) > 0 {
		return backend.config.NodeNameOverride, nil
	} else {
		// fallback to container/hostname
		agentHostname, err := os.Hostname()
		if err != nil {
			return "", err
		}
		return agentHostname, nil
	}
}

// Find a node by the pod hostname or current object hostname.
func (backend *KubernetesBackend) findCurrentKubeNode() (*corev1.Node, error) {

	// find current hostname
	agentHostname, err := backend.getAgentHostname()
	if err != nil {
		return nil, err
	}

	if backend.runMode == RUN_MODE_INTERNAL {

		pods, err := backend.client.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}

		// find node name by pod name, this requires that the pod hostname is unchanged
		var nodeName string;
		for _, pod := range pods.Items {
			if pod.Name == agentHostname {
				nodeName = pod.Spec.NodeName
				break;
			}
		}
		if nodeName == "" {
			return nil, errors.New("failed to find current pod, runMode: cluster-internal, agent pod hostname: " + agentHostname)
		}

		// find node
		nodes, err := backend.client.CoreV1().Nodes().List(metav1.ListOptions{})
		for _, node := range nodes.Items {
			if node.Name == nodeName {
				return &node, nil
			}
		}

		return nil, errors.New("failed to find current node, runMode: cluster-internal, node name derived from pod: " + nodeName)

	} else if backend.runMode == RUN_MODE_EXTERNAL {

		// find node matching hostname
		nodes, err := backend.client.CoreV1().Nodes().List(metav1.ListOptions{})
		if err != nil {
			errors.New("failed to find current node, runMode: cluster-external, agent hostname: " + agentHostname)
		}

		// find node
		for _, node := range nodes.Items {
			if node.Name == agentHostname {
				return &node, nil
			}
		}

		return nil, errors.New("failed to find current node, runMode: cluster-external, agent hostname: " + agentHostname)

	} else {
		return nil, errors.New("failed to determine current run mode, mode=" + strconv.Itoa(backend.runMode))
	}
}
