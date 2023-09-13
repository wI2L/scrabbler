Output scrabbler.gif

Set FontFamily "SF Mono"
Set FontSize 28
Set Width 1600
Set Height 800
Set Padding 20
Set Framerate 1
Set CursorBlink false
Set PlaybackSpeed 1
Set WindowBar Colorful
Set WindowBarSize 45
Set LineHeight 1

Set Theme { "name": "Challenger Deep", "background": "#1b182c" }

Hide
Type "go build -o scrabbler"
Enter
Type "./scrabbler"
Enter
Sleep 3

Show
Sleep 2
Type "?"
Sleep 2
Type "?"
Sleep 2
Right
Sleep 1
Enter@1.5 2
Sleep 1
Left
Sleep 1
Enter
Sleep 2
Ctrl+G
Sleep 2
Ctrl+G
Sleep 3

Hide
Ctrl+C
Type "rm -f scrabbler"
Enter