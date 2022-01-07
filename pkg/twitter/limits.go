package twitter

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type rateLimit struct {
	remaining int
	limit     int
	reset     int
}

var limits rateLimit

func getLimit(resp *http.Response) (int, int, int) {
	// x-rate-limit-limit: the rate limit ceiling for that given endpoint
	// x-rate-limit-remaining: the number of requests left for the 15-minute window
	// x-rate-limit-reset: the remaining window before the rate limit resets, in UTC epoch seconds
	remaining, _ := strconv.Atoi(resp.Header.Get("x-rate-limit-remaining"))
	limit, _ := strconv.Atoi(resp.Header.Get("x-rate-limit-limit"))
	reset, _ := strconv.Atoi(resp.Header.Get("x-rate-limit-reset"))

	return remaining, limit, reset
}

func prettyLimit(resp *http.Response) string {
	remaining, limit, reset := getLimit(resp)

	resetTs := time.Unix(int64(reset), 0)

	//now := time.Now()
	//zone, _:= now.Zone()
	loc, err := time.LoadLocation("Europe/Zurich")
	if err != nil {
		log.Fatal("time.LoadLocation: ", err)
	}
	//	utc := resetTs.In(loc)

	//utc.In(a.Local())
	nowUnix := time.Now().Unix()
	log.Print(int64(reset) - nowUnix)

	return fmt.Sprintf("remaining=%d limit=%d reset=%s", remaining, limit, resetTs.UTC().In(loc))

}
