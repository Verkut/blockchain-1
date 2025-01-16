package main

import (
	"blockchain-1/blockchain"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

/** FIRST COMMENT */
/** SECOND COMMENT AND TAG */
/** !! */

/** SIX COMMENT */
/** SEVEN COMMENT */

type BlockchainServer struct {
	bc *blockchain.Blockchain
}

func (s *BlockchainServer) getBlockchain(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GET /blocks")
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(s.bc.Chain)
	if err != nil {
		return
	}
}

func (s *BlockchainServer) writeBlock(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Data string `json:"data"`
	}

	_ = json.NewDecoder(r.Body).Decode(&data)

	var timeNow = time.Now()

	var newBlock blockchain.Block = s.bc.GenerateBlock(data.Data)
	var wg sync.WaitGroup

	//for i := 0; i < 100; i++ {
	//	wg.Add(1)
	//	go func() {
	//		fmt.Printf("goroutine #%d start\n", i)
	//		var timeGoroutine = time.Now()
	//		var block = s.bc.GenerateBlock(data.Data)
	//		fmt.Printf("goroutine #%d; result: %s; duration: %s;\n", i, block.Hash, time.Since(timeGoroutine))
	//		if newBlock.Hash == "" {
	//			newBlock = block
	//		}
	//		wg.Done()
	//	}()
	//}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			fmt.Printf("goroutine #%d start\n", i)
			var timeGoroutine = time.Now()
			response, err := http.Get("https://baconipsum.com/api/?type=all-meat&paras=1000&start-with-lorem=1")
			if err != nil {
				return
			}
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					fmt.Println(err.Error())
				}
			}(response.Body)

			// Читаем данные из response.Body
			data, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Printf("goroutine #%d error reading response: %s\n", i, err.Error())
				return
			}

			fileName := fmt.Sprintf("goroutine-%d.json", i)

			if SaveToJsonFile(fileName, data) {
				fmt.Printf("goroutine #%d; status: %s; duration: %s; file: %s\n", i, response.Status, time.Since(timeGoroutine), fileName)
			} else {
				fmt.Printf("goroutine #%d error saving file: %s\n", i, fileName)
			}

			fmt.Printf("goroutine #%d; status: %s; duration: %s;\n", i, response.Status, time.Since(timeGoroutine))
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Printf("all time: %d\n", time.Since(timeNow))

	fmt.Printf("POST /mine: %+v\n", newBlock)

	s.bc.AddBlock(newBlock)

	err := json.NewEncoder(w).Encode(newBlock)
	if err != nil {
		return
	}
}

func main() {
	bc := blockchain.NewBlockchain()
	server := &BlockchainServer{bc: bc}

	r := mux.NewRouter()
	r.HandleFunc("/blocks", server.getBlockchain).Methods("GET")
	r.HandleFunc("/mine", server.writeBlock).Methods("POST")

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}

func SaveToJsonFile(fileName string, data []byte) bool {
	// Определяем директорию
	storageDir := "storage"

	// Проверяем и создаём директорию, если её нет
	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		if err := os.Mkdir(storageDir, os.ModePerm); err != nil {
			return false
		}
	}

	// Полный путь к файлу
	fullPath := filepath.Join(storageDir, fileName)

	// Создаём файл
	file, err := os.Create(fullPath)
	if err != nil {
		return false
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(file)

	// Записываем данные
	_, err = file.Write(data)
	if err != nil {
		return false
	}

	return true
}
