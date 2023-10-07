package main
import (
  "os"
  "fmt"
  "time"
  
  tea "github.com/charmbracelet/bubbletea"
)

var messages chan string = make(chan string)

type ChangeSongMsg *Song
type ChangeDurationMsg struct {
  duration time.Duration
  position time.Duration
}

type model struct {
  song *Song
  duration time.Duration
  position time.Duration
}

func initialModel() model {
  return model{
    song: nil,
  }
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch msg.String() {
    case "q":
      return m, tea.Quit
    case "p":
      messages <- "pause"
    case "n":
      messages <- "next"
    case "b":
      messages <- "back"
    }
  case ChangeSongMsg:
    m.song = msg
  case ChangeDurationMsg:
    m.duration = msg.duration
    m.position = msg.position
  }

  return m, nil
}

func (m model) View() string {
  if(m.song != nil) {
    return fmt.Sprintf(
      "%s - %s [%s/%s]",
      m.song.Title,
      m.song.Artist,
      m.position,
      m.duration,
    )
  }

  return ""
}

func main() {
  // Initialize the playlist
  InitializePlaylist([]string{
    "/home/cheshire/Music/Brand New/Deja Entendu/01 Tautou.mp3",
    "/home/cheshire/Music/Brand New/Deja Entendu/07 Jaws Theme Swimming.mp3",
    "/home/cheshire/Music/Modest Mouse/Interstate 8/Modest Mouse - Interstate 8 - 11 - Edit the Sad Parts.mp3",
  })
  
  p := tea.NewProgram(initialModel())
  
  go PlayerLoop(p, messages)

  if _, err := p.Run(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1) 
  }
}
