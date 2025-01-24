package port

type HealthCheckService interface {
	Ping() string
	Status() (map[string]interface{}, error)
}
