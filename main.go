package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var flag string
var block string
var padding string
var offset int
var alphabet []rune

// flag : 32 F{...}

// 53 // crabby patty // 31 // 7(picoCT) //

// 11

func main() {
	flag = "picoCTF{"
	block = "ode is: "
	padding = strings.Repeat("A", 11)
	offset = 41
	alphabet = []rune("!$%&/{}[]=?^-_+#@*.abcdefghijklmnopqrstuvwxyz0123456789")

	fmt.Print(flag)
	for i := 1; i <= 29; i++ {
		letter := bruteforce(i)
		flag += letter
	}
	fmt.Printf("\r%s}\n\n", flag)
}

func sendMsg(msg string) string {
	conn, err := net.Dial("tcp", "2018shell.picoctf.com:34490")

	if err != nil {
		log.Fatal(err)
	}

	if _, err = bufio.NewReader(conn).ReadString('\n'); err != nil {
		log.Fatal(err)
	}

	buff := bytes.NewBufferString(msg + "\n")

	buff.WriteTo(conn)

	enc, _ := bufio.NewReader(conn).ReadString('\n')

	return enc
}

func bruteforce(position int) string {
	var wg sync.WaitGroup
	ch := make(chan string)

	for _, letter := range alphabet {
		wg.Add(1)
		go guessLetter(position, string(letter), &wg, ch)
	}

	wg.Add(1)
	go printAnimation(&wg, ch)

	letter := <-ch
	wg.Wait()

	return letter
}

func printAnimation(wg *sync.WaitGroup, ch chan string) {
	for _, letter := range alphabet {
		select {
		case _ = <-ch:
			wg.Done()
			return
		default:
			fmt.Printf("\r%s", flag+string(letter))
			time.Sleep(100 * time.Millisecond)
		}
	}
	_ = <-ch
	wg.Done()
}

func guessLetter(position int, letter string, wg *sync.WaitGroup, ch chan string) {
	payload := padding + strings.Repeat("B", 32-position) + block + flag + letter + strings.Repeat("A", offset-position)
	enc := []rune(sendMsg(payload))
	guess := string(enc[6*32 : 7*32])
	real := string(enc[11*32 : 12*32])

	// fmt.Printf("%s %s %s %s\n", string(alphabet[i]), guess, real, string(enc))

	if guess == real {
		ch <- string(letter)
		ch <- string(letter)
	}
	wg.Done()
}
