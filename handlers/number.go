package handlers

import (
	"bytes"
	"fmt"
	"awesomeProject/repos"
	"time"
)

type NumberHandler struct {
	repo repos.NumberRepo
}

func NewNumberHandler(repo repos.NumberRepo) *NumberHandler{
	h := new(NumberHandler)
	h.repo = repo
	return h
}

func (h *NumberHandler) Save(values chan string) {

	//buffer messages
	var buffer bytes.Buffer

	//statistics
	var newMessages = 0
	var duplicates = 0
	var uniques = make(map[string]string)

	//flush control, to save to our file
	var flush = false

	//monitor
	go func() {
		for{
			time.Sleep(10 * time.Second)
			fmt.Printf(
				"[%v] processed %v new unique numbers, %v duplicates. Unique total: %v \n",
				time.Now().Format("Mon Jan _2 2006 15:04:05 "),
					newMessages,
						duplicates,
							len(uniques))
			if newMessages > 0 {
				flush = true
			}
		}
	}()

	//read off the channel
	for msg := range values {

		_, exist := uniques[msg]
		if !exist {
			buffer.WriteString(msg+"\n")//@todo use server newline
			uniques[msg] = msg
		} else {
			duplicates++
		}

		if flush {
			var err = h.repo.Save(buffer.String())

			if isError(err) {
				//retry... yes it could be better :)
				h.repo.Save(buffer.String())
			}

			flush = false
			newMessages = 0
			duplicates = 0
		}

		newMessages++
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}