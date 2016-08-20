package main

import "bytes"

func subscribeFeed(feed string, user int) (string, error) {
	var b bytes.Buffer

	feedName, err := getFeedName(feed)
	if err != nil {
		return "", err
	}

	// Subscribe the user to the feed
	err = newSubscription(feed, user)
	if err != nil {
		return "", err
	}

	if err = templates.ExecuteTemplate(&b, "subscribe.tpl", feedName); err != nil {
		return "", err
	}

	return b.String(), nil
}

func unsubscribeFeed(feed string, user int) (string, error) {
	var b bytes.Buffer

	feedName, err := getFeedName(feed)
	if err != nil {
		return "", err
	}

	// Unsubscribe the user
	err = removeSubscription(feed, user)
	if err != nil {
		return "", err
	}

	if err = templates.ExecuteTemplate(&b, "unsubscribe.tpl", feedName); err != nil {
		return "", err
	}

	return b.String(), nil
}
