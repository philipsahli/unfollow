package twitter

import (
	"time"

	go_twitter "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type User struct {
	ID          int64 `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
	LastTweetAt time.Time
	LastLikeAt  time.Time
	ScreenName  string
}

type Follower struct {
	ID int64 `gorm:"primary_key"`
}

//var client *go_twitter.Client

func MakeClient(consumerKey *string, consumerSecret *string, accessToken *string, accessSecret *string) *go_twitter.Client {
	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := go_twitter.NewClient(httpClient)
	return client
}
