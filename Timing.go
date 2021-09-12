package morsica

// Timing for dits and dash
type Timing struct {
	ditMs      int
	ditsInADah float32
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
