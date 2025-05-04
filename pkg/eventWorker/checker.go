package eventworker

import "time"

func checkDeltaTime(startTime time.Duration, delta time.Duration) bool {
	if len(scheduledStartTime) == 0 {
		scheduledStartTime = append(scheduledStartTime, startTime)

		return true
	}

	divTime := startTime - scheduledStartTime[len(scheduledStartTime)-1]

	scheduledStartTime = append(scheduledStartTime, startTime)

	return divTime == delta
}

func checkActualTime(scheduledStartTime, actualStartTime time.Duration) bool {
	return actualStartTime >= scheduledStartTime
}

func checkStartTime(actualTime, startTime time.Duration) bool {
	return actualTime >= startTime
}
