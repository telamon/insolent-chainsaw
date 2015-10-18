package main

import (
	. "github.com/telamon/wharfmaster"
)

func main() {
	_, err := RegenerateConf()
	if err != nil {
		panic(err)
	}
	nginx, err := StartNginx()
	if err != nil {
		panic(err)
	}
	defer StopNginx()
	nginx.Wait()
}
