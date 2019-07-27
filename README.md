# go-chatroom
> This is a chatroom by golang on linux.

ip#message (send message to ip)

server.go
ready logger -> server.listen -> go processInfo (conn in connMap) -> go handleMessage (message in queue)

client.go
client.dial -> go sendMessage (bufio.NewReader(os.Stdin)) -> show message
