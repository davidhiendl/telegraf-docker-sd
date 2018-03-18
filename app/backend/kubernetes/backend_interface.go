package kubernetes

import (
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"github.com/davidhiendl/telegraf-docker-sd/app/backend"
	"github.com/davidhiendl/telegraf-docker-sd/app/utils"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

const (
	NAME = "kubernetes"
)

type KubernetesBackend struct {
	commonConfig     *config.ConfigSpec
	templates        map[string]*sdtemplate.Template
	telegrafReloader *utils.TelegrafReloader

	config  *KubernetesConfigSpec
	client  *kubernetes.Clientset
	tags    []string
	runMode int
	node    *corev1.Node
}

func NewBackend() *KubernetesBackend {
	return &KubernetesBackend{
	}
}

func (backend *KubernetesBackend) Name() string {
	return NAME
}

func (backend *KubernetesBackend) Status() int {
	return 1
}

func (backend *KubernetesBackend) Init(spec *backend.BackendConfigSpec) {
	backend.commonConfig = spec.Config
	backend.templates = spec.Templates
	backend.telegrafReloader = spec.Reloader
	backend.config = LoadConfig()

	// create client
	client, err := backend.createKubeClient()
	if err != nil {
		logger.Fatalf("failed to create kubernetes client: %+v", err)
	}
	backend.client = client

	// find current node
	node, err := backend.findCurrentKubeNode()
	if err != nil {
		logger.Fatalf("failed to find current node: %v", err)
	}
	backend.node = node
}

func (backend *KubernetesBackend) Run() {

	pods, err := backend.client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	logger.Infof("There are %d pods in the cluster", len(pods.Items))
	backend.processPodsOnCurrentKubeNode()
}

func (backend *KubernetesBackend) Clean() {
	backend.telegrafReloader.ShouldReload = true
}
