package main

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"strings"

	"github.com/apostolistselios/atm-simulator-rpc/api"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("error connecting to the server:", err)
	}
	fmt.Println("WELCOME TO THE ATM.")

	username, err := getCredentials()
	if err != nil {
		log.Fatal(err)
	}

	// Make the user request to the server.
	if err := verifyUser(client, username); err != nil {
		log.Fatal(err)
	}

	mainLoop(client, username)
}

// getCredentials gets the username from the user and returns it.
func getCredentials() (string, error) {
	var username string
	fmt.Print("USERNAME: ")
	if _, err := fmt.Scanf("%s\n", &username); err != nil {
		return "", errors.New("error parsing the username")
	}
	return username, nil
}

// verifyUser makes a RPC call to the server in order to verify
// if that specific username exists in the ATM's Database.
func verifyUser(client *rpc.Client, username string) error {
	var reply int
	err := client.Call("ATM.VerifyUser", username, &reply)
	if err != nil {
		return err
	}
	return nil
}

// mainLoop function executes the main loop of the client handling
// all the user's actions after his verification.
func mainLoop(client *rpc.Client, username string) {
	answer := "y"
	for answer == "y" || answer == "Y" {
		userPrompt()
		var choice string
		if _, err := fmt.Scanf("%s\n", &choice); err != nil {
			log.Println("error incorrect choice")
			continue
		}

		choice = strings.ToLower(choice)
		if choice == "w" || choice == "d" {
			var amount int
			fmt.Print("PLEASE ENTER THE AMOUNT: ")
			if _, err := fmt.Scanf("%d\n", &amount); err != nil {
				log.Println("error incorrect amount")
				continue
			}

			// Check if the transaction is in the correct form.
			if err := checkTransaction(choice, amount); err != nil {
				log.Println(err)
				continue
			}

			// Make the transaction object.
			transaction := api.Transaction{
				UserID: username,
				Type:   choice,
				Amount: amount,
			}

			// Make the transaction request to the server.
			if err := executeTransaction(client, transaction); err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("TRANSACTION COMPLETE")
		} else if choice == "b" {
			// Make the balance request to the server.
			balance, err := getBalance(client, username)
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("YOUR BALANCE IS:", balance)
		} else if choice == "e" {
			fmt.Println("BYE BYE")
			break
		}
		// Check if the user wants to continue.
		fmt.Print("WOULD YOU LIKE TO CONTINUE (Y/N): ")
		fmt.Scanf("%s\n", &answer)
	}
}

func userPrompt() {
	fmt.Println("1. W TO WITHDRAW AN AMOUNT")
	fmt.Println("2. D TO DEPOSIT AN AMOUNT")
	fmt.Println("3. B TO SEE YOUR BALANCE")
	fmt.Println("4. E TO EXIT")
	fmt.Print("PLEASE CHOOSE YOUR ACTION: ")
}

// checkTransaction checks if the transaction is in the correct form.
// The transaction type (tranType) has to be w/W or d/D.
// The amount has to be a multiple of 20 or 50.
func checkTransaction(tranType string, amount int) error {
	if !(tranType == "w" || tranType == "W" || tranType == "d" || tranType == "D") {
		return errors.New("the transaction has to be between w/W or d/D, try again")
	}

	if amount%20 != 0 && amount%50 != 0 {
		return errors.New("the amount has to be multiple of 20 or multiple 50, try again")
	}
	return nil
}

// executeTransaction function makes a RPC call to the server to execute
// a transaction in the ATM's database.
func executeTransaction(client *rpc.Client, transaction api.Transaction) error {
	var reply int
	err := client.Call("ATM.Transaction", &transaction, &reply)
	if err != nil {
		return err
	}
	return nil
}

// getBalance function makes a RPC call to the server requesting the balance
// of the specific user.
func getBalance(client *rpc.Client, username string) (int, error) {
	var balance int
	err := client.Call("ATM.Balance", username, &balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
