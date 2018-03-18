package kubernetes

import (
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"path/filepath"
)

// TrackedPod is used to maintain state about already processed containers and to be able to remove their configurations easily
type TrackedPod struct {
	backend    *KubernetesBackend
	configFile string

	UID  types.UID
	Name string
	Pod  *corev1.Pod

	Config map[string]string
	Tags   map[string]string
}

// Create new config and populate it from environment
func NewTrackedPod(backend *KubernetesBackend, pod *corev1.Pod) *TrackedPod {
	tp := TrackedPod{
		UID:     pod.GetUID(),
		Name:    pod.Name,
		backend: backend,
		Pod:     pod,

		Config: make(map[string]string),
		Tags:   make(map[string]string),
	}

	// parse config
	tp.parseAnnotationsAsConfig()

	// parse tags, precedence order: system tags > annotations > labels
	tp.parseLabelsAsTags()
	tp.parseAnnotationsAsTags()
	tp.parseSystemTags()

	// debug
	logger.Debugf(LOG_PREFIX + "[%v] config: %+v", tp.Name, tp.Config)
	logger.Debugf(LOG_PREFIX + "[%v] labels: %+v", tp.Name, tp.Pod.Labels)
	logger.Debugf(LOG_PREFIX + "[%v] annotations: %+v", tp.Name, tp.Pod.Annotations)
	logger.Debugf(LOG_PREFIX + "[%v] tags: %+v", tp.Name, tp.Tags)

	return &tp
}

func (tp *TrackedPod) ConfigFile() string {
	if tp.configFile == "" {
		file, _ := filepath.Abs(
			tp.backend.commonConfig.ConfigDir +
				"/" + tp.backend.commonConfig.AutoConfPrefix +
				tp.backend.Name() + "_" +
				tp.Name +
				tp.backend.commonConfig.AutoConfExtension)
		tp.configFile = file
	}
	return tp.configFile
}

func (tp *TrackedPod) TargetIP() string {
	return tp.Pod.Status.PodIP
}
