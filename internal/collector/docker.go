package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Container struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Image   string  `json:"image"`
	Status  string  `json:"status"`
	State   string  `json:"state"`
	Ports   string  `json:"ports"`
	Created string  `json:"created"`
	CPU     float64 `json:"cpu_percent"`
	MemPct  float64 `json:"mem_percent"`
	MemUsed uint64  `json:"mem_used"`
	MemLim  uint64  `json:"mem_limit"`
}

func GetContainers() ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	out, err := exec.CommandContext(ctx, "docker", "ps", "-a",
		"--format", `{"id":"{{.ID}}","name":"{{.Names}}","image":"{{.Image}}","status":"{{.Status}}","state":"{{.State}}","ports":"{{.Ports}}","created":"{{.CreatedAt}}"}`).Output()
	if err != nil {
		return nil, fmt.Errorf("docker not available: %w", err)
	}

	var containers []Container
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		var c Container
		if err := json.Unmarshal([]byte(line), &c); err == nil {
			containers = append(containers, c)
		}
	}

	// Enrich with stats (non-blocking)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel2()
	statsOut, err := exec.CommandContext(ctx2, "docker", "stats", "--no-stream",
		"--format", `{{.ID}}\t{{.CPUPerc}}\t{{.MemPerc}}\t{{.MemUsage}}`).Output()
	if err == nil {
		type stat struct{ cpu, memPct float64; memUsed, memLim uint64 }
		statsMap := make(map[string]stat)
		for _, line := range strings.Split(strings.TrimSpace(string(statsOut)), "\n") {
			parts := strings.Split(line, "\t")
			if len(parts) < 4 {
				continue
			}
			id := parts[0]
			cpu, _ := strconv.ParseFloat(strings.TrimSuffix(parts[1], "%"), 64)
			memPct, _ := strconv.ParseFloat(strings.TrimSuffix(parts[2], "%"), 64)
			// parse mem usage like "128MiB / 2GiB"
			var used, lim uint64
			memParts := strings.Split(parts[3], " / ")
			if len(memParts) == 2 {
				used = parseMemStr(memParts[0])
				lim = parseMemStr(memParts[1])
			}
			statsMap[id] = stat{cpu, memPct, used, lim}
		}
		for i, c := range containers {
			if s, ok := statsMap[c.ID]; ok {
				containers[i].CPU = s.cpu
				containers[i].MemPct = s.memPct
				containers[i].MemUsed = s.memUsed
				containers[i].MemLim = s.memLim
			}
		}
	}
	return containers, nil
}

func parseMemStr(s string) uint64 {
	s = strings.TrimSpace(s)
	multipliers := map[string]uint64{
		"B": 1, "KiB": 1024, "MiB": 1024 * 1024, "GiB": 1024 * 1024 * 1024,
		"KB": 1000, "MB": 1000 * 1000, "GB": 1000 * 1000 * 1000,
	}
	for suffix, mult := range multipliers {
		if strings.HasSuffix(s, suffix) {
			val, _ := strconv.ParseFloat(strings.TrimSuffix(s, suffix), 64)
			return uint64(val * float64(mult))
		}
	}
	val, _ := strconv.ParseUint(s, 10, 64)
	return val
}

func ContainerAction(id, action string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "docker", action, id).Run()
}

func GetContainerLogs(id string, lines int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, "docker", "logs", "--tail", strconv.Itoa(lines), id).CombinedOutput()
	return string(out), err
}
