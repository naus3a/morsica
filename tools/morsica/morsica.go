package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/naus3a/morsica"
)

//ToolMode is the mode the tool is operating
type ToolMode int

//mode the tool is operating
const (
	ModeEncoding ToolMode = iota
	ModeDecoding
	ModeNone
)

var txtInput string
var sCmdOn string = ""
var sCmdOff string = ""
var ditMs int = -1
var curMode ToolMode = ModeNone
var bVerbose bool = true
var bPlayback bool = false

///
/// CLI
///

func printHelp() {
	fmt.Println("Morsica is a morse code encoder/decoder")
	fmt.Println("Usage:")
	fmt.Println("\tmorsica [optional params] [optional command] [text to encode or decode]")
	fmt.Println("Commands:")
	fmt.Println("\t-e: encode plaintext to morse code")
	fmt.Println("\t-d: decode morse code to plaintext")
	fmt.Println("Misc Params:")
	fmt.Println("\t-h: print this help screen")
	fmt.Println("\t-s: silent mode; no verbose messages, just the output")
	fmt.Println("Morse Format Params:")
	fmt.Println("\t-t [ms]: sets the duration of a dit in ms (DEFAULT: 50)")
	fmt.Println("\t-sW [num]: sets the number of spaces between words in Morse code (DEFAULT: 7)")
	fmt.Println("\t-sS [num]: sets the number of spaces between symbols within a single word in Morse code (DEFAULT: 3)")
	fmt.Println("Playback Params:")
	fmt.Println("\t-p: playback encoded message")
	fmt.Println("\t-cOn [cmd]: command to call when signal starts")
	fmt.Println("\t-cOff [cmd]: command to call when signal stops")
}

func printVerbose(msg string) {
	if !bVerbose {
		return
	}
	fmt.Println(msg)
}

func parseCliArgs(alphabet *morsica.Alphabet) {
	args := os.Args[1:]
	if len(args) < 1 {
		printHelp()
		return
	}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h":
			printHelp()
		case "-s":
			bVerbose = false
		case "-e":
			curMode = ModeEncoding
		case "-d":
			curMode = ModeDecoding
		case "-t":
			val, err := getIntArgValue(&args, i)
			if err == nil {
				ditMs = val
				i++
			}
		case "-sW":
			val, err := getIntArgValue(&args, i)
			if err == nil {
				alphabet.SetInterWordSpaces(val)
				i++
			}
		case "-sS":
			val, err := getIntArgValue(&args, i)
			if err == nil {
				alphabet.SetInterSymbolSpaces(val)
				i++
			}
		case "-p":
			bPlayback = true
		case "-cOn":
			sCmdOn = getArgValue(&args, i)
			if sCmdOn != "" {
				i++
			}
		case "-cOff":
			sCmdOff = getArgValue(&args, i)
			if sCmdOff != "" {
				i++
			}
		default:
			txtInput = args[i]
		}
	}
}

func argHasValue(args *[]string, idx int) bool {
	return (len(*args) > idx)
}

func getArgValue(args *[]string, idx int) string {
	if argHasValue(args, idx) {
		return (*args)[idx+1]
	}
	return ""

}

func getIntArgValue(args *[]string, idx int) (int, error) {
	s := getArgValue(args, idx)
	if s == "" {
		err := errors.New("Arg has no value")
		i := -1
		return i, err
	}
	return strconv.Atoi(s)
}

func guessModeIfNeeded() {
	if curMode != ModeNone {
		return
	}
	printVerbose("No operating mode specified: trying to guess...")
	if morsica.DoesItLookLikeMorse(txtInput) {
		curMode = ModeDecoding
	} else {
		curMode = ModeEncoding
	}
}

///
/// playback
///

func playback(morse string) {
	printVerbose("Starting Playback!")
	timing := morsica.NewTiming()
	if ditMs >= 0 {
		timing.SetDitMs(ditMs)
	}
	player := morsica.NewIntervalSequencePlayer(timing.MorseMessageToIntervalSequence(morse))
	var cmdOn *exec.Cmd = nil
	var cmdOff *exec.Cmd = nil
	var outOn []byte
	var outOff []byte
	setupCommand(&sCmdOn, cmdOn, &outOn, &player.OnSignalOn)
	setupCommand(&sCmdOff, cmdOff, &outOff, &player.OnSignalOff)
	player.Start()
}

func setupCommand(sCmd *string, cmd *exec.Cmd, output *[]byte, cbk *morsica.SignalCallback) {
	if *sCmd == "" {
		return
	}
	printVerbose("Registering command: " + *sCmd)
	ss := strings.Split(*sCmd, " ")
	var sMainCmd string
	if len(ss) > 1 {
		sMainCmd = ss[0]
		var sCmdArgs []string = ss[1:len(ss)]
		cmd = exec.Command(sMainCmd, sCmdArgs[0:]...)
	} else {
		sMainCmd = *sCmd
		cmd = exec.Command(sMainCmd)
	}
	var err error
	*output, err = cmd.Output()
	if err != nil {
		fmt.Println("Cannout create an output for " + *sCmd)
	}
	*cbk = func() {
		cmd.Run()
		fmt.Println(string(*output))
	}
}

///
/// endocding/decoding
///

func performEncoding(alphabet *morsica.Alphabet) {
	printVerbose("Encoding message:")
	output := alphabet.Encode(txtInput)
	fmt.Println(output)
	if bPlayback {
		playback(output)
	}
}

func performDecoding(alphabet *morsica.Alphabet) {
	printVerbose("Decoding message:")
	output := alphabet.Decode(txtInput)
	fmt.Println(output)
}

///
/// main
///

func main() {
	alphabet := morsica.NewAlphabet()
	parseCliArgs(alphabet)
	guessModeIfNeeded()
	switch curMode {
	case ModeNone:
		printHelp()
	case ModeEncoding:
		performEncoding(alphabet)
	case ModeDecoding:
		performDecoding(alphabet)
	}
}
