package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/indiependente/gw-example/rpc"
	"google.golang.org/grpc"
)

const (
	QUITCMD    = `!quit`
	ECHOCMD    = `!echo`
	REVERSECMD = `!reverse`
)

func main() {

	opts := grpc.WithInsecure()
	address := getDefault("SERVER_ADDR", "0.0.0.0")
	port := getDefault("PORT", "9090")

	log.Printf("Connecting to gRPC server @ %s:%s\n", address, port)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", address, port), opts)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cl := rpc.NewMessageAPIClient(conn)
	funcs := map[string]func(context.Context, *rpc.StringMessage, ...grpc.CallOption) (*rpc.StringMessage, error){
		"!echo":    cl.Echo,
		"!reverse": cl.Reverse,
	}

	stdin := bufio.NewReader(os.Stdin)
	fmt.Println("Commands: !echo, !reverse, !quit")
	for cmd := readString(*stdin, "cmd> "); !strings.EqualFold(cmd, QUITCMD); cmd = readString(*stdin, "cmd> ") {
		f, ok := funcs[cmd]
		if !ok {
			log.Fatal("Unknown command")
		}
		msgstring := readString(*stdin, "Type a message > ")
		log.Println("client >>> " + msgstring)

		msg, err := f(ctx, &rpc.StringMessage{
			Value: msgstring,
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("server >>> " + msg.Value)
	}

	log.Println("Client shutting down...")
}

func getDefault(name, defaultValue string) string {
	value := os.Getenv(name)
	if value != "" {
		return value
	}
	return defaultValue
}

func readString(br bufio.Reader, message string) string {
	fmt.Print(message)
	msg, _ := br.ReadString('\n')
	return strings.Trim(msg, "\n")
}
