package store

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gopanel/gopanel/internal/collector"
	"github.com/gopanel/gopanel/internal/config"
	"github.com/gopanel/gopanel/internal/websocket"
)

func Init(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_journal=WAL&_timeout=5000")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp INTEGER NOT NULL,
			cpu_percent REAL,
			mem_percent REAL,
			disk_percent REAL,
			net_recv INTEGER,
			net_sent INTEGER
		);
		CREATE INDEX IF NOT EXISTS idx_metrics_ts ON metrics(timestamp);
		CREATE TABLE IF NOT EXISTS alerts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp INTEGER NOT NULL,
			type TEXT NOT NULL,
			value REAL,
			threshold REAL,
			message TEXT
		);
	`)
	return db, err
}

func SaveMetrics(db *sql.DB, snap collector.MetricsSnapshot) {
	var maxDisk float64
	for _, p := range snap.Disk.Partitions {
		if p.UsedPercent > maxDisk {
			maxDisk = p.UsedPercent
		}
	}
	var totalRecv, totalSent uint64
	for _, iface := range snap.Network.Interfaces {
		totalRecv += iface.BytesRecv
		totalSent += iface.BytesSent
	}
	db.Exec(`INSERT INTO metrics (timestamp,cpu_percent,mem_percent,disk_percent,net_recv,net_sent) VALUES (?,?,?,?,?,?)`,
		snap.Timestamp, snap.CPU.UsagePercent, snap.Memory.UsedPercent, maxDisk, totalRecv, totalSent)

	// Prune old data (keep 7 days)
	cutoff := time.Now().Add(-7 * 24 * time.Hour).Unix()
	db.Exec(`DELETE FROM metrics WHERE timestamp < ?`, cutoff)
}

func GetMetricsHistory(db *sql.DB, hours int) ([]map[string]interface{}, error) {
	since := time.Now().Add(-time.Duration(hours) * time.Hour).Unix()
	rows, err := db.Query(`SELECT timestamp,cpu_percent,mem_percent,disk_percent FROM metrics WHERE timestamp>? ORDER BY timestamp ASC`, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var ts int64
		var cpu, mem, disk float64
		rows.Scan(&ts, &cpu, &mem, &disk)
		result = append(result, map[string]interface{}{
			"timestamp": ts, "cpu": cpu, "memory": mem, "disk": disk,
		})
	}
	return result, nil
}

func GetAlerts(db *sql.DB, limit int) ([]map[string]interface{}, error) {
	rows, err := db.Query(`SELECT id,timestamp,type,value,threshold,message FROM alerts ORDER BY timestamp DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []map[string]interface{}
	for rows.Next() {
		var id int64
		var ts int64
		var atype, msg string
		var value, threshold float64
		rows.Scan(&id, &ts, &atype, &value, &threshold, &msg)
		result = append(result, map[string]interface{}{
			"id": id, "timestamp": ts, "type": atype, "value": value, "threshold": threshold, "message": msg,
		})
	}
	return result, nil
}

// cooldown prevents repeated alerts (1 per 10 mins per type)
var alertCooldown = make(map[string]time.Time)

func checkAlert(db *sql.DB, cfg *config.Config, alertType string, value, threshold float64) {
	if threshold <= 0 || value < threshold {
		return
	}
	// Cooldown check
	if t, ok := alertCooldown[alertType]; ok && time.Since(t) < 10*time.Minute {
		return
	}
	alertCooldown[alertType] = time.Now()

	msg := fmt.Sprintf("%s 使用率 %.1f%% 超过阈值 %.0f%%", alertType, value, threshold)
	db.Exec(`INSERT INTO alerts (timestamp,type,value,threshold,message) VALUES (?,?,?,?,?)`,
		time.Now().Unix(), alertType, value, threshold, msg)

	if cfg.Alert.Webhook != "" {
		go sendWebhook(cfg.Alert.Webhook, alertType, value, threshold, msg)
	}
}

func sendWebhook(webhookURL, alertType string, value, threshold float64, msg string) {
	payload := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]string{"content": fmt.Sprintf("⚠️ GoPanel 告警\n%s", msg)},
	}
	b, _ := json.Marshal(payload)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewReader(b))
	if err != nil {
		log.Printf("webhook error: %v", err)
		return
	}
	resp.Body.Close()
}

func StartCollector(db *sql.DB, hub *websocket.Hub, interval time.Duration) {
	// Ensure minimum 2s interval to keep resource usage low
	if interval < 2*time.Second {
		interval = 2 * time.Second
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		snap := collector.CollectAll()
		SaveMetrics(db, snap)
		hub.Broadcast("metrics", snap)
	}
}

func StartAlertChecker(db *sql.DB, cfg *config.Config) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		cpu := collector.GetCPUStats()
		mem := collector.GetMemoryStats()
		disk := collector.GetDiskStats()
		checkAlert(db, cfg, "CPU", cpu.UsagePercent, cfg.Alert.CPU)
		checkAlert(db, cfg, "内存", mem.UsedPercent, cfg.Alert.Memory)
		for _, p := range disk.Partitions {
			checkAlert(db, cfg, "磁盘("+p.Mountpoint+")", p.UsedPercent, cfg.Alert.Disk)
		}
	}
}
