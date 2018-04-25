package kubernetes

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"path/filepath"
	"github.com/sirupsen/logrus"
)

// TrackedPod is used to maintain state about already processed containers and to be able to remove their configurations easily
type TrackedPod struct {
	backend    *KubernetesBackend
	configFile string

	UID   types.UID
	Name  string
	Pod   *corev1.Pod
	PodIP string

	Config map[string]string
	Tags   map[string]string
	Env    map[string]string
}

// Create new config and populate it from environment
func NewTrackedPod(backend *KubernetesBackend, pod *corev1.Pod) *TrackedPod {
	tp := TrackedPod{
		backend: backend,
		Env:     backend.commonConfig.EnvMap,
	}

	tp.importPodData(pod)

	return &tp
}

func (tp *TrackedPod) importPodData(pod *corev1.Pod) {
	// extract basic info
	tp.Pod = pod
	tp.UID = pod.GetUID()
	tp.Name = pod.Name
	tp.PodIP = pod.Status.PodIP

	// init maps
	tp.Config = make(map[string]string)
	tp.Tags = make(map[string]string)

	// parse config
	tp.parseAnnotationsAsConfig()

	// parse tags, precedence order: system tags > annotations > labels
	tp.parseLabelsAsTags()
	tp.parseAnnotationsAsTags()
	tp.parseSystemTags()

	// debug
	logrus.Debugf(LOG_PREFIX+"[%v] config: %+v", tp.Name, tp.Config)
	logrus.Debugf(LOG_PREFIX+"[%v] labels: %+v", tp.Name, tp.Pod.Labels)
	logrus.Debugf(LOG_PREFIX+"[%v] annotations: %+v", tp.Name, tp.Pod.Annotations)
	logrus.Debugf(LOG_PREFIX+"[%v] tags: %+v", tp.Name, tp.Tags)
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
	return tp.PodIP
}
