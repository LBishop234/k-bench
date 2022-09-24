package configHandler

type connType string

const (
	CLUSTER  connType = "cluster"
	MINIKUBE connType = "minikube"
)

type connStruct struct {
	Type       connType
	Kubeconfig string
}
