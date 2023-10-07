package main
import (
  "fmt"
)

func main() {
  // Initialize the playlist
  InitializePlaylist([]string{
    "/home/cheshire/Music/Brand New/Deja Entendu/01 Tautou.mp3",
    "/home/cheshire/Music/Brand New/Deja Entendu/07 Jaws Theme Swimming.mp3",
    "/home/cheshire/Music/Modest Mouse/Interstate 8/Modest Mouse - Interstate 8 - 11 - Edit the Sad Parts.mp3",
    "/home/cheshire/Music/Testing/futuristic.wav",
  })

  messages := make(chan string)
  go PlayerLoop(messages)

  for {
    var first string

    fmt.Print("Enter command: ")
    fmt.Scanln(&first)
    if(first == "exit") {
      break
    }

    if(first == "pause") {
      messages <- "pause"
    }

    if(first == "next") {
      messages <- "next"
    }

    if(first == "back") {
      messages <- "back"
    }
  }
}
