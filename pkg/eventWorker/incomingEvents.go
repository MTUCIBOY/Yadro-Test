package eventworker

import (
	"fmt"
)

var incomingEvents = map[string]func(string, string, string) string{
	"1":  regComp,
	"2":  setTime,
	"3":  compOnStartLine,
	"4":  compIsStart,
	"5":  compOnFiringRange,
	"6":  targetHit,
	"7":  compLeftFiringRange,
	"8":  compEnterPenaltyLap,
	"9":  compLeftPenaltyLap,
	"10": compEndMainLap,
	"11": compCantContinue,
}

func regComp(eventTime string, copmID string, _ string) string {
	return fmt.Sprintf("%s The competitor(%s) registered", eventTime, copmID)
}

func setTime(eventTime string, compID string, startTime string) string {
	return fmt.Sprintf(
		"%s The start time for the competitor(%s) was set by a draw to %s",
		eventTime, compID, startTime,
	)
}

func compOnStartLine(eventTime string, copmID string, _ string) string {
	return fmt.Sprintf("%s The competitor(%s) is on the start line", eventTime, copmID)
}

func compIsStart(eventTime string, copmID string, _ string) string {
	return fmt.Sprintf("%s The competitor(%s) has started", eventTime, copmID)
}

func compOnFiringRange(eventTime string, compID string, firingRange string) string {
	return fmt.Sprintf(
		"%s The competitor(%s) is on the firing range(%s)",
		eventTime, compID, firingRange,
	)
}

func targetHit(eventTime string, compID string, targetID string) string {
	return fmt.Sprintf(
		"%s The target(%s) has been hit by competitor(%s)",
		eventTime, targetID, compID,
	)
}

func compLeftFiringRange(eventTime string, copmID string, _ string) string {
	return fmt.Sprintf("%s The competitor(%s) left the firing range", eventTime, copmID)
}

func compEnterPenaltyLap(eventTime string, copmID string, _ string) string {
	return fmt.Sprintf("%s The competitor(%s) entered the penalty laps", eventTime, copmID)
}

func compLeftPenaltyLap(eventTime string, copmID string, _ string) string {
	return fmt.Sprintf("%s The competitor(%s) left the penalty laps", eventTime, copmID)
}

func compEndMainLap(eventTime string, copmID string, _ string) string {
	return fmt.Sprintf("%s The competitor(%s) ended the main lap", eventTime, copmID)
}

func compCantContinue(eventTime string, compID string, comment string) string {
	return fmt.Sprintf(
		"%s The competitor(%s) can`t continue: %s",
		eventTime, compID, comment,
	)
}
