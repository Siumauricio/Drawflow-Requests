package main

func main() {

	mux := Route()
	server := NewServer(mux)
	server.Run()
}