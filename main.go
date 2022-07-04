package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}

	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch args["operation"] {
	case "add":
		return add(args, writer)
	case "list":
		return list(args, writer)
	case "findById":
		return findById(args, writer)
	case "remove":
		return remove(args, writer)
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	id := flag.String("id", "", "")
	operation := flag.String("operation", "", "")
	item := flag.String("item", "", "")
	fileName := flag.String("fileName", "", "")
	flag.Parse()
	return Arguments{
		"id":        *id,
		"operation": *operation,
		"item":      *item,
		"fileName":  *fileName,
	}
}

func add(args Arguments, writer io.Writer) error {
	if args["item"] == "" {
		return errors.New("-item flag has to be specified")
	}

	//////////////////////////////////////
	jsonFile, err := os.OpenFile(args["fileName"], os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var users []User
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &users)
		if err != nil {
			return err
		}
	}
	//////////////////////////////////////

	var item User
	err = json.Unmarshal([]byte(args["item"]), &item)
	if err != nil {
		return err
	}
	for _, value := range users {
		if value.Id == item.Id {
			errMessage := fmt.Sprintf("Item with id %s already exists", item.Id)
			_, err = writer.Write([]byte(errMessage))
			if err != nil {
				return err
			}
		}
	}

	users = append(users, item)

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		return err
	}

	err = jsonFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = jsonFile.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = jsonFile.Write(jsonUsers)
	if err != nil {
		return err
	}

	return nil
}

func remove(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}

	//////////////////////////////////////
	jsonFile, err := os.OpenFile(args["fileName"], os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var users []User
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &users)
		if err != nil {
			return err
		}
	}
	//////////////////////////////////////

	// If user exists, then json object should be written in `io.Writer`
	for key, value := range users {
		fmt.Println(value.Id)
		if value.Id == args["id"] {
			users = append(users[:key], users[key+1:]...)
			jsonUsers, err := json.Marshal(value)
			if err != nil {
				return err
			}
			_, err = writer.Write(jsonUsers)
			if err != nil {
				return err
			}
			return nil
		}
	}

	errMessage := fmt.Sprintf("Item with id %s not found", args["id"])
	_, err = writer.Write([]byte(errMessage))
	if err != nil {
		return err
	}
	return nil
}

func findById(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}

	//////////////////////////////////////
	jsonFile, err := os.OpenFile(args["fileName"], os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var users []User
	if len(byteValue) > 0 {
		err = json.Unmarshal(byteValue, &users)
		if err != nil {
			return err
		}
	}
	//////////////////////////////////////

	// If user exists, then json object should be written in `io.Writer`
	for _, value := range users {
		if value.Id == args["id"] {
			jsonUsers, err := json.Marshal(value)
			if err != nil {
				return err
			}
			_, err = writer.Write(jsonUsers)
			if err != nil {
				return err
			}
			return nil
		}
	}

	// If user with specified id does not exist in the users.json file, then empty string has to be written to  the `io.Writer`
	_, err = writer.Write([]byte(""))
	if err != nil {
		return err
	}

	return nil
}

func list(args Arguments, writer io.Writer) error {
	//////////////////////////////////////
	jsonFile, err := os.OpenFile(args["fileName"], os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	/*
		var users []User
		if len(byteValue) > 0 {
			err = json.Unmarshal(byteValue, &users)
			if err != nil {
				return err
			}
		}
	*/
	//////////////////////////////////////

	_, err = writer.Write(byteValue)
	if err != nil {
		return err
	}
	return nil
}
