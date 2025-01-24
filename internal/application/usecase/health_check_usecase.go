package usecase

import (
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// HealthCheckResult é o modelo da resposta do health check
type HealthCheckResult struct {
	Status    string  `json:"status"`
	CPUUsage  float64 `json:"cpu_usage"`
	MemUsage  float64 `json:"mem_usage"`
	Timestamp string  `json:"timestamp"`
}

// HealthCheckUseCase é responsável por realizar o health check
type HealthCheckUseCase struct{}

// NewHealthCheckUseCase cria uma nova instância do use case
func NewHealthCheckUsecase() *HealthCheckUseCase {
	return &HealthCheckUseCase{}
}

// Check retorna informações de saúde do sistema
func (h *HealthCheckUseCase) Check() (*HealthCheckResult, error) {
	// Obter o uso da CPU
	cpuUsage, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	// Obter informações de memória
	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	// Retornar o status de saúde
	result := &HealthCheckResult{
		Status:    "OK",
		CPUUsage:  cpuUsage[0],
		MemUsage:  vm.UsedPercent,
		Timestamp: time.Now().Format(time.RFC3339),
	}
	return result, nil
}
