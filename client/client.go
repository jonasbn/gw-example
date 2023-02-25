package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	service "github.com/indiependente/gw-example/rpc/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	QUITCMD    = `!quit`
	ECHOCMD    = `!echo`
	REVERSECMD = `!reverse`
)

func main() {

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	address := getDefault("SERVER_ADDR", "0.0.0.0")
	port := getDefault("PORT", "9090")

	log.Printf("Connecting to gRPC server @ %s:%s\n", address, port)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", address, port), opts)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	cl := service.NewMessageAPIServiceClient(conn)
	supportedCmds := map[string]struct{}{
		"!echo":    {},
		"!reverse": {},
	}

	stdin := bufio.NewReader(os.Stdin)
	fmt.Println(cmds())
	for cmd := readString(*stdin, "cmd> "); !strings.EqualFold(cmd, QUITCMD); cmd = readString(*stdin, "cmd> ") {
		if _, ok := supportedCmds[cmd]; !ok {
			log.Println("unsupported command")

			continue
		}

		msgstring := readString(*stdin, "Type a message > ")
		log.Println("client >>> " + msgstring)

		var msgResponse string
		switch cmd {
		case ECHOCMD:
			msg, err := cl.Echo(ctx, &service.EchoRequest{Value: msgstring})
			if err != nil {
				log.Println(err)
			}
			msgResponse = msg.Value

		case REVERSECMD:
			msg, err := cl.Reverse(ctx, &service.ReverseRequest{Value: msgstring})
			if err != nil {
				log.Println(err)
			}
			msgResponse = msg.Value
		}
		log.Println("server >>> " + msgResponse)
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

func cmds() string {
	sb := strings.Builder{}
	sb.WriteString("Commands list: ")
	sb.WriteString(ECHOCMD)
	sb.WriteString(", ")
	sb.WriteString(REVERSECMD)
	sb.WriteString(", ")
	sb.WriteString(QUITCMD)
	return sb.String()
}
