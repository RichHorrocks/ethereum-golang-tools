package main

/*
 * Installation:
 *  go get github.com/syndtr/goleveldb/leveldb
 *  go get github.com/syndtr/goleveldb/leveldb/errors
 *  go get github.com/ethereum/go-ethereum/rlp
 *  go get github.com/syndtr/goleveldb/leveldb/opt
 */


import (
//    "errors"
    "fmt"
    "log"
//    "os"
//    "os/user"

    "github.com/syndtr/goleveldb/leveldb"
    "github.com/syndtr/goleveldb/leveldb/opt"
//    "github.com/syndtr/goleveldb/leveldb/errors"
//    "github.com/ethereum/go-ethereum/rlp"
//    "github.com/ethereum/go-ethereum/common"
)

type LDBDatabase struct {
    db   *leveldb.DB
    comp  bool
}

func getLDBDatabase(file string) (*LDBDatabase, error) {
    // Open the db
    db, err := leveldb.OpenFile(file, 
                                &opt.Options{ErrorIfMissing: true})
    if err != nil {
        log.Println(err)
        return nil, err
    }

    database := &LDBDatabase{db: db, comp: false}

    return database, nil
}

func main() {
    // Open the database.
    db, err = getLDBDatabase("chaindata")
    if err != nil {
        fmt.Println(err)
    }

    // 
}
