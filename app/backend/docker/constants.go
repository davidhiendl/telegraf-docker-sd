package docker

const (
	NAME       = "docker"
	LOG_PREFIX = "[docker]"
)

var SWARM_LABELS = []string{
	"com.docker.stack.namespace",
	"com.docker.swarm.node.id",
	"com.docker.swarm.service.id",
	"com.docker.swarm.service.name",
	"com.docker.swarm.task",
	"com.docker.swarm.task.id",
	"com.docker.swarm.task.name",
}
