# Remote control cheap window openers

A CLI tool to wirelessly open / close cheap window openers

## Intro

So i wanted to be able to open and close my windows automatically and found these relatively cheap window openers on aliexpress. They arrived togeather with a small controller which has simple open-close buttons on them. When wired up you'll have a simple window opener solution - they also include a remote control which transmits in the 433MHz band.

This app mimics these signals and is able to control any amount of window openers. 

As to hardware I use a raspberry pi and bitbang the signals out to one of the included remote which i hacked: i cut the connection from the MCU to the radio and attached the radio to a GPIO pin on the PI - then i'm powering the whole thing from the PI's 3v3 rail. 

Get yours here https://www.aliexpress.com/item/Chain-Supply-Electric-Window-Opener-Office-Building-Window-Opener-Christmas-Home-Automatic-Window-Opener-100-400mm/32868724598.html

## Quick and dirty getting started

### Prepare hardware

Start off by wiring everything up, the window-opener to the included wall controller and the wall controller to the mains - you should now be able to open and close the chain. 

Next up, pull the included remote apart and figure out where the DATA pin is going from the MCU to the 433MHz radio - it's easy to spot - cut it somewhere and prepare 3 leads - one for ground, one for power and one for DATA into the radio. Hook these up to the raspberry pi

Btw - I think any 433MHz transmitter should be able to work just fine - please let me know. I don't know much about sending radio signals but my research leads me to believe what we are dealing with is Manchester encoded ASK/OOK at about 1200 baud.

### Configure window-openers
This app takes a simple yml configuration file describing names of your window openers and their ID. Their ID is a 5 byte value where you should choose 5 random bytes and think of these as your "password". 
The window controller will be paired to these 5 bytes later.

Its absolutely not a secure system as anyone will be able to capture these payloads wirelessly and replay them whenever they want. 

Example:
```
pin: GPIO21
windows:
  Bedroom:
    ID: 177 167 68 228 115
  Otherroom:
    ID: 177 167 132 139 166
```
pin specifies the GPIO pin of your choosing on the raspberry pi and windows specifies your window openers.

Save this file in `/etc/window-openers.yml` on the raspberry

### Build and install window-openers
Its outside the scope of this README to talk about how to acquire golang environment, have a look at https://golang.org if you need this - otherwise:

```
$ go get .
$ GOOS=linux GOARCH=arm go build .
$ scp window-openers pi@<raspberrypi>: 
```

Then login to the raspberry and move the binary somewhere that's in your `path` :)

### Pairing
If everything went well you should now be able to pair with the wall controller - this is done by pressing and holding the Stop key for 10 seconds or so - it will indicate when it's ready for pairing by blinking to you - once it does this go to your terminal and run the pair action

```
# Hold Stop until the controller is in pairing mode
$ window-openers Bedroom pair
```
repeat until all your configured windows are paired.

### Have fun 
You should now be able to control your windows from the command line

```
$ window-openers Bedroom open
$ window-openers Bedroom stop
$ window-openers Bedroom close
```

I personally like to have it configured in cron to wake me up to fresh air in the morning :)

```
0 7 * * * window-openers Bedroom open
0 8 * * * window-openers Bedroom close
```