package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func formatUptime(totalSeconds int64) string {
	days := totalSeconds / 86400
	hours := (totalSeconds % 86400) / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if days > 0 {
		return fmt.Sprintf("%dj %dh %dmin", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dmin %ds", hours, minutes, seconds)
	}
	return fmt.Sprintf("%dmin %ds", minutes, seconds)
}

func getVPSUptime() string {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "erreur"
	}
	fields := strings.Fields(string(data))
	totalSeconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return "erreur"
	}
	return formatUptime(int64(totalSeconds))
}

func getMinecraftUptime() string {
	cmd := exec.Command("systemctl", "show", "minecraft", "--property=ActiveEnterTimestamp")
	cmd.Env = append(os.Environ(), "LC_ALL=C", "LANG=C")
	out, err := cmd.Output()
	if err != nil {
		return "offline"
	}
	line := strings.TrimSpace(string(out))
	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 || parts[1] == "" {
		return "offline"
	}
	t, err := time.Parse("Mon 2006-01-02 15:04:05 MST", parts[1])
	if err != nil {
		return "offline"
	}
	return formatUptime(int64(time.Since(t).Seconds()))
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.POST("/start", func(c *gin.Context) {
		if err := exec.Command("systemctl", "start", "minecraft").Run(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "started"})
	})

	r.POST("/stop", func(c *gin.Context) {
		if err := exec.Command("systemctl", "stop", "minecraft").Run(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "stopped"})
	})

	r.GET("/status", func(c *gin.Context) {
		out, _ := exec.Command("systemctl", "is-active", "minecraft").Output()
		c.JSON(200, gin.H{"active": strings.TrimSpace(string(out))})
	})

	r.GET("/uptime", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"vps_uptime":       getVPSUptime(),
			"minecraft_uptime": getMinecraftUptime(),
		})
	})

	r.GET("/logs", func(c *gin.Context) {
		cmd := exec.Command("journalctl", "-u", "minecraft", "-n", "100", "--no-pager")
		out, err := cmd.Output()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		lines := []string{}
		scanner := bufio.NewScanner(strings.NewReader(string(out)))
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		c.JSON(200, lines)
	})

	r.GET("/logs/stream", func(c *gin.Context) {
		cmd := exec.Command("journalctl", "-u", "minecraft", "-f", "--no-pager")
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		if err := cmd.Start(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		defer cmd.Process.Kill()

		scanner := bufio.NewScanner(stdout)
		c.Stream(func(w io.Writer) bool {
			if scanner.Scan() {
				c.SSEvent("log", scanner.Text())
				return true
			}
			return false
		})
	})

	socketPath := "/tmp/mc-agent.sock"
	os.Remove(socketPath)

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatalf("❌ Socket: %v", err)
	}
	os.Chmod(socketPath, 0666)

	log.Println("✅ mc-agent démarré sur", socketPath)
	http.Serve(listener, r)
}
