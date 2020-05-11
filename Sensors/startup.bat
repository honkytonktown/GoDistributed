go build sensor.go

FOR /l %%x in (1, 1, 3) DO (
    start cmd.exe /c sensor -freq=1
)

