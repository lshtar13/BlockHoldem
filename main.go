package main

func main() {
	bc, _ := NewBlockchain()
	defer bc.db.Close()

	cli := CLI{bc}
	cli.Run()
}
