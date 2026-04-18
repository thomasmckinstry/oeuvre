# MediaLogger TUI

This project is making a terminal user interface for a user to track movies, shows, and other media they have on their watchlist. It'll allow users to keep track of the status of a given work (Completed, Pending, etc.) as well as taking notes and writing reviews on works.

# Architecture
## Packages:
### Bubbletea:
Bubbletea is a golang framework for making terminal user interfaces.
#### Lipgloss
Lipgloss is a package built for use with bubbletea that allows for styling of the user interface similar to CSS.

### SQLite:
SQLite will provide an isolated local database to store the data.

# Directory Structure:

# Debugging:
```
DEBUG=1 go run main.go
```
```
if len(os.Getenv("DEBUG")) > 0 {
    log.Println("{ DEBUG MESSAGE }")
}
```
