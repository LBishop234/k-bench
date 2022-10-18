package deploy

type manifestType string

// TODO: List of types to handle, remove when done.
const (
	deployment  manifestType = "deployment"
	daemonSet   manifestType = "daemonSet"
	replicSet   manifestType = "replicaSet"
	statefulSet manifestType = "statefulSet"

	configMap             manifestType = "configMap"
	persistentVolume      manifestType = "persistentVolume"
	persistentVolumeClaim manifestType = "persistentVolumeClaim"
	secret                manifestType = "secret"
	service               manifestType = "service"

	unknown manifestType = "unknown"
)
