package main

import (
	"bufio"
	"fmt"
	"log"
	"mighty/services"
	"os"
	"strings"
	"testing"
)

func TestGame(t *testing.T) {
	var game services.Game

	game.Id = 1

	game.Command("init 0")
	//game.Print()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("start")
	fmt.Println("---------------------")

	str := ""
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		//text = strings.Replace(text, "\n", "", -1)

		if strings.Compare("exit\n", text) == 0 {
			fmt.Println("end")
			break
		}

		if strings.Compare("print\n", text) == 0 {
			game.Print()
		} else {
			log.Println("+++++++++++++++++++++++++++++++++++++++++")
			log.Println(str)
			log.Println("+++++++++++++++++++++++++++++++++++++++++")
			game.Command(str)
		}
	}
}
