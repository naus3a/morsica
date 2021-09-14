# Morsica
*CLI* tool to encode/decode Morse code.
Usage:
```
morsica [optional params] [optional commands] [message to encode/decode]
```

## Commands
* `-e`: encode plaintext to Morse
* `-d`: decode Morse to plaintext

If no command is specified, `morsica` will try to guess what you want based on your input

## Misc Params
* `-h`: prints a help screen
* `-s`: silent mode, disables all messages and spits out your nakes output

## Morse format Params
* `-t [ms]`: sets the dit duration in ms. **Default is 50ms**
* `-sW [num]`: sets the inter word number of spaces in Morse code. **Default is 7 spaces**.
* `-sS [num]`: sets the inter symbol number of spaces within a single word in Morse code. **Default is 3 spaces**

## Playback params
* `-p`: enable playback of a Morse sequence (alone it will do very little)
* `-cOn [cmd]`: command called when signal is on
* `-cOff [cmd]`: command called when signal is off 