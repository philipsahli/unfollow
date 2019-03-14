package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
)

// func SynchronizeFollowers() {
// 	f, _, err := client.Followers.IDs(&twitter.FollowerIDParams{
// 		ScreenName: "philipsahli",
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	for _, f := range f.IDs {
// 		fm := Follower{ID: f}
// 		fm.Save()
// 	}
// }

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}

// SynchronizeFriend synchronizes friend to database
func SynchronizeFriend(id int64) {

	user := User{ID: id}

	// Save at the end
	defer user.Save()

	// Read User
	user.Get()
	then := time.Now().Add(time.Duration(-4) * time.Hour)
	if user.UpdatedAt.Before(then) {
		fmt.Print("r")
	} else {
		fmt.Print("n")
		return
	}

	// ScreenName from Twitter
	if user.ScreenName == "" {
		uo, resp, err := client.Users.Show(&twitter.UserShowParams{
			UserID: id,
		})

		user.ScreenName = uo.ScreenName
		printLimit(resp)
		if err != nil {
			log.Print(err)
		}

	}

	// Last Tweet
	// lt, resp, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
	// 	UserID: id,
	// 	Count:  1,
	// })

	// printLimit(resp)
	// if err != nil {
	// 	log.Print(err)
	// }

	// if len(lt) > 0 {
	// 	user.LastTweetAt, err = lt[0].CreatedAtTime()
	// 	if err != nil {
	// 		log.Print(err)
	// 	}
	// }

	// Last Like
	ll, _, err := client.Favorites.List(&twitter.FavoriteListParams{
		UserID: id,
		Count:  1,
	})
	if err != nil {
		log.Print(err)
	}
	if len(ll) > 0 {
		user.LastLikeAt, err = ll[0].CreatedAtTime()
		if err != nil {
			log.Print(err)
		}
	}

}

func printLimit(resp *http.Response) {
	fmt.Println(fmt.Sprintf("%s/%s until %s", resp.Header.Get("x-rate-limit-remaining"), resp.Header.Get("x-rate-limit-limit"), resp.Header.Get("x-rate-limit-reset")))
}

// Synchronize data to database
func Synchronize(client *twitter.Client) {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Print("Synchronize ")
	screenName := "philipsahli"
	ids, resp, err := client.Friends.IDs(&twitter.FriendIDParams{
		ScreenName: screenName,
	})
	if err != nil {
		fmt.Println(err)
	}
	printLimit(resp)

	var wg sync.WaitGroup
	// get all friends
	for _, id := range ids.IDs {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			SynchronizeFriend(id)
		}(id)
	}
	wg.Wait()
}

// Unfollow lists friends which should unfollowed
func Unfollow() {
	twoYears := time.Now().Add(time.Duration(0.5*-8760) * time.Hour)
	users := []User{}
	db.Where("last_tweet_at < ?", twoYears).Find(&users)
	for _, u := range users {
		fmt.Print("https://twitter.com/", u.ScreenName, " ")
	}
	fmt.Println("You could unfollow", len(users), "accounts")
}
