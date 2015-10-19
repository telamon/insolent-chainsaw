package main

import (
	. "github.com/telamon/wharfmaster"
)

func main() {
	New().Start(5000)
	/*_, err := RegenerateConf()
	if err != nil {
		panic(err)
	}
	nginx, err := StartNginx()
	if err != nil {
		panic(err)
	}
	defer StopNginx()
	nginx.Wait()*/
}
