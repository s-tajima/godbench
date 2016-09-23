package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jessevdk/go-flags"
	"log"
	"os"
	"sync"
	"time"
)

var (
	db    *sql.DB
	wg    *sync.WaitGroup
	rw    *sync.RWMutex
	stats *Stats
)

var options struct {
	BenchmarkFile string `short:"b" long:"bench" required:"true" value-name:"<bench>" description:"benchmark file name."`
	//Envfile  string `short:"e" long:"envfile" required:"false" default:".env" value-name:"<path/to/envfile>" description:"env file path."`
	Debug bool `short:"d" long:"debug"`
}

func main() {
	p := flags.NewParser(&options, flags.Default)

	_, err := p.ParseArgs(os.Args[1:])
	errorIf(err, false)

	fp, err := os.Open(options.BenchmarkFile)
	errorIf(err, true)

	db, err = sql.Open("mysql", "bench:bench@tcp(localhost:3306)/bench")
	errorIf(err, true)

	db.SetMaxOpenConns(10)

	stats = &Stats{}
	rw = &sync.RWMutex{}

	wg = &sync.WaitGroup{}
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		wg.Add(1)
		go query(scanner.Text())
	}
	wg.Wait()

	c := stats.Count()
	s := stats.Sum()

	log.Printf("count: %d, sum: %f", c, s)

}

func query(sql string) {
	defer wg.Done()
	defer rw.Unlock()

	start := time.Now()

	id := 0

	rows, err := db.Query(sql)
	errorIf(err, true)

	rows.Next()
	rows.Scan(&id)
	rows.Close()

	duration := time.Since(start)

	rw.Lock()
	stats.AddTime(duration.Seconds())
}

func errorIf(err error, echoError bool) {
	if err == nil {
		return
	}

	if echoError {
		fmt.Println(err)
	}
	os.Exit(1)
}
