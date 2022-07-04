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

func Perform(args Arguments, writer io.Writer) error {
	if len(args["operation"]) == 0 {
		return errors.New("`-operation` flag has to be specified")
	}
	if len(args["fileName"]) == 0 {
		return errors.New("`-fileName` flag has to be specified")
	}

	jsonFile, err := os.OpenFile(args["fileName"], os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {

		}
	}(jsonFile)

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var users []map[string]string
	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		return err
	}

	switch args["operation"] {
	case "add":
		if len(args["item"]) == 0 {
			return errors.New("`-item` flag has to be specified")
		}

		var item map[string]string
		err := json.Unmarshal([]byte(args["item"]), &item)
		if err != nil {
			return err
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
	case "list":
		_, err := writer.Write(byteValue)
		if err != nil {
			return err
		}
	case "findById":
		// If `-id` flag is not provided error «-id flag has to be specified» should be returned from Perform function
		if len(args["id"]) == 0 {
			return errors.New("`-id` flag has to be specified")
		}

		// If user exists, then json object should be written in `io.Writer`
		for _, value := range users {
			if value["id"] == args["id"] {
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
		_, err := writer.Write([]byte(""))
		if err != nil {
			return err
		}
	case "remove":
		// If `-id` flag is not provided error «-id flag has to be specified» should be returned from Perform function
		if len(args["id"]) == 0 {
			return errors.New("`-id` flag has to be specified")
		}

		// If user exists, then json object should be written in `io.Writer`
		for key, value := range users {
			if value["id"] == args["id"] {
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

		// If user with id `"2"`, for example, does not exist, Perform functions should print message to the `io.Writer` «Item with id 2 not found».
		return fmt.Errorf("item with id '%s' not found", args["id"])

	}
	return nil
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
