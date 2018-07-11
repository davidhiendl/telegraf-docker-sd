package kubernetes

import (
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"github.com/davidhiendl/telegraf-docker-sd/app/backend"
	"github.com/davidhiendl/telegraf-docker-sd/app/utils"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
)

type KubernetesBackend struct {
	commonConfig     *config.ConfigSpec
	templates        map[string]*sdtemplate.Template
	telegrafReloader *utils.TelegrafReloader

	config      *KubernetesConfigSpec
	client      *kubernetes.Clientset
	tags        []string
	runMode     int
	node        *corev1.Node
	trackedPods map[types.UID]*TrackedPod
}

func NewBackend() *KubernetesBackend {
	return &KubernetesBackend{
		trackedPods: make(map[types.UID]*TrackedPod),
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

	// print config
	m := structs.Map(backend.config)
	for key, value := range m {
		logrus.WithFields(logrus.Fields{"key": key, "value": value}).Infof(LOG_PREFIX + " configuration loaded")
	}

	backend.configureClient()
}

func (backend *KubernetesBackend) Run() {
	if backend.node == nil {
		logrus.Infof(LOG_PREFIX + " kubernetes node is not set, attempting to find current node")
		backend.configureKubeNode()

		// do not attempt to process containers if node was not found
		if backend.node == nil {
			return
		}
	}

	backend.processPodsOnCurrentKubeNode()
}

func (backend *KubernetesBackend) Clean() {
	for _, pod := range backend.trackedPods {
		backend.cleanupTrackedPod(pod)
	}
	backend.telegrafReloader.RequestReload()
}
