package get

import (
	"fmt"
	"sync"
	"time"
)

// Poll polls the main URL. It uses getPage to fetch each page.
// Poll sends orders, GET times, and errors into worker channels.
func (w *worker) Poll() {
	fmt.Printf("%-6d Spawn (%s) \n", w.ID(), w.Chain())
	// Wait for count to be set by PollCount()
	for !countSet {
		fmt.Printf("%-6d Count not set\n", w.ID())
		time.Sleep(1 * time.Second)
	}
	var wg sync.WaitGroup
	pages := getPages(count, w.limit)
	start := time.Now()
	for i := 1; i <= pages; i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()
			if err := getPage(
				w.chain,
				page,
				w.limit,
				w.orderChan,
			); err != nil {
				w.errorChan <- err
			}
		}(i)
	}
	wg.Wait()
	getTime := time.Since(start)
	fmt.Printf(
		"%6d GET time:\t%6v\n",
		w.ID(),
		getTime,
	)
	w.timeChan <- getTime
}
