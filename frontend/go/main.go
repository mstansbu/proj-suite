package main

var gameServer *GameServer

func main() {
	gameServer = NewGameServer()
	go gameServer.run()

	router := NewRouter()
	router.Run(":3000")

	//router := mux.NewRouter()
}
