package cache

import (
	"sync"
	"time"

	"github.com/gopanel/gopanel/internal/collector"
)

var (
	dockerMu        sync.RWMutex
	dockerData      []collector.Container
	dockerUpdatedAt time.Time

	servicesMu        sync.RWMutex
	servicesData      []collector.SystemdService
	servicesUpdatedAt time.Time
)

// Start 启动后台定时刷新
func Start(interval time.Duration) {
	refreshDocker()
	refreshServices()
	go func() {
		t := time.NewTicker(interval)
		for range t.C {
			refreshDocker()
			refreshServices()
		}
	}()
}

func refreshDocker() {
	data, err := collector.GetContainers()
	if err != nil { return }
	dockerMu.Lock()
	dockerData = data
	dockerUpdatedAt = time.Now()
	dockerMu.Unlock()
}

func refreshServices() {
	data, err := collector.GetServices()
	if err != nil { return }
	servicesMu.Lock()
	servicesData = data
	servicesUpdatedAt = time.Now()
	servicesMu.Unlock()
}

func GetDockerContainers() ([]collector.Container, bool) {
	dockerMu.RLock()
	defer dockerMu.RUnlock()
	return dockerData, dockerData != nil
}

func GetServices() ([]collector.SystemdService, bool) {
	servicesMu.RLock()
	defer servicesMu.RUnlock()
	return servicesData, servicesData != nil
}

func InvalidateDocker() { go refreshDocker() }
func InvalidateServices() { go refreshServices() }
