package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"

	// "github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HealthCheck(c *gin.Context) {
	CPUCheck(c)
	RamCheck(c)
	DiskCheck(c)
}

func CPUCheck(c *gin.Context) {
	physical_cores, _ := cpu.Counts(false)
	logical_cores, _ := cpu.Counts(true)

	// load_avg, _ := load.Avg()

	cpu, err := cpu.Percent(0, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "CPU check failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         "CPU check passed",
		"physical_cores": physical_cores,
		"logical_cores":  logical_cores,
		"cpu_usage":      len(cpu),
		// "load_avg":       load_avg,
	})
}

func RamCheck(c *gin.Context) {
	v, _ := mem.VirtualMemory()

	c.JSON(http.StatusOK, gin.H{
		"status":    "Ram check passed",
		"total":     v.Total,
		"Avaliable": v.Available,
		"used":      v.Used,
	})
}

func DiskCheck(c *gin.Context) {
	u, err := disk.Usage("/")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "Disk check failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "Disk check passed",
		"total":  u.Total,
		"free":   u.Free,
		"used":   u.Used,
	})
}
