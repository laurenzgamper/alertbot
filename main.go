package main

func main() {
	config := readConfigFile()

	messages := make(chan string, 100)

	go listenWebHook(config.Listen, messages)
	startBot(config, messages)
}