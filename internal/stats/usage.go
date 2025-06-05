package stats

import (
	"sync"
	"time"
)

type UsageStats struct {
	mu         sync.Mutex
	Downloads  map[string]int
	Publishes  map[string]int
	LastAccess map[string]time.Time
}

var stats = &UsageStats{
	Downloads:  make(map[string]int),
	Publishes:  make(map[string]int),
	LastAccess: make(map[string]time.Time),
}

func IncDownload(pkg string) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.Downloads[pkg]++
	stats.LastAccess[pkg] = time.Now()
}

func IncPublish(pkg string) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.Publishes[pkg]++
	stats.LastAccess[pkg] = time.Now()
}

func GetStats() UsageStats {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	return *stats
}