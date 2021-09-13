package morsica

import "time"

// SignalCallback is the callback called when a signal is started or stopped
type SignalCallback func()

// Interval stores information about an element (dit, dah, space) interval
type Interval struct {
	Signal bool
	Ms     int
}

// IntervalSequencePlayer is able to play a seuqnce of intervals
type IntervalSequencePlayer struct {
	Sequence    []Interval
	OnSignalOn  SignalCallback
	OnSignalOff SignalCallback
	timer       *time.Timer
	curSeq      int
	bRunning    bool
}

// Timing for dits and dash
type Timing struct {
	ditMs      int
	ditsInADah float32
}

// NewTiming creates a pointer to a new Timing
func NewTiming() *Timing {
	t := new(Timing)
	t.ditMs = 50
	t.ditsInADah = 3
	return t
}

// GetDitMs returns the dit time in ms
func (t *Timing) GetDitMs() int {
	return t.ditMs
}

// SetDitMs sets the dit time in ms
func (t *Timing) SetDitMs(ms int) {
	if ms < 0 {
		return
	}
	t.ditMs = ms
}

// GetDitsInADah returns the length of a dah in dits
func (t *Timing) GetDitsInADah() float32 {
	return t.ditsInADah
}

// SetDitsInADah sets the length of a dah in dits
func (t *Timing) SetDitsInADah(f float32) {
	if f < 1 {
		return
	}
	t.ditsInADah = f
}

// GetDahMs returns the in dah time in ms
func (t *Timing) GetDahMs() int {
	return (t.ditMs * int(t.ditsInADah))
}

// GetInterelementMs returns the ms between dits and dahs in the same symbol
func (t *Timing) GetInterelementMs() int {
	return t.ditMs
}

// GetSpaceMs returns the ms for a space
func (t *Timing) GetSpaceMs() int {
	return t.ditMs
}

//GetInterSymbolMs returns the ms between symbols withing a specific word, given an Alphabet
func (t *Timing) GetInterSymbolMs(a *Alphabet) int {
	return (t.GetSpaceMs() * a.GetInterSymbolSpaces())
}

//GetInterWordMs returns the ms between words, given an Alphabet
func (t *Timing) GetInterWordMs(a *Alphabet) int {
	return (t.GetSpaceMs() * a.GetInterWordSpaces())
}

///
/// Intervals
///

func (t *Timing) makeDitInterval() Interval {
	return Interval{
		Signal: true,
		Ms:     t.ditMs,
	}
}

func (t *Timing) makeDahInterval() Interval {
	return Interval{
		Signal: true,
		Ms:     t.GetDahMs(),
	}
}

func (t *Timing) makeInterelementInterval() Interval {
	return Interval{
		Signal: false,
		Ms:     t.GetInterelementMs(),
	}
}

func (t *Timing) makeSpaceInterval(numSpaces int) Interval {
	return Interval{
		Signal: false,
		Ms:     (t.GetSpaceMs() * numSpaces),
	}
}

func countContinousSpaces(msg *string, startIdx int) int {
	var n int = 0
	for i := startIdx; i < len(*msg); i++ {
		if (*msg)[i] == ' ' {
			n++
		} else {
			return n
		}
	}
	return n
}

func doesItNeedInterelementSpace(msg *string, idx int) bool {
	if idx >= len(*msg)-1 {
		return false
	}
	if (*msg)[idx+1] == ' ' {
		return false
	}
	return true
}

// MorseMessageToIntervalSequence returns an Interval sequence for a Morse message
func (t *Timing) MorseMessageToIntervalSequence(msg string) []Interval {
	seq := make([]Interval, 0)
	for i := 0; i < len(msg); i++ {
		if msg[i] == '.' {
			seq = append(seq, t.makeDitInterval())
			if doesItNeedInterelementSpace(&msg, i) {
				seq = append(seq, t.makeInterelementInterval())
			}
		} else if msg[i] == '-' {
			seq = append(seq, t.makeDahInterval())
			if doesItNeedInterelementSpace(&msg, i) {
				seq = append(seq, t.makeInterelementInterval())
			}
		} else {
			if msg[i] == ' ' {
				numSpaces := countContinousSpaces(&msg, i)
				seq = append(seq, t.makeSpaceInterval(numSpaces))
				i = i + numSpaces - 1
			}
		}
	}
	return seq
}

///
/// Interval playback
///

// NewIntervalSequencePlayer returns a new IntervalSequencePlayer pointer
func NewIntervalSequencePlayer(seq []Interval) *IntervalSequencePlayer {
	player := new(IntervalSequencePlayer)
	player.Sequence = seq
	player.bRunning = false
	player.curSeq = 0
	return player
}

// Stop stops a player if it's running
func (p *IntervalSequencePlayer) Stop() {
	if p.timer == nil {
		return
	}
	p.bRunning = false
	p.timer.Stop()
}

// Start starts a player
func (p *IntervalSequencePlayer) Start() {
	if len(p.Sequence) < 1 {
		return
	}
	p.Stop()
	p.curSeq = 0
	p.startTimerForInterval(p.curSeq)
}

func (p *IntervalSequencePlayer) startTimerForInterval(idx int) {
	if idx < 0 || (idx >= len(p.Sequence)-1) {
		p.Stop()
		return
	}
	p.curSeq = idx
	if p.Sequence[p.curSeq].Signal {
		if p.OnSignalOn != nil {
			p.OnSignalOn()
		}
	}
	p.timer = time.NewTimer(time.Millisecond * time.Duration(p.Sequence[p.curSeq].Ms))
	//go func() {
	<-p.timer.C
	if p.Sequence[p.curSeq].Signal {
		if p.OnSignalOff != nil {
			p.OnSignalOff()
		}
	}
	p.startTimerForInterval(p.curSeq + 1)
	//}()
}
