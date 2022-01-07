package twitter

import (
	"fmt"
	"log"
	"strings"
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
func SynchronizeFriend(workerId int, client *twitter.Client, ids <-chan int64, results chan<- int64, errors chan<- error) {

	for id := range ids {
		fmt.Println("worker", workerId, "started  job", id)
		if client == nil {
			log.Fatal("client get niled")
		}
		user := User{ID: id}

		// Save at the end
		defer user.Save()

		// Read User
		err := user.Get()
		if err != nil {
			errors <- err
		}
		then := time.Now().Add(time.Duration(-4) * time.Hour)
		if user.UpdatedAt.Before(then) {
			fmt.Print("r")
		} else {
			fmt.Print("n")
			return
		}

		time.Sleep(1 * time.Second)

		// ScreenName from Twitter
		if user.ScreenName == "" {
			uo, resp, err := client.Users.Show(&twitter.UserShowParams{
				UserID: id,
			})

			user.ScreenName = uo.ScreenName
			if err != nil {
				log.Fatal(err)
				fmt.Print("Limits:")
				fmt.Println(prettyLimit(resp))
				//limits, _, lerr := client.RateLimits.Status(&twitter.RateLimitParams{Resources: []string{"Users"}})
				log.Fatal(err)
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
		ll, resp, err := client.Favorites.List(&twitter.FavoriteListParams{
			UserID: id,
			Count:  1,
		})
		if err != nil {
			if strings.Contains(err.Error(), "88") {
				log.Print(err, prettyLimit(resp))
			}
			errors <- err
			log.Print(err)
			return
		}
		if len(ll) > 0 {
			user.LastLikeAt, err = ll[0].CreatedAtTime()
			if err != nil {
				log.Print(err)
			}
		}

		fmt.Println("worker", workerId, "ended job", id)
		results <- id
	}

}

// Synchronize data to database
func Synchronize(client *twitter.Client) error {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Print("Synchronize ")
	screenName := "philipsahli"

	// Get IDs for all friend
	ids, resp, err := client.Friends.IDs(&twitter.FriendIDParams{
		ScreenName: screenName,
	})
	log.Print(prettyLimit(resp))
	if err != nil {
		fmt.Println(prettyLimit(resp))
		log.Fatal("Cannot get Friends.IDs:", err)
	}
	//printLimit(resp)

	numJobs := len(ids.IDs)
	jobs := make(chan int64, numJobs)
	results := make(chan int64, numJobs)
	errors := make(chan error)

	// Start n workers
	for w := 1; w <= 5; w++ {
		go SynchronizeFriend(w, client, jobs, results, errors)
	}

	// Hand-over work
	idsWork := ids.IDs
	for i, id := range idsWork {
		fmt.Printf("Kickoff %d/%d to sync (id=%d)\n", i+1, len(idsWork), id)
		jobs <- id
		/*
			go func(id int64) {
				defer wg.Done()
				SynchronizeFriend(id, client)
			}(id)
		*/

	}
	time.Sleep(1 * time.Second)
	close(results)

	// Get results
	for result := range results {
		fmt.Println("Result:", result)
	}

	return nil
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
