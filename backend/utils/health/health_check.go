package health

import (
	"context"
	"time"

	"github.com/imankit007/url-shortner/infrastructure"
)

// HealthStatus represents the health status of a service
type HealthStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

// CheckRedisHealth checks if Redis is healthy
func CheckRedisHealth() HealthStatus {
	if infrastructure.RedisClient == nil {
		return HealthStatus{
			Status:  "unhealthy",
			Message: "Redis client not initialized",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := infrastructure.RedisClient.Ping(ctx).Result()
	if err != nil {
		return HealthStatus{
			Status:  "unhealthy",
			Message: "Redis connection failed: " + err.Error(),
		}
	}

	return HealthStatus{
		Status: "healthy",
	}
}

// CheckMongoHealth checks if MongoDB is healthy
func CheckMongoHealth() HealthStatus {
	if infrastructure.MongoClient == nil {
		return HealthStatus{
			Status:  "unhealthy",
			Message: "MongoDB client not initialized",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := infrastructure.MongoClient.Ping(ctx, nil)
	if err != nil {
		return HealthStatus{
			Status:  "unhealthy",
			Message: "MongoDB connection failed: " + err.Error(),
		}
	}

	return HealthStatus{
		Status: "healthy",
	}
}

// OverallHealth represents the overall health of all services
type OverallHealth struct {
	Status   string                  `json:"status"`
	Services map[string]HealthStatus `json:"services"`
}

// CheckOverallHealth checks the health of all services
func CheckOverallHealth() OverallHealth {
	redisHealth := CheckRedisHealth()
	mongoHealth := CheckMongoHealth()

	services := map[string]HealthStatus{
		"redis": redisHealth,
		"mongo": mongoHealth,
	}

	status := "healthy"
	for _, service := range services {
		if service.Status != "healthy" {
			status = "unhealthy"
			break
		}
	}

	return OverallHealth{
		Status:   status,
		Services: services,
	}
}
