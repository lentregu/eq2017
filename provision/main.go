package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"tools-support/screen"

	"github.com/TDAF/gologops"
)

var clear screen.ClearWindow

var actions map[int]string

var stdin *bufio.Reader

var (
	bdServer string
	bdPort   int
	bdName   string
	bdColl   string
)

func init() {

	var fatalErr error
	defer func() {
		if fatalErr != nil {
			flag.PrintDefaults()
			log.Fatalln(fatalErr)
		}
	}()

	flag.Parse()

	if fatalErr != nil {
		return
	}

	clear = screen.NewClearScreenFunction(screen.DARWIN)
	stdin = bufio.NewReader(os.Stdin)

	actions = map[int]string{
		1: "createFaceList",
		2: "listFacesList",
		3: "facesInAList",
		4: "addFace",
		5: "whois?",
		6: "whoisB64?",
		7: "createSpeakerProfile",
		8: "end",
	}
}

func main() {

	var option string
	for option != "end" {
		option = menu()
		switch {
		case option == "createFaceList":
			id, err := createFaceList()
			if err != nil {
				log.Fatal(err)
			} else {
				gologops.Infof("The list %s has been created", id)
			}
		case option == "listFacesList":
			list, err := getFaceList()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("Lists: %s\n", list)
			}
		case option == "facesInAList":
			list, err := getFacesInAList()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("Lists: %s\n", list)
			}
		case option == "addFace":
			photoID, err := addFace()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("The PhotoID is: %s\n", photoID)
			}
		case option == "whois?":
			userInfo, err := whois()
			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("The UserInfo is: %s\n", userInfo)
			}
		case option == "whoisB64?":
			whoIsBase64()
		case option == "createSpeakerProfile":
			id, err := createSpeakerProfile()
			if err != nil {
				log.Fatal(err)
			} else {
				gologops.Infof("The Profile %s has been created ss", id)
			}

		case option == "end":
			return
		}
		fmt.Println("\nPress <ENTER>......")
		stdin.ReadLine()
	}
	os.Exit(0)
}

func menu() string {
	option := 0
	for option < 1 || option > 8 {
		clear := screen.NewClearScreenFunction(screen.DARWIN)
		clear()
		fmt.Println("1. Create Face List")
		fmt.Println("2. List of Faces lists")
		fmt.Println("3. List Faces in a list")
		fmt.Println("4. Add Face")
		fmt.Println("5. whois??")
		fmt.Println("6. whoisB64??")
		fmt.Println("7. Create Speaker Profile")
		fmt.Println("8. Exit")
		fmt.Printf("\nChoose an option....:")
		if _, err := fmt.Fscanf(stdin, "%d", &option); err != nil {
			// In case of not introducing a number
			option = 0
		}
		stdin.ReadLine() //This line is necessary to flush the buffer because there is a "\n" left

	}
	return actions[option]
}
