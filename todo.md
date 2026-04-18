# Tasks:
## Frontend:
- [ ] Create global styles.
- [x] Implement home page mockup in Bubbletea.
    - [x] Style Add Panel
    - [x] Style Filter Panel
    - [x] Style Sort Panel
- [x] Implement cursor moving between sections.
    - This should use H, J, K, L in the directions expected by nvim.
    - [x] Cursor moves between add and list.
    - [x] Allow Cursor to move vertically.
    - [x] Cursor move from any sidebar partial and the list, and correctly select the same sidebar view when moving back.
- [ ] Add relevant columns to table.
- [ ] Make table size columns automatically, truncate past a certain amount of characters.
- [ ] Make main.go use an array of tea.Models to send messages
    - Most key messages can probably just generically go to whichever view is active.
- [ ] Figure out blink on the filter text inputs

## Backend:
- [ ] Create SQL Database
    - [ ] Write code for generating SQLite database.
    - [ ] Write code for "connecting" to SQLite database.
- [ ] Swap in JSON object code for arrays.
    - https://www.sqlite.org/json1.html 

## External:
- [ ] Keep processing my existing watchlist.
- [ ] Look into a copyleft license.
