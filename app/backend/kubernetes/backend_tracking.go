package kubernetes

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"github.com/davidhiendl/telegraf-docker-sd/app/utils"
	"bytes"
	"github.com/sirupsen/logrus"
)

func (backend *KubernetesBackend) processPodsOnCurrentKubeNode() error {

	pods, err := backend.client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	logrus.Debugf(LOG_PREFIX+" there are %d pods in the cluster", len(pods.Items))

	uidToPod := make(map[types.UID]*corev1.Pod)

	// process trackedPods
	for key, _ := range pods.Items {

		// create explicit reference to value as value from range would be a copy that cannot be referenced
		pod := &pods.Items[key]

		// skip trackedPods on other nodes
		if pod.Spec.NodeName != backend.node.Name {
			continue
		}

		// map trackedPods by UID for cleanup later on
		uidToPod[pod.UID] = pod

		// check if already tracked and if it should be tracked
		_, alreadyTracked := backend.trackedPods[pod.UID]
		shouldTrack := backend.shouldTrackPod(pod)
		if !alreadyTracked && shouldTrack {
			backend.startTrackingPod(pod)
		}
	}

	logrus.Debugf(LOG_PREFIX+" there are %d pods on current node", len(uidToPod))

	// check on already tracked pods
	for uid, trackedPod := range backend.trackedPods {
		pod, ok := uidToPod[uid]

		// check for trackedPods that are still tracked but shouldn't be
		shouldTrack := backend.shouldTrackPod(pod)
		if !ok || !shouldTrack {
			logrus.Infof(LOG_PREFIX+" stopping to track pod: found_in_cluster=%v, shouldTrack=%v", ok, shouldTrack)
			backend.cleanupTrackedPod(trackedPod)

			// do not execute further checks if pod was cleaned up
			continue
		}

		// check if IP changed and generate the configuration accordingly
		if ok && pod.Status.PodIP != trackedPod.PodIP {
			// renew tracking if IP changed
			logrus.Infof(LOG_PREFIX + " tracked pod ip changed, re-configuring")
			trackedPod.importPodData(pod)

			// process templates
			backend.executeTemplatesAgainstTrackedPod(trackedPod)

			// request reload
			backend.telegrafReloader.RequestReload()
		}
	}

	return nil
}

func (backend *KubernetesBackend) cleanupTrackedPod(tp *TrackedPod) {
	logrus.Infof(LOG_PREFIX+"[%v] cleaning up no longer tracked pod", tp.Name)
	utils.RemoveConfigFile(tp.ConfigFile())
	delete(backend.trackedPods, tp.UID)
}

func (backend *KubernetesBackend) shouldTrackPod(pod *corev1.Pod) bool {
	if pod.Status.Phase != corev1.PodRunning {
		return false
	}
	// TODO check if pod is using HostPorts/HostNetworking and if it is reachable (agent on host => yes, agent in cluster: probably no?)
	return true
}

// initialize pod tracking
func (backend *KubernetesBackend) startTrackingPod(pod *corev1.Pod) {
	logrus.Infof(LOG_PREFIX+"[%v] starting to track pod", pod.Name)
	podIP := pod.Status.PodIP

	logrus.Debugf(LOG_PREFIX+"[%v] detected IP: %v", pod.Name, podIP)

	// create object
	trackedPod := NewTrackedPod(backend, pod)

	// process templates
	backend.executeTemplatesAgainstTrackedPod(trackedPod)

	// add to map
	backend.trackedPods[trackedPod.UID] = trackedPod

	// request reload
	backend.telegrafReloader.RequestReload()
}

// process templates for pod and write config file
func (backend *KubernetesBackend) executeTemplatesAgainstTrackedPod(trackedPod *TrackedPod) {
	// process template(s) for container
	configBuffer := new(bytes.Buffer)
	for _, template := range backend.templates {
		logrus.Debugf(LOG_PREFIX+"[%v] running against template: %v", trackedPod.Name, template.FileName)
		err := template.Execute(configBuffer, trackedPod)
		if err != nil {
			logrus.Fatalf(LOG_PREFIX+"[%v] error during template execution: %+v", trackedPod.Name, err)
		}
	}

	// write config
	utils.WriteConfigFile(trackedPod.ConfigFile(), configBuffer.String())
}
