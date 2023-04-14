package main

import (
	"github.com/gorilla/mux"
	"io"
	"fmt"
	"log"
	"time"
	"crypto/md5"
	"crypto/sha256"
	"net/http"
	"encoding/json"
	"encoding/hex"
)

type Block struct {
	Pos          int
	Data         BookCheckout
	TimeStamp    string
	Hash         string
	PreviousHash string
}

type BookCheckout struct {
	BookID       string `json:"bookId"`
	UserID       string `json:"userId"`
	CheckOutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"isGenesis"`
}

type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishDate string `json:"publish_date"`
	ISBN        string `json:"isbn"`
}

type Blockchain struct {
	blocks []*Block
}

var Blockchain *Blockchain


func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.data)


	data := string(b.pos) + b.TimeStamp + string(bytes) + b.PreviousHash

     hash := sha256.New()
	hash.Write([]byte(data))

	b.Hash = hex.EncodeToString(hash.sum(nil))


}

func CreateBlock(prevBlock *Block, checkoutitem, BookCheckout ) {
	block := &Block{}
	block.Pos = prevBlock.Pos + 1
	block.TimeStamp = time.Now().String()
	block.PreviousHash = prevBlock.Hash()
	block.generateHash()
}

func (bc *Blockchain) AddBlock(data BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBLock(prevBlock, data)

	if validBlock(block, prevBlock) {
       bc.blocks = append(bc.blocks, block)
	}
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutitem BookCheckout 

	if err := json.NewDecoder(r.Body).Decode(&checkoutitem); err != nil {
		r.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not write block: %v", err)
		w.Write([]byte("could not write block"))
	}
	Blockchain.AddBlock(checkoutitem)
}


 func newBook(w http.ResponseWriter , r *http.Request) {
      var book Book

	 if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Could not create: %v", err)
		w.Write([]byte("could not create new book"))
		return
	 }

	 h := md5.New()
	 io.WriteString(h, book.ISBN+book.PublishDate)
	 book.ID = fmt.Sprintf("%x", h.Sum(nil))

	 resp, err := json.MarshalIndent(book, "", " ")
	 if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		w.Write([]byte("could not save Book data"))

		return
	 }

	 w.WriteHeader(http.StatusOK)
	 w.Write(resp)
 }

func main() {
	r := mux.NewRouter()
	r.Handler('/', getBlockchain).Methods("GET")
	r.Handler('/', writeBlock).Methods("POST")
	r.Handler('/', newBook).Methods("GET")

	log.Println("Listening on port 3000")

	log.Fatal(http.ListenAndServe(":3000", r))
}
