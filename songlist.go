/* This file controls the "playlist" that songs are played from.
   It keeps track of upcoming songs as well as history of songs that have been played.
   It provides functions for going to the next or previous song, as well as
   toggling Shuffle and Repeat functionality. 
*/

package main
import (
  "os"
  "path/filepath"
  "math/rand"
  "time"
  
  "github.com/dhowden/tag"
)

// Song: keeps track of file and metadata for every song
type Song struct {
  FilePath string
  Title string
  Artist string
  Album string
}

// Stacks for upcoming songs and previously played songs
type stack struct {
  data []*Song
}

func (s *stack) push(elem *Song) {
  s.data = append(s.data, elem)
}
func (s *stack) pop() *Song {
  if(len(s.data) == 0) {
    return nil
  }
  n := len(s.data) - 1
  elem := s.data[n]
  s.data[n] = nil
  s.data = s.data[:n]
  return elem
}
func (s *stack) shuffle() {
  rand.Seed(time.Now().UnixNano())
  n := len(s.data)
  for i:=n-1; i>0; i-- {
    j := rand.Intn(i)
    temp := s.data[j]
    s.data[j] = s.data[i]
    s.data[i] = temp
  }
}

// keeping track of different things
var playlist []*Song = nil
//var songplaylist []song = make([]song, 0)
var nextUp *stack = &stack{data:make([]*Song, 0)}
var history *stack = &stack{data:make([]*Song, 0)}
var currentSong *Song = nil
var shuffle bool = false // not yet implemented
var repeat bool = true

/* Initializes everything from an initial array of songs (file paths)
   This sets the `playlist` variable which lets us know which songs
   we have to work with and will not change unless InitializePlaylist is called again.
   It also initializes nextUp with all of the songs. 
 */
func InitializePlaylist(initialPlaylist []string) {
  playlist = make([]*Song, len(initialPlaylist))
  for i:=len(initialPlaylist)-1; i >= 0; i-- {
    f, err := os.Open(initialPlaylist[i])
    if err != nil {
      continue
    }
    m, err := tag.ReadFrom(f)

    var title string
    var artist string
    var album string
    
    if err != nil {
      title = filepath.Base(initialPlaylist[i])
      artist = "No Data"
      album = "No Data"
    } else {
      title = m.Title()
      if title == "" {
        title = filepath.Base(initialPlaylist[i])
      }

      artist = m.Artist()
      if artist == "" {
        artist = "No Data"
      }

      album = m.Album()
      if album == "" {
        album = "No Data"
      }
    }

    var song *Song = &Song{
      FilePath: initialPlaylist[i],
      Title: title,
      Artist: artist,
      Album: album,
    }

    playlist[i] = song
    
    f.Close()
  }
  
  populateNextUp()
}

func populateNextUp() {
  // Going through array backwards because we are pushing onto a stack
  for i := len(playlist) - 1; i >= 0; i-- {
    nextUp.push(playlist[i])
  }

  if shuffle {
    nextUp.shuffle()
  }
}

func ToggleRepeat() bool {
  repeat = !repeat
  return repeat
}

func ToggleShuffle() bool {
  shuffle = !shuffle
  nextUp.shuffle()
  return shuffle
}

func GetNextSong() *Song {
  if(currentSong != nil) {
    history.push(currentSong)
  }

  currentSong = nextUp.pop()
  if(currentSong == nil && repeat) {
    // We ran out of songs and need to repeat the playlist
    populateNextUp()
    currentSong = nextUp.pop()
  } 
  return currentSong
}

func GetPrevSong() *Song {
  if(currentSong == nil) {
    return nil
  }

  nextUp.push(currentSong)
  currentSong = history.pop()
  return currentSong
}
