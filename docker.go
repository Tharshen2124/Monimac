package main

import (
	"bufio"
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

// Container holds information about a single Docker container.
type Container struct {
	ID     string
	Name   string
	Image  string
	Status string
	CPU    string
	Mem    string
}

// dockerPSEntry matches the JSON output of `docker ps --format '{{json .}}'`.
type dockerPSEntry struct {
	ID     string `json:"ID"`
	Names  string `json:"Names"`
	Image  string `json:"Image"`
	Status string `json:"Status"`
	State  string `json:"State"`
}

// dockerStatsEntry matches the JSON output of `docker stats --no-stream --format '{{json .}}'`.
type dockerStatsEntry struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	CPUPerc  string `json:"CPUPerc"`
	MemUsage string `json:"MemUsage"`
}

// ListContainers returns all running Docker containers with stats merged in.
// Returns (nil, error) if Docker is unavailable.
func ListContainers() ([]Container, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	psOut, err := exec.CommandContext(ctx, "docker", "ps", "--format", "{{json .}}").Output()
	if err != nil {
		return nil, err
	}

	psEntries := make(map[string]dockerPSEntry)
	scanner := bufio.NewScanner(strings.NewReader(string(psOut)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry dockerPSEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue
		}
		psEntries[entry.ID] = entry
	}

	if len(psEntries) == 0 {
		return []Container{}, nil
	}

	statsCtx, statsCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer statsCancel()

	statsOut, err := exec.CommandContext(statsCtx, "docker", "stats", "--no-stream", "--format", "{{json .}}").Output()
	statsMap := make(map[string]dockerStatsEntry)
	if err == nil {
		statsScanner := bufio.NewScanner(strings.NewReader(string(statsOut)))
		for statsScanner.Scan() {
			line := strings.TrimSpace(statsScanner.Text())
			if line == "" {
				continue
			}
			var entry dockerStatsEntry
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				continue
			}
			statsMap[entry.ID] = entry
		}
	}

	containers := make([]Container, 0, len(psEntries))
	for id, ps := range psEntries {
		name := strings.TrimPrefix(ps.Names, "/")
		c := Container{
			ID:     id,
			Name:   name,
			Image:  ps.Image,
			Status: ps.Status,
			CPU:    "--",
			Mem:    "--",
		}
		if stats, ok := statsMap[id]; ok {
			c.CPU = stats.CPUPerc
			c.Mem = stats.MemUsage
		}
		containers = append(containers, c)
	}

	return containers, nil
}

// StopContainer stops the Docker container with the given ID.
func StopContainer(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "docker", "stop", id).Run()
}
