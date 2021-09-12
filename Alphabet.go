package morsica

import (
	"strconv"
	"strings"
)

// Alphabet (and digit) conversion between ascii and morse
type Alphabet struct {
	morseDigits      [10]string
	morseLetters     [26]string
	interLetterSpace int
	interWordSpace   int
}

// NewAlphabet creates a pointer to a new Alphabet
func NewAlphabet() *Alphabet {
	a := new(Alphabet)
	a.makeMorseDigits()
	a.makeMorseLetters()
	a.interLetterSpace = 3
	a.interWordSpace = 7
	return a
}

///
/// encoding
///

//Encode plaintext into a morse sequence
func (a *Alphabet) Encode(txt string) string {
	var morse string
	lTxt := strings.ToLower(txt)
	words := strings.Split(lTxt, " ")
	numWords := len(words)
	for w := 0; w < numWords; w++ {
		word := words[w]
		wordLen := len(word)
		for i := 0; i < wordLen; i++ {
			if isDigitASCII(word[i]) {
				morse += a.encodeDigitASCII(word[i])
			} else if isLetterASCII(word[i]) {
				morse += a.encodeLetterASCII(word[i])
			}
			if i < (wordLen - 1) {
				morse += a.encodeInterLetterSpace()
			}
		}
		if w < (numWords - 1) {
			morse += a.encodeInterWordSpace()
		}
	}
	return morse
}

// GetInterWordSpaces returns the number of spaces between words in Morse code
func (a *Alphabet) GetInterWordSpaces() int {
	return a.interWordSpace
}

// GetInterSymbolSpaces returns the numver of spaces between letter/numbers withing a single word in Morse code
func (a *Alphabet) GetInterSymbolSpaces() int {
	return a.interLetterSpace
}

// SetInterWordSpaces sets the number of spaces between words in Morse code
func (a *Alphabet) SetInterWordSpaces(n int) {
	if n < 2 {
		return
	}
	a.interWordSpace = n
}

// SetInterSymbolSpaces sets the numver of spaces between letter/numbers withing a single word in Morse code
func (a *Alphabet) SetInterSymbolSpaces(n int) {
	if n < 1 {
		return
	}
	a.interLetterSpace = n
}

func isDigitASCII(char byte) bool {
	return (char >= 48 && char <= 57)
}

func isLetterASCII(char byte) bool {
	return (char >= 97 && char <= 122)
}

func (a *Alphabet) encodeDigitASCII(char byte) string {
	idx := char - 48
	return a.morseDigits[idx]
}

func (a *Alphabet) encodeLetterASCII(char byte) string {
	idx := char - 97
	return a.morseLetters[idx]
}

func encodeSpace(spaceLen int) string {
	var s string
	for i := 0; i < spaceLen; i++ {
		s += " "
	}
	return s
}

func (a *Alphabet) encodeInterLetterSpace() string {
	return encodeSpace(a.interLetterSpace)
}

func (a *Alphabet) encodeInterWordSpace() string {
	return encodeSpace(a.interWordSpace)
}

///
/// decoding
///

// Decode morse sequence to plain text
func (a *Alphabet) Decode(morse string) string {
	iWspace := a.encodeInterWordSpace()
	words := strings.Split(morse, iWspace)
	numWords := len(words)
	if numWords < 1 {
		return ""
	}
	var txt string
	for i := 0; i < numWords; i++ {
		txt += a.decodeWord(words[i])
		if i < (numWords-1) && txt != "" {
			txt += " "
		}
	}
	return txt
}

func (a *Alphabet) decodeWord(word string) string {
	iLspace := a.encodeInterLetterSpace()
	symbols := strings.Split(word, iLspace)
	numSymbols := len(symbols)
	if numSymbols < 1 {
		return ""
	}
	var txt string
	for i := 0; i < numSymbols; i++ {
		symbol := symbols[i]
		if isValidMorseSequence(symbol) {
			if isMorseDigit(symbol) {
				txt += decodeDigit(symbol)
			} else {
				txt += a.decodeLetter(symbol)
			}
		}
	}
	return txt
}

func decodeDigit(morseSeq string) string {
	var val int
	if morseSeq[0] == '.' {
		val = 1
		for i := 1; i < len(morseSeq); i++ {
			if morseSeq[i] == '.' {
				val++
			} else {
				i = len(morseSeq) + 2
			}
		}
	} else {
		val = 6
		for i := 1; i < len(morseSeq); i++ {
			if morseSeq[i] == '-' {
				val++
			} else {
				i = len(morseSeq) + 2
			}
			if val > 9 {
				val = 0
			}
		}
	}
	return strconv.Itoa(val)
}

func (a *Alphabet) decodeLetter(morseSeq string) string {
	//TODO: group letters by len, so we can do a faster searc
	for i := 0; i < len(a.morseLetters); i++ {
		if morseSeq == a.morseLetters[i] {
			asciiCode := 97 + i
			asciiRune := rune(asciiCode)
			return string(asciiRune)
		}
	}
	return ""
}

func isMorseDigit(morseSeq string) bool {
	return (len(morseSeq) == 5)
}

func isValidMorseSymbol(char byte) bool {
	return (char == '.' || char == '-')
}

func isValidMorseSequence(morseSeq string) bool {
	seqLen := len(morseSeq)
	if seqLen < 2 || seqLen > 5 {
		return false
	}
	for i := 0; i < seqLen; i++ {
		if !isValidMorseSymbol(morseSeq[i]) {
			return false
		}
	}
	return true
}

// DoesItLookLikeMorse returns true if a text starts with a sequence of Morse code.
// It's meant for a quick guess, it does not evaluate the whole text
func DoesItLookLikeMorse(txt string) bool {
	split := strings.Split(txt, " ")
	if len(txt) > 0 {
		return isValidMorseSequence(split[0])
	} else {
		return false
	}
}

///
/// alphabet init
///

func (a *Alphabet) makeMorseDigits() {
	a.morseDigits = [10]string{
		"-----", //0
		".----", //1
		"..---", //2
		"...--", //3
		"....-", //4
		".....", //5
		"-....", //6
		"--...", //7
		"---..", //8
		"----.", //9
	}
}

func (a *Alphabet) makeMorseLetters() {
	a.morseLetters = [26]string{
		".-",   //a
		"-...", //b
		"-.-.", //c
		"-..",  //d
		".",    //e
		"..-.", //f
		"--.",  //g
		"....", //h
		"..",   //i
		".---", //j
		"-.-",  //k
		".-..", //l
		"--",   //m
		"-.",   //n
		"---",  //o
		".--.", //p
		"--.-", //q
		".-.",  //r
		"...",  //s
		"-",    //t
		"..-",  //u
		"...-", //v
		".--",  //w
		"-..-", //x
		"-.-",  //y
		"--.",  //z
	}
}
