package main

import (
	"guzey-timer/pkg/phases"
	"os"
	"os/exec"

	"time"

	tm "github.com/buger/goterm"
	"github.com/fatih/color"

	"github.com/common-nighthawk/go-figure"
)

// Guzey Timer
// Timer should alternate between WORK and REST phases
// as defined here (https://guzey.com/productivity/):
//  - work for 25 minutes from :05 to :30
//  - take a 5 minute break from :30 to :35
//  - work for 25 minutes from :35 to :00
//  - take a 5 minute break from :00 to :05
//	- every three hours (at 12-3-6-9) the :05-:30 work cycle is substituted for a break, which lasts 35 minutes.

var phaseColor string = phases.WORK_COLOR
var phaseColors map[phases.Phase]string = map[phases.Phase]string{phases.WORK: "green", phases.BREAK: "red"}
var opposite map[phases.Phase]phases.Phase = map[phases.Phase]phases.Phase{phases.WORK: phases.BREAK, phases.BREAK: phases.WORK}

func main() {
	var end time.Time
	var currentPhase phases.Phase
	var nextPhase phases.Phase
	// Track Time for PhaseTimers
	go func() {
		for {
			pt := getPhaseTimer(time.Now())
			if pt == nil {
				continue
			}
			// change display info
			phaseColor = phaseColors[pt.Phase()]
			end = pt.EndingTime()
			currentPhase = pt.Phase()
			nextPhase = pt.State.NextPhase
			// announce start of phase
			go func() {
				cmd := exec.Command("afplay", pt.MP3()) // clears the scrollback buffer as well as the screen.
				cmd.Stdout = os.Stdout
				cmd.Run()
			}()
			<-pt.Timer().C
		}
	}()

	time.Sleep(time.Millisecond * 500)
	// Display Time
	for {
		DisplayClock(currentPhase, nextPhase, end, time.Millisecond*100, phaseColor)
	}

}

func getPhaseTimer(current_time time.Time) *phases.PhaseTimer {

	hour := current_time.Hour()
	minute := current_time.Minute()
	second := current_time.Second()
	switch hour % 3 {
	case 0: // 6a,9a,12p,3p,6p,...
		// minute starting a break phase
		if minute == 0 {
			return phases.NewLongBreak()
		}
		// minute starting a work phase
		if minute == 35 {
			return phases.NewWorkPhase()
		}
		// minutes during a break phase
		if minute > 0 && minute < 35 {
			minDiff := time.Minute * time.Duration(35-minute)
			secDiff := time.Second * time.Duration(second)
			duration := (minDiff - secDiff)
			return phases.NewGenericBreak(duration)
		}
		// minutes during a work phase
		if minute > 35 && minute < 60 {
			minDiff := time.Minute * time.Duration(60-minute)
			secDiff := time.Second * time.Duration(second)
			duration := (minDiff - secDiff)
			return phases.NewGenericWork(duration)

		}
	default:
		// minutes starting a break phase
		if minute == 0 || minute == 30 {
			return phases.NewShortBreak()
		}
		if minute > 0 && minute < 5 {
			minDiff := time.Minute * time.Duration(5-minute)
			secDiff := time.Second * time.Duration(second)
			duration := (minDiff - secDiff)
			return phases.NewGenericBreak(duration)
		}
		// minutes during a working phase
		if minute > 5 && minute < 30 {
			minDiff := time.Minute * time.Duration(30-minute)
			secDiff := time.Second * time.Duration(second)
			duration := (minDiff - secDiff)
			return phases.NewGenericWork(duration)
		}
		// minutes starting a work phase
		if minute == 5 || minute == 35 {
			return phases.NewWorkPhase()
		}
		// minutes during a break phase
		if minute > 30 && minute < 35 {
			minDiff := time.Minute * time.Duration(35-minute)
			secDiff := time.Second * time.Duration(second)
			duration := (minDiff - secDiff)
			return phases.NewGenericBreak(duration)
		}
		// minutes during a working phase
		if minute > 35 && minute < 60 {
			minDiff := time.Minute * time.Duration(60-minute)
			secDiff := time.Second * time.Duration(second)
			duration := (minDiff - secDiff)
			return phases.NewGenericWork(duration)
		}

	}
	return nil
}

func DisplayClock(
	currentPhase phases.Phase,
	nextPhase phases.Phase,
	end time.Time,
	delay time.Duration,
	phaseColor string,
) {
	//Display Time
	tm.MoveCursor(1, 1)
	tm.Println("Current phase: ", currentPhase, " / Next phase starting @ ", end.Format(time.Kitchen), " - ", nextPhase)
	myFigure := figure.NewColorFigure(time.Now().Format(time.Kitchen), "", phaseColor, true)
	var colorFunc func(a ...interface{}) string
	switch phaseColor {
	case phases.WORK_COLOR:
		colorFunc = color.New(color.FgGreen).SprintFunc()
	case phases.BREAK_COLOR:
		colorFunc = color.New(color.FgRed).SprintFunc()
	}
	tm.Println(colorFunc(myFigure.String()))
	time.Sleep(delay)
	clearScreen()
	tm.Flush()
}

var clearScreen = func() {
	cmd := exec.Command(`printf '\33c\e[3J'`) // clears the scrollback buffer as well as the screen.
	cmd.Stdout = os.Stdout
	cmd.Run()
}
