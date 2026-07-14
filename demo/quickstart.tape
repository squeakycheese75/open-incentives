# Regenerate: brew install vhs && vhs demo/quickstart.tape
# Requires the stack already built at least once (docker compose build).
Output demo/quickstart.gif

Set Shell "bash"
Set FontSize 18
Set Width 1200
Set Height 650
Set Theme "Dracula"
Set TypingSpeed 45ms
Set Padding 20

Type "docker compose --profile demo-store up -d"
Enter
Sleep 3s

Type "clear"
Enter
Sleep 500ms

Type "./demo/quickstart.sh"
Enter
Sleep 4500ms
