package jobs

func Pipe() {
	chPocketId := make(chan string)

	go FetchPacketsList(chPocketId)

	for i := 0; i < 10; i ++ {
		go FetchBulkCacheDownload(chPocketId)
	}
}
