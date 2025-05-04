package eventworker

import (
	"YadroTest/pkg/config"
	"fmt"
	"log"
	"strings"
	"time"
)

// Слайс функций ивентов.
var incomingEvents = []func(map[string]*competotor, *config.Config, string, string, string) string{
	regComp,
	setTime,
	compOnStartLine,
	compIsStart,
	compOnFiringRange,
	targetHit,
	compLeftFiringRange,
	compEnterPenaltyLap,
	compLeftPenaltyLap,
	compEndMainLap,
	compCantContinue,
}

// Слайс расписания стартов. Нужен для проверки startDelta.
var scheduledStartTime []time.Duration

// Структура участника соревнования.
type competotor struct {
	scheduledStartTime time.Duration
	actualStartTime    time.Duration

	scheduledFinishTime time.Duration
	actualFinishTime    time.Duration

	// Состояние участника. Если участник пришел успешно к финишу,
	// то mark замениться на финишное время.
	mark string

	ml            []mainLap
	pl            []penaltyLap
	numberOfShots []int

	// Количество оставшихся кругов.
	laps int

	// Количество оставшихся участков со стрельбой.
	firings int
}

// mainLap структура для информации главного круга.
type mainLap struct {
	timeCompleteLap time.Duration
	averageSpeedLap float64
}

// mainLap структура для информации штрафного круга.
type penaltyLap struct {
	timeCompleteLap time.Duration
	averageSpeedLap float64
}

// regComp Регестрация участника.
func regComp(cm map[string]*competotor, cfg *config.Config, eventTime, compID, _ string) string {
	cm[compID] = &competotor{mark: "NotStarted", laps: cfg.Laps, firings: cfg.FiringLines}

	return fmt.Sprintf("%s The competitor(%s) registered", eventTime, compID)
}

// setTime Время старта по жеребью.
func setTime(cm map[string]*competotor, cfg *config.Config, eventTime, compID, startTime string) string {
	comp := cm[compID]

	startTimeDuration := parseDuration(startTime)

	// Проверка соответствует ли startDelta
	if !checkDeltaTime(startTimeDuration, parseDuration(cfg.StartDelta)) {
		log.Fatal("wrong delta time")
	}

	// Проверка соответствует ли start
	if !checkStartTime(startTimeDuration, parseDuration(cfg.Start)) {
		log.Fatal("wrong start time")
	}

	comp.scheduledStartTime = startTimeDuration

	return fmt.Sprintf(
		"%s The start time for the competitor(%s) was set by a draw to %s",
		eventTime, compID, startTime,
	)
}

// compOnStartLine Участник встал на стартовую линию.
func compOnStartLine(_ map[string]*competotor, _ *config.Config, eventTime, compID, _ string) string {
	return fmt.Sprintf("%s The competitor(%s) is on the start line", eventTime, compID)
}

// compIsStart Участник начал соревнование.
func compIsStart(cm map[string]*competotor, _ *config.Config, eventTime, compID, _ string) string {
	comp := cm[compID]

	actualStartTimeDuration := parseDuration(eventTime)

	comp.actualStartTime = actualStartTimeDuration
	comp.ml = append(comp.ml, mainLap{comp.scheduledStartTime, 0})

	// Проверка начал ли участник в положенное время.
	if !checkActualTime(comp.scheduledStartTime, actualStartTimeDuration) {
		comp.laps = -1

		return fmt.Sprintf("%s The competitor(%s) is disqualified", eventTime, compID)
	}

	comp.mark = "NotFinished"

	return fmt.Sprintf("%s The competitor(%s) has started", eventTime, compID)
}

// compOnFiringRange Участник в зоне стрельбы.
func compOnFiringRange(cm map[string]*competotor, _ *config.Config, eventTime, compID, firingRange string) string {
	comp := cm[compID]
	comp.numberOfShots = append(comp.numberOfShots, 0)

	return fmt.Sprintf(
		"%s The competitor(%s) is on the firing range(%s)",
		eventTime, compID, firingRange,
	)
}

