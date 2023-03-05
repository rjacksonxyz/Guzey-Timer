package phases

import (
	"fmt"
	"time"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)

type Phase string

const (
	WORK_COLOR           = "green"
	BREAK_COLOR          = "red"
	WORK_PHASE_DURATION  = time.Minute * 25
	LONG_BREAK_DURATION  = time.Minute * 35
	SHORT_BREAK_DURATION = time.Minute * 5
	MESSAGE_TEMPLATE     = "starting %s phase"
	WORK_MP3_FILENAME    = "work_announcement"
	BREAK_MP3_FILENAME   = "break_announcement"
)

const (
	WORK  Phase = "WORK"
	BREAK Phase = "BREAK"
)

var (
	WORK_MP3_PATH  = ""
	BREAK_MP3_PATH = ""
)

type PhaseMetadata struct {
	PhaseType Phase
	Color     string
	NextPhase Phase
	Mp3       string
}

type PhaseTimer struct {
	timer *time.Timer
	end   time.Time
	State PhaseMetadata
}

var SPEECH htgotts.Speech

func init() {
	SPEECH = htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}
	w, _ := SPEECH.CreateSpeechFile(fmt.Sprintf(MESSAGE_TEMPLATE, WORK), WORK_MP3_FILENAME)
	WORK_MP3_PATH = w
	b, _ := SPEECH.CreateSpeechFile(fmt.Sprintf(MESSAGE_TEMPLATE, BREAK), BREAK_MP3_FILENAME)
	BREAK_MP3_PATH = b
}

func NewPhaseTimer(
	t time.Duration,
	phase Phase,
	color, filepath string,
	nextPhase Phase,
) *PhaseTimer {
	return &PhaseTimer{
		timer: time.NewTimer(t),
		end:   time.Now().Add(t),
		State: PhaseMetadata{
			PhaseType: phase,
			Color:     color,
			Mp3:       filepath,
			NextPhase: nextPhase,
		},
	}
}

func (s *PhaseTimer) Reset(t time.Duration) {
	s.timer.Reset(t)
	s.end = time.Now().Add(t)
}

func (s *PhaseTimer) Stop() {
	s.timer.Stop()
}

func (s *PhaseTimer) TimeRemaining() time.Duration {
	return time.Until(s.end)
}

func (s *PhaseTimer) Phase() Phase {
	return s.State.PhaseType
}

func (s *PhaseTimer) Timer() *time.Timer {
	return s.timer
}

func (s *PhaseTimer) EndingTime() time.Time {
	return s.end
}

func (s *PhaseTimer) MP3() string {
	return s.State.Mp3
}

func NewLongBreak() *PhaseTimer {
	return NewPhaseTimer(LONG_BREAK_DURATION, BREAK, BREAK_COLOR, BREAK_MP3_PATH, WORK)
}

func NewShortBreak() *PhaseTimer {
	return NewPhaseTimer(SHORT_BREAK_DURATION, BREAK, BREAK_COLOR, BREAK_MP3_PATH, WORK)
}

func NewWorkPhase() *PhaseTimer {
	return NewPhaseTimer(WORK_PHASE_DURATION, WORK, WORK_COLOR, WORK_MP3_PATH, BREAK)
}

func NewGenericBreak(duration time.Duration) *PhaseTimer {
	return NewPhaseTimer(duration, BREAK, BREAK_COLOR, BREAK_MP3_PATH, WORK)
}

func NewGenericWork(duration time.Duration) *PhaseTimer {
	return NewPhaseTimer(duration, WORK, WORK_COLOR, WORK_MP3_PATH, BREAK)
}
