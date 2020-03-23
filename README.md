## bt-autoconnect
This program will detect bluetooth event changes via `/dev/input/event0` and
if a change occurs, it will try to connect to one of the paired bluetooth
devices.

Before this tool is able to work, you have to manually pair and trust all
your devices.

#### Cross compilation
On your development machine run the following line to build a `bt-autoconnect` binary that can run on a
raspberry pi:

    GOOS=linux GOARCH=arm go build -o bt-autoconnect -v bt-autoconnect.go