// targetHit Участник попал в мишень.
func targetHit(cm map[string]*competotor, _ *config.Config, eventTime, compID, targetID string) string {
	comp := cm[compID]
	comp.numberOfShots[len(comp.numberOfShots)-1]++

	return fmt.Sprintf(
		"%s The target(%s) has been hit by competitor(%s)",
		eventTime, targetID, compID,
	)
}

// compLeftFiringRange Участник вышел из зоны стрельбы.
func compLeftFiringRange(cm map[string]*competotor, _ *config.Config, eventTime, compID, _ string) string {
	comp := cm[compID]
	comp.firings--

	return fmt.Sprintf("%s The competitor(%s) left the firing range", eventTime, compID)
}

// compEnterPenaltyLap Участник в штрафной зоне.
func compEnterPenaltyLap(cm map[string]*competotor, _ *config.Config, eventTime, compID, _ string) string {
	comp := cm[compID]

	comp.pl = append(comp.pl, penaltyLap{parseDuration(eventTime), 0})

	return fmt.Sprintf("%s The competitor(%s) entered the penalty laps", eventTime, compID)
}

// compLeftPenaltyLap Участник вышел из штрафной зоны.
func compLeftPenaltyLap(cm map[string]*competotor, cfg *config.Config, eventTime, compID, _ string) string {
	comp := cm[compID]

	eventTimeDuration := parseDuration(eventTime)

	// Считаем среднюю скорость
	comp.pl[len(comp.pl)-1].averageSpeedLap =
		float64(cfg.PenaltyLen) / (eventTimeDuration.Seconds() - comp.pl[len(comp.pl)-1].timeCompleteLap.Seconds())

	// Считаем время в штрафной зоне
	comp.pl[len(comp.pl)-1].timeCompleteLap =
		eventTimeDuration - comp.pl[len(comp.pl)-1].timeCompleteLap

	return fmt.Sprintf("%s The competitor(%s) left the penalty laps", eventTime, compID)
}

// compEndMainLap Участник закончил круг.
func compEndMainLap(cm map[string]*competotor, cfg *config.Config, eventTime, compID, _ string) string {
	comp := cm[compID]

	eventTimeDuration := parseDuration(eventTime)

	// Считаем среднюю скорость
	comp.ml[len(comp.ml)-1].averageSpeedLap =
		float64(cfg.LapLen) / (eventTimeDuration.Seconds() - comp.ml[len(comp.ml)-1].timeCompleteLap.Seconds())

	// Считаем время круга
	comp.ml[len(comp.ml)-1].timeCompleteLap =
		eventTimeDuration - comp.ml[len(comp.ml)-1].timeCompleteLap

	comp.laps--

	// Если это не последний круг, то подготавливаем данные
	if comp.laps != 0 {
		comp.ml = append(comp.ml, mainLap{eventTimeDuration, 0})
	}

	// Если это последний круг, то пишем результ соревнования
	if comp.laps == 0 {
		comp.scheduledFinishTime = eventTimeDuration - comp.scheduledStartTime
		comp.actualFinishTime = eventTimeDuration - comp.actualStartTime
	}

	return fmt.Sprintf("%s The competitor(%s) ended the main lap", eventTime, compID)
}

// compCantContinue Участник не может продолжить соревнование.
func compCantContinue(cm map[string]*competotor, _ *config.Config, eventTime, compID, comment string) string {
	comp := cm[compID]
	comp.laps = -1

	return fmt.Sprintf(
		"%s The competitor(%s) can`t continue: %s",
		eventTime, compID, comment,
	)
}

// parseDuration перевод времени ивента из строки в time.Duration.
func parseDuration(str string) time.Duration {
	str = strings.Trim(str, "[]")

	parts := strings.Split(str, ":")

	formatted := fmt.Sprintf("%sh%sm%ss", parts[0], parts[1], parts[2])

	result, err := time.ParseDuration(formatted)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
