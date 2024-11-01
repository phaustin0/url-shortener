package main

func main() {
	listenAddr := ":8000"
	s := NewServer(listenAddr)
	s.Listen()
}
