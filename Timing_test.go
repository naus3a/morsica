package morsica

import (
	"fmt"
	"testing"
	"time"
)

func TestIntervalSequenceCreation(t *testing.T) {
	want := make([]Interval, 0)
	//S
	want = append(want, Interval{
		Signal: true,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: false,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: true,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: false,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: true,
		Ms:     50,
	})
	// [space]
	want = append(want, Interval{
		Signal: false,
		Ms:     150,
	})
	// O
	want = append(want, Interval{
		Signal: true,
		Ms:     150,
	})
	want = append(want, Interval{
		Signal: false,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: true,
		Ms:     150,
	})
	want = append(want, Interval{
		Signal: false,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: true,
		Ms:     150,
	})
	// [space]
	want = append(want, Interval{
		Signal: false,
		Ms:     150,
	})
	//S
	want = append(want, Interval{
		Signal: true,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: false,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: true,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: false,
		Ms:     50,
	})
	want = append(want, Interval{
		Signal: true,
		Ms:     50,
	})

	morse := "...   ---   ..."
	timing := NewTiming()
	seq := timing.MorseMessageToIntervalSequence(morse)
	if len(seq) != len(want) {
		t.Fatalf(`MorseMessageToIntervalSequence(): expected %d intervals, got %d`, len(want), len(seq))
	} else {
		for i := 0; i < len(seq); i++ {
			if seq[i] != want[i] {
				t.Fatalf(`MorseMessageToIntervalSequence(): different interval @ pos %d: expected Signal: %t Ms: %d; got: Signal: %t Ms: %d`, i, want[i].Signal, want[i].Ms, seq[i].Signal, seq[i].Ms)
			}
		}
	}
}

func TestIntervalPlayer(t *testing.T) {
	morse := "..."
	timing := NewTiming()
	seq := timing.MorseMessageToIntervalSequence(morse)
	if len(seq) != 5 {
		t.Fatalf(`MorseMessageToIntervalSequence("..."): expected 5 intervals, go %d`, len(seq))
	}
	player := NewIntervalSequencePlayer(seq)
	player.OnSignalOn = func() {
		fmt.Println("ON ", time.Now())
	}
	player.OnSignalOff = func() {
		fmt.Println("OFF ", time.Now())
	}
	player.Start()
}
