package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/antihax/optional"
	"github.com/kylebevans/twitapi"
)

// PrintTweet prints out a tweet followed by a little red line.
func PrintTweet(tweet twitapi.FilteredStreamingTweet) {
	fmt.Printf("%v\n\u001b[31m--------\n\u001b[0m", tweet.Data.Text)
}

// Grab tweet data from the last 7 days to seed the word cloud.
func ConvoControlTweets(ctx context.Context, s string, apiClient *twitapi.APIClient) {
	var searchOpts *twitapi.TweetsRecentSearchOpts
	searchOpts = new(twitapi.TweetsRecentSearchOpts)
	var tweets twitapi.TweetSearchResponse
	var nextToken optional.String
	var tweetFields optional.Interface
	var err error

	tweetFields = optional.NewInterface([]string{"reply_settings"})
	searchOpts.TweetFields = tweetFields

	// API is paginated, so print pages of recent tweets in the last week
	// that match the query until they are all done
	for ok := true; ok; ok = (tweets.Meta.NextToken != "") {

		tweets, _, err = apiClient.SearchApi.TweetsRecentSearch(ctx, s, searchOpts)

		if err != nil {
			log.Printf("Could not seed data: %v", err)
			return
		}

		if tweets.Data.ReplySettings == MENTIONED_USERS || tweets.Data.ReplySettings == FOLLOWING {
			for _, v := range tweets.Data {
				PrintTweet(v.Text, w)
			}
		}

		nextToken = optional.NewString(tweets.Meta.NextToken)
		searchOpts.NextToken = nextToken
		time.Sleep(2 * time.Second) // Twitter API is rate limited to 450 requests per 15 min.
	}
}

func main() {
	ctx := context.WithValue(context.Background(), twitapi.ContextAccessToken, os.Getenv("TWITTER_BEARER_TOKEN"))
	cfg := twitapi.NewConfiguration()
	apiClient := twitapi.NewAPIClient(cfg)
	searchFilter := "tech"

	ConvoControlTweets(ctx, searchFilter, apiClient)
}
