package jobs

const WORKERS_NUM = 10

func Pipe() {
	chPocketId := make(chan string)

	go FetchPacketsList(chPocketId)

	for i := 0; i < WORKERS_NUM; i++ {
		go FetchBulkCacheDownload(chPocketId)
	}
}
