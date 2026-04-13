package health

import "time"

type Status struct {
	Status    string `json:"status"`
	Check     string `json:"check"`
	Timestamp string `json:"timestamp"`
}

func Health(now time.Time) Status {
	return Status{
		Status:    "ok",
		Check:     "liveness",
		Timestamp: now.UTC().Format(time.RFC3339),
	}
}

func Ready(now time.Time) Status {
	return Status{
		Status:    "ok",
		Check:     "readiness",
		Timestamp: now.UTC().Format(time.RFC3339),
	}
}
