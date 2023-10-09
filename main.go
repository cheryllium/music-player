package main
import (
  "os"
  "fmt"
  "time"
  "path/filepath"
  "strings"
  
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/bubbles/progress"
  "github.com/charmbracelet/bubbles/help"
  "github.com/charmbracelet/bubbles/key"
)

var messages chan string = make(chan string)

type ChangeSongMsg *Song
type ChangeDurationMsg struct {
  duration time.Duration
  position time.Duration
}

type keyMap struct {
  P key.Binding
  N key.Binding
  B key.Binding
  Q key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
  return []key.Binding{k.P, k.N, k.B, k.Q}
}
func (k keyMap) FullHelp() [][]key.Binding {
  return [][]key.Binding{
    {k.P, k.Q}, {k.N}, {k.B},
  }
}

var keys = keyMap{
  P: key.NewBinding(
    key.WithKeys("p"),
    key.WithHelp("p", "pause/unpause"),
  ),
  N: key.NewBinding(
    key.WithKeys("n"),
    key.WithHelp("n", "next"),
  ),
  B: key.NewBinding(
    key.WithKeys("b"),
    key.WithHelp("b", "back"),
  ),
  Q: key.NewBinding(
    key.WithKeys("q"),
    key.WithHelp("q", "quit"),
  ),
}

type model struct {
  keys keyMap
  help help.Model
  song *Song
  progress progress.Model
  duration time.Duration
  position time.Duration
}

func initialModel() model {
  return model{
    keys: keys,
    help: help.New(),
    song: nil,
    progress: progress.New(progress.WithDefaultGradient(), progress.WithoutPercentage()),
  }
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    m.help.Width = msg.Width
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
    percentage := float64(m.position) / float64(m.duration)
    cmd := m.progress.SetPercent(percentage)
    return m, cmd
  case progress.FrameMsg:
    progressModel, cmd := m.progress.Update(msg)
    m.progress = progressModel.(progress.Model)
    return m, cmd
  }

  return m, nil
}

func (m model) View() string {
  if(m.song != nil) {
    return fmt.Sprintf(
      "%s - %s\n%s [%s / %s]\n%s",
      m.song.Title,
      m.song.Artist,
      m.progress.View(),
      m.position,
      m.duration,
      m.help.FullHelpView(m.keys.FullHelp()),
    )
  }

  return ""
}

func main() {
  // Read in the first argument to the command; the directory to play music from
  args := os.Args[1:]
  directory := args[0]

  // Attempt to read the directory
  files, err := os.ReadDir(directory)
  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  // Get files from the directory
  var filenames = make([]string, 0)
  for _, file := range files {
    // Skip if the DirEntry is a directory
    if file.IsDir() {
      continue
    }
    // Only add to playlist if it is an mp3 or wav file
    if strings.HasSuffix(file.Name(), ".mp3") || strings.HasSuffix(file.Name(), ".wav") {
      filenames = append(filenames, filepath.Join(directory, file.Name()))
    }
  }
  
  // Initialize the playlist
  InitializePlaylist(filenames)
  
  p := tea.NewProgram(initialModel())
  
  go PlayerLoop(p, messages)

  if _, err := p.Run(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1) 
  }
}
