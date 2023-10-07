package main
import (
  "os"
  "log"
  "time"
  "strings"

  "github.com/faiface/beep"
  "github.com/faiface/beep/wav"
  "github.com/faiface/beep/mp3"
  "github.com/faiface/beep/speaker"
)

// Globals
var ctrl *beep.Ctrl
var format beep.Format
var streamer beep.StreamSeekCloser

// Play a song
func playSong(song *Song, done chan bool) {
  f, err := os.Open(song.FilePath)
  if err != nil {
    log.Fatal(err)
  }

  switch {
  case strings.HasSuffix(song.FilePath, ".mp3"):
    streamer, format, err = mp3.Decode(f)
  case strings.HasSuffix(song.FilePath, ".wav"):
    streamer, format, err = wav.Decode(f)
  default:
    f.Close()
    return
  }

  if err != nil {
    log.Fatal(err)
  }

  speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
  ctrl = &beep.Ctrl{Streamer: streamer, Paused: false}
  
  speaker.Play(beep.Seq(ctrl, beep.Callback(func () {
    done <- true
  })))
}

// Main function
func PlayerLoop(messages chan string){
  // Create the channel
  playNext := make(chan bool, 1)
  playNext <- true

  // Play the first song; send done when done
  for {
    select {
    case <-playNext:
      if(streamer != nil) {
        streamer.Close()
      }
      // Play the next song if available
      nextSong := GetNextSong()
      if(nextSong == nil) {
        messages <- "pause"
      }
      playSong(nextSong, playNext)
    case message := <- messages:
      switch message {
      case "pause": 
        speaker.Lock()
        ctrl.Paused = !ctrl.Paused
        speaker.Unlock()
      case "next":
        if(streamer != nil) {
          streamer.Close()
        }
        nextSong := GetNextSong()
        if(nextSong == nil) {
          messages <- "pause"
        }
        playSong(nextSong, playNext)
      case "back":
        if(streamer != nil) {
          streamer.Close()
        }
        prevSong := GetPrevSong()
        if(prevSong == nil) {
          messages <- "pause"
        }
        playSong(prevSong, playNext)
      }
    case <-time.After(time.Second):
      speaker.Lock()

      position := format.SampleRate.D(streamer.Position()).Round(time.Second)
      length := format.SampleRate.D(streamer.Len()).Round(time.Second)
      
      speaker.Unlock()
    }
  }

}
