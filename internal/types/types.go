package types

// ContainerData will maintain the id and status of docker containers
type ContainerData struct {
	ID          int
	ContainerID string
	Running     bool
}
