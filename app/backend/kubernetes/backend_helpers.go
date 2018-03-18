package kubernetes

import (
	"strings"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"flag"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func (backend *KubernetesBackend) processConfig() {
	backend.processConfigLabelsAsTags()
}

func (backend *KubernetesBackend) processConfigLabelsAsTags() {
	labelsRaw := strings.Split(backend.config.TagsFromLabels, ",")

	labelsClean := make([]string, len(labelsRaw))
	i := 0
	for _, label := range labelsRaw {
		labelsClean[i] = label
		i++
	}

	backend.tags = labelsClean
}

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
	var kubeconfig *string
	if home := os.Getenv("HOME"); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	rest.InClusterConfig()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	return kubernetes.NewForConfig(config)
}
