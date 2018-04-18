package main

import (
	"flag"

	"github.com/smallnest/blockchain"
	"github.com/smallnest/blockchain/store"
	"github.com/smallnest/log"
)

var (
	privateKey = flag.String("privateKey", "", "private key")
	addr       = flag.String("addr", ":8972", "listened address")
	dataFile   = flag.String("data", "./data", "data file")
)

func main() {
	flag.Parse()
	if *privateKey == "" {
		log.Info("请使用key命令行生成你自己的私钥，并且妥善保存。一旦丢失，无法找回!")
		return
	}

	store, err := store.NewLevelDBStore(*dataFile)
	if err != nil {
		log.Fatalf("failed to create leveldb store: %v", err)
	}
	defer store.Close()

	// 创建一个区块链
	var bc = &blockchain.Blockchain{
		Store:      store,
		Difficulty: 5,
		PrefixZero: "00000",
	}

	err = bc.LoadFromStore()
	if err != nil {
		log.Fatal(err)
	}

	if len(bc.Blocks) == 0 {
		bc.GenerateGenesisBlock()
	}

	// 创建 rpc server
	var server = blockchain.NewServer(*addr, bc)

	// 启动服务
	if err := server.Serve(); err != nil {
		log.Errorf("failed to serve: %v", err)
		return
	}

	log.Info("exit mormally")
}
