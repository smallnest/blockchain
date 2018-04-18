package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Server 提供区块链rpc服务，比如增加区块，查看区块等.
type Server struct {
	Addr       string
	server     *http.Server // rpc server
	Blockchain *Blockchain
}

// NewServer 创建一个新的blockchain服务器.
func NewServer(addr string, bc *Blockchain) *Server {
	return &Server{
		Addr:       addr,
		Blockchain: bc,
	}
}

// Serve 开启http rpc server.
func (s *Server) Serve() error {
	r := s.configRouter()
	ss := &http.Server{
		Addr:           s.Addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.server = ss

	return s.server.ListenAndServe()
}

func (s *Server) configRouter() http.Handler {
	r := httprouter.New()
	r.GET("/blocks", s.handleGetBlockchain)
	r.POST("/blocks", s.handleWriteBlock)
	return r
}

func (s *Server) handleGetBlockchain(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	startHeight := r.FormValue("start")
	start := 0
	var err error
	if startHeight != "" {
		start, err = strconv.Atoi(startHeight)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	bytes, err := json.MarshalIndent(s.Blockchain.Blocks[start:], "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.Write(bytes)
}

func (s *Server) handleWriteBlock(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	s.Blockchain.Lock()
	defer s.Blockchain.Unlock()

	prevBlock := s.Blockchain.Blocks[len(s.Blockchain.Blocks)-1]
	newBlock := s.Blockchain.generateBlock(prevBlock, data)

	if validateBlock(newBlock, prevBlock) {
		s.Blockchain.AddBlock(newBlock)
		respondJSON(w, r, http.StatusOK, newBlock)
		return
	}

	http.Error(w, "invalid new block", http.StatusInternalServerError)
}

func respondJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
