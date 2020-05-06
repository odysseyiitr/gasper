package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sdslabs/gasper/lib/utils"

	"github.com/sdslabs/gasper/configs"
	"github.com/sdslabs/gasper/lib/docker"
	"github.com/sdslabs/gasper/types"
)

// CreateRedisDBContainer  creates a RedisDB container
func CreateRedisDBContainer(db types.Database) error {
	storepath, _ := os.Getwd()
	var err error
	port, err := utils.GetFreePort()

	if err != nil {
		return fmt.Errorf("Error while getting free port for container : %s", err)
	}

	storedir := filepath.Join(storepath, "redis-storage", db.GetName())

	if err := os.MkdirAll(storedir, 0755); err != nil {
		return fmt.Errorf("Error while creating the directory : %s", err)
	}

	containerID, err := docker.CreateDatabaseContainer(&types.DatabaseContainer{
		Image:         configs.ImageConfig.Redis,
		ContainerPort: port,
		DatabasePort:  6379,
		Env:           configs.ServiceConfig.Kaen.Redis.Env,
		WorkDir:       "/data/",
		StoreDir:      filepath.Join(storepath, "redis-storage", db.GetName()),
		Name:          db.GetName(),
		Cmd:           []string{"redis-server", "--requirepass", db.GetPassword()},
	})

	if err != nil {
		return types.NewResErr(500, "container not created", err)
	}

	db.SetContainerPort(port)
	if err := docker.StartContainer(containerID); err != nil {
		return types.NewResErr(500, "container not started", err)
	}

	return nil
}

// DeleteRedisDBContainer deletes RedisDB container
func DeleteRedisDBContainer(containerID string) error {
	if err := docker.DeleteContainer(containerID); err != nil {
		return types.NewResErr(500, "container not deleted", err)
	}

	storepath, _ := os.Getwd()
	storedir := filepath.Join(storepath, "redis-storage", containerID)

	if err := os.RemoveAll(storedir); err != nil {
		return fmt.Errorf("Error while deleting the database directory : %s", err)
	}
	return nil
}
