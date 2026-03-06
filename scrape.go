package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jacobhuneke/gator/internal/database"
	"github.com/jacobhuneke/gator/internal/rss"
)

func scrapeFeeds(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	feed, err := s.db.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	rss.CleanFeed(rssFeed)

	for _, item := range rssFeed.Channel.Items {

		parsedTime, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			return err
		}

		newPostParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: parsedTime,
			FeedID:      feed.ID,
		}
		_, err = s.db.CreatePost(context.Background(), newPostParams)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	return nil
}
