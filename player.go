package main
import (
  "os"
  "log"
  "time"
  "fmt"
  "strings"

  tea "github.com/charmbracelet/bubbletea"
  
  "github.com/faiface/beep"
  "github.com/faiface/beep/wav"
  "github.com/faiface/beep/mp3"
  "github.com/faiface/beep/speaker"
)

// Globals
var ctrl *beep.Ctrl
var format beep.Format
var streamer beep.StreamSeekCloser
var initialized bool = false

// Play a song
func playSong(song *Song, playNext chan bool, p *tea.Program) {
  f, err := os.Open(song.FilePath)
  if err != nil {
    fmt.Println(err)
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

  speaker.Clear()
  if(!initialized) {
    speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
    initialized = true
  }

  speaker.Lock()

  ctrl = &beep.Ctrl{
    Streamer: beep.Seq(streamer, beep.Callback(func () {
      // Hacky workaround for erroneous callbacks
      // Sometimes this callback runs for a skipped song when I press the skip button, idk why
      // This will simply ignore that case
      position := format.SampleRate.D(streamer.Position()).Round(time.Second)
      length := format.SampleRate.D(streamer.Len()).Round(time.Second)
      if(position != length) {
        return
      }
      // -- End hacky workaround
      
      streamer.Close()
      playNext <- true
    })),
    Paused: false,
  }
  speaker.Unlock()
  
  speaker.Play(ctrl)

  p.Send(ChangeSongMsg(song))
  position := format.SampleRate.D(streamer.Position()).Round(time.Second)
  length := format.SampleRate.D(streamer.Len()).Round(time.Second)
  
  p.Send(ChangeDurationMsg{
    duration: length,
    position: position,
  })
}

// Main function
func PlayerLoop(p *tea.Program, messages chan string){
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
      playSong(nextSong, playNext, p)
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
        playSong(nextSong, playNext, p)
      case "back":
        if(streamer != nil) {
          streamer.Close()
        }
        prevSong := GetPrevSong()
        if(prevSong == nil) {
          messages <- "pause"
        }
        playSong(prevSong, playNext, p)
      }
    case <-time.After(time.Second):
      speaker.Lock()

      position := format.SampleRate.D(streamer.Position()).Round(time.Second)
      length := format.SampleRate.D(streamer.Len()).Round(time.Second)
      
      p.Send(ChangeDurationMsg{
        duration: length,
        position: position,
      })
      
      speaker.Unlock()
    }
  }

}
