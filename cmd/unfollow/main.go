package main

import (
	"flag"
	"log"
	"os"

	"github.com/coreos/pkg/flagutil"
	go_twitter "github.com/dghubble/go-twitter/twitter"
	"github.com/philipsahli/unfollow/pkg/app"
	"github.com/philipsahli/unfollow/pkg/twitter"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type persistance interface {
	Save()
}

func main() {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")
	web := flags.Bool("web", false, "start web")
	sync := flags.Bool("sync", false, "start sync")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	//dbName := "unfollow.db?cache=shared&mode=rwc"

	if *web {
		go app.Start()
	}

	if *sync {

		if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
			log.Fatal("Consumer key/secret and Access token/secret required")
		}

		// config := oauth1.NewConfig(*consumerKey, *consumerSecret)
		// token := oauth1.NewToken(*accessToken, *accessSecret)
		// // OAuth1 http.Client will automatically authorize Requests
		// httpClient := config.Client(oauth1.NoContext, token)

		// // Twitter client
		// client := twitter.NewClient(httpClient)
		client := twitter.MakeClient(consumerKey, consumerSecret, accessToken, accessSecret)

		// Verify Credentials
		verifyParams := &go_twitter.AccountVerifyParams{
			SkipStatus:   go_twitter.Bool(true),
			IncludeEmail: go_twitter.Bool(true),
		}
		client.Accounts.VerifyCredentials(verifyParams)
		// fmt.Printf("User's ACCOUNT:\n%+v\n", user)

		// // Home Timeline
		// homeTimelineParams := &twitter.HomeTimelineParams{
		// 	Count:     2,
		// 	TweetMode: "extended",
		// }
		// tweets, _, _ := client.Timelines.HomeTimeline(homeTimelineParams)
		// fmt.Printf("User's HOME TIMELINE:\n%+v\n", tweets)

		// // Mention Timeline
		// mentionTimelineParams := &twitter.MentionTimelineParams{
		// 	Count:     2,
		// 	TweetMode: "extended",
		// }
		// tweets, _, _ = client.Timelines.MentionTimeline(mentionTimelineParams)
		// fmt.Printf("User's MENTION TIMELINE:\n%+v\n", tweets)

		// // Retweets of Me Timeline
		// retweetTimelineParams := &twitter.RetweetsOfMeTimelineParams{
		// 	Count:     2,
		// 	TweetMode: "extended",
		// }
		// tweets, _, _ = client.Timelines.RetweetsOfMeTimeline(retweetTimelineParams)
		// fmt.Printf("User's 'RETWEETS OF ME' TIMELINE:\n%+v\n", tweets)

		// Update (POST!) Tweet (uncomment to run)
		// tweet, _, _ := client.Statuses.Update("just setting up my twttr", nil)
		// fmt.Printf("Posted Tweet\n%v\n", tweet)

		// screenName := "philipsahli"
		// ids, _, _ := client.Friends.IDs(&twitter.FriendIDParams{
		// 	ScreenName: screenName,
		// })
		// for _, id := range ids.IDs {
		// 	fmt.Println(id)
		// 	user := User{ID: id}
		// 	db.FirstOrCreate(&user, user)

		// 	uo, _, _ := client.Users.Show(&twitter.UserShowParams{
		// 		UserID: id,
		// 	})

		// 	user.ScreenName = uo.ScreenName
		// 	db.FirstOrCreate(&user, user)

		// }
		twitter.Synchronize(client)

		//tw.Unfollow()

		// friendsListParams := &twitter.FriendListParams{
		// 	ScreenName: screenName,
		// 	Cursor:     -1,
		// }
		// friends, _, err := client.Friends.List(friendsListParams)
		// if err != nil {
		// 	fmt.Println(err)
		// }

		// a := 0
		// // fmt.Println(friends)

		// for _, f := range friends.Users {
		// 	fmt.Println("User's Friends:", a, f.ScreenName)
		// 	a++
		// }
		// // fmt.Println("Hello")

		// for {
		// 	fmt.Println(friends.NextCursor)
		// 	friends, _, _ := client.Friends.List(&twitter.FriendListParams{
		// 		ScreenName: screenName,
		// 		Cursor:     friends.NextCursor,
		// 	})

		// 	for _, f := range friends.Users {
		// 		fmt.Println("User's Friends:", a, f.ScreenName)
		// 		a++
		// 	}
		// 	if friends.NextCursor == 0 {
		// 		fmt.Println("Last user")
		// 		os.Exit(0)
		// 	}
		// }
	}

}