package interfaces

type IBot interface {
	Connect()
	Join(string)
	Broadcast(string)
	Run()
}