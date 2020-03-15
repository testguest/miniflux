// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package client // import "miniflux.app/client"

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
)

// Client holds API procedure calls.
type Client struct {
	request *request
}

// Me returns the logged user information.
func (c *Client) Me() (*User, error) {
	body, err := c.request.Get("/v1/me")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var user *User
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&user); err != nil {
		return nil, fmt.Errorf("miniflux: json error (%v)", err)
	}

	return user, nil
}

// Users returns all users.
func (c *Client) Users() (Users, error) {
	body, err := c.request.Get("/v1/users")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var users Users
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return users, nil
}

// UserByID returns a single user.
func (c *Client) UserByID(userID int64) (*User, error) {
	body, err := c.request.Get(fmt.Sprintf("/v1/users/%d", userID))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var user User
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&user); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return &user, nil
}

// UserByUsername returns a single user.
func (c *Client) UserByUsername(username string) (*User, error) {
	body, err := c.request.Get(fmt.Sprintf("/v1/users/%s", username))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var user User
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&user); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return &user, nil
}

// CreateUser creates a new user in the system.
func (c *Client) CreateUser(username, password string, isAdmin bool) (*User, error) {
	body, err := c.request.Post("/v1/users", &User{Username: username, Password: password, IsAdmin: isAdmin})
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var user *User
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&user); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return user, nil
}

// UpdateUser updates a user in the system.
func (c *Client) UpdateUser(userID int64, userChanges *UserModification) (*User, error) {
	body, err := c.request.Put(fmt.Sprintf("/v1/users/%d", userID), userChanges)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var u *User
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&u); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return u, nil
}

// DeleteUser removes a user from the system.
func (c *Client) DeleteUser(userID int64) error {
	body, err := c.request.Delete(fmt.Sprintf("/v1/users/%d", userID))
	if err != nil {
		return err
	}
	body.Close()
	return nil
}

// Discover try to find subscriptions from a website.
func (c *Client) Discover(url string) (Subscriptions, error) {
	body, err := c.request.Post("/v1/discover", map[string]string{"url": url})
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var subscriptions Subscriptions
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&subscriptions); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return subscriptions, nil
}

// Categories gets the list of categories.
func (c *Client) Categories() (Categories, error) {
	body, err := c.request.Get("/v1/categories")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var categories Categories
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&categories); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return categories, nil
}

// CreateCategory creates a new category.
func (c *Client) CreateCategory(title string) (*Category, error) {
	body, err := c.request.Post("/v1/categories", map[string]interface{}{
		"title": title,
	})

	if err != nil {
		return nil, err
	}
	defer body.Close()

	var category *Category
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&category); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return category, nil
}

// UpdateCategory updates a category.
func (c *Client) UpdateCategory(categoryID int64, title string) (*Category, error) {
	body, err := c.request.Put(fmt.Sprintf("/v1/categories/%d", categoryID), map[string]interface{}{
		"title": title,
	})

	if err != nil {
		return nil, err
	}
	defer body.Close()

	var category *Category
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&category); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return category, nil
}

// DeleteCategory removes a category.
func (c *Client) DeleteCategory(categoryID int64) error {
	body, err := c.request.Delete(fmt.Sprintf("/v1/categories/%d", categoryID))
	if err != nil {
		return err
	}
	defer body.Close()

	return nil
}

// Feeds gets all feeds.
func (c *Client) Feeds() (Feeds, error) {
	body, err := c.request.Get("/v1/feeds")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var feeds Feeds
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&feeds); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return feeds, nil
}

// Export creates OPML file.
func (c *Client) Export() ([]byte, error) {
	body, err := c.request.Get("/v1/export")
	if err != nil {
		return nil, err
	}
	defer body.Close()

	opml, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	return opml, nil
}

// Import imports an OPML file.
func (c *Client) Import(f io.ReadCloser) error {
	_, err := c.request.PostFile("/v1/import", f)
	return err
}

// Feed gets a feed.
func (c *Client) Feed(feedID int64) (*Feed, error) {
	body, err := c.request.Get(fmt.Sprintf("/v1/feeds/%d", feedID))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var feed *Feed
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&feed); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return feed, nil
}

// CreateFeed creates a new feed.
func (c *Client) CreateFeed(url string, categoryID int64) (int64, error) {
	body, err := c.request.Post("/v1/feeds", map[string]interface{}{
		"feed_url":    url,
		"category_id": categoryID,
	})
	if err != nil {
		return 0, err
	}
	defer body.Close()

	type result struct {
		FeedID int64 `json:"feed_id"`
	}

	var r result
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&r); err != nil {
		return 0, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return r.FeedID, nil
}

// UpdateFeed updates a feed.
func (c *Client) UpdateFeed(feedID int64, feedChanges *FeedModification) (*Feed, error) {
	body, err := c.request.Put(fmt.Sprintf("/v1/feeds/%d", feedID), feedChanges)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var f *Feed
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&f); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return f, nil
}

// RefreshFeed refresh a feed.
func (c *Client) RefreshFeed(feedID int64) error {
	body, err := c.request.Put(fmt.Sprintf("/v1/feeds/%d/refresh", feedID), nil)
	if err != nil {
		return err
	}
	body.Close()
	return nil
}

// DeleteFeed removes a feed.
func (c *Client) DeleteFeed(feedID int64) error {
	body, err := c.request.Delete(fmt.Sprintf("/v1/feeds/%d", feedID))
	if err != nil {
		return err
	}
	body.Close()
	return nil
}

// FeedIcon gets a feed icon.
func (c *Client) FeedIcon(feedID int64) (*FeedIcon, error) {
	body, err := c.request.Get(fmt.Sprintf("/v1/feeds/%d/icon", feedID))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var feedIcon *FeedIcon
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&feedIcon); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return feedIcon, nil
}

// FeedEntry gets a single feed entry.
func (c *Client) FeedEntry(feedID, entryID int64) (*Entry, error) {
	body, err := c.request.Get(fmt.Sprintf("/v1/feeds/%d/entries/%d", feedID, entryID))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var entry *Entry
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&entry); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return entry, nil
}

// Entry gets a single entry.
func (c *Client) Entry(entryID int64) (*Entry, error) {
	body, err := c.request.Get(fmt.Sprintf("/v1/entries/%d", entryID))
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var entry *Entry
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&entry); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return entry, nil
}

// Entries fetch entries.
func (c *Client) Entries(filter *Filter) (*EntryResultSet, error) {
	path := buildFilterQueryString("/v1/entries", filter)

	body, err := c.request.Get(path)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var result EntryResultSet
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return &result, nil
}

// FeedEntries fetch feed entries.
func (c *Client) FeedEntries(feedID int64, filter *Filter) (*EntryResultSet, error) {
	path := buildFilterQueryString(fmt.Sprintf("/v1/feeds/%d/entries", feedID), filter)

	body, err := c.request.Get(path)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	var result EntryResultSet
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("miniflux: response error (%v)", err)
	}

	return &result, nil
}

// UpdateEntries updates the status of a list of entries.
func (c *Client) UpdateEntries(entryIDs []int64, status string) error {
	type payload struct {
		EntryIDs []int64 `json:"entry_ids"`
		Status   string  `json:"status"`
	}

	body, err := c.request.Put("/v1/entries", &payload{EntryIDs: entryIDs, Status: status})
	if err != nil {
		return err
	}
	body.Close()

	return nil
}

// ToggleBookmark toggles entry bookmark value.
func (c *Client) ToggleBookmark(entryID int64) error {
	body, err := c.request.Put(fmt.Sprintf("/v1/entries/%d/bookmark", entryID), nil)
	if err != nil {
		return err
	}
	body.Close()

	return nil
}

// New returns a new Miniflux client.
func New(endpoint, username, password string) *Client {
	return &Client{request: &request{endpoint: endpoint, username: username, password: password}}
}

func buildFilterQueryString(path string, filter *Filter) string {
	if filter != nil {
		values := url.Values{}

		if filter.Status != "" {
			values.Set("status", filter.Status)
		}

		if filter.Direction != "" {
			values.Set("direction", filter.Direction)
		}

		if filter.Order != "" {
			values.Set("order", filter.Order)
		}

		if filter.Limit >= 0 {
			values.Set("limit", strconv.Itoa(filter.Limit))
		}

		if filter.Offset >= 0 {
			values.Set("offset", strconv.Itoa(filter.Offset))
		}

		if filter.After > 0 {
			values.Set("after", strconv.FormatInt(filter.After, 10))
		}

		if filter.AfterEntryID > 0 {
			values.Set("after_entry_id", strconv.FormatInt(filter.AfterEntryID, 10))
		}

		if filter.Before > 0 {
			values.Set("before", strconv.FormatInt(filter.Before, 10))
		}

		if filter.BeforeEntryID > 0 {
			values.Set("before_entry_id", strconv.FormatInt(filter.BeforeEntryID, 10))
		}

		if filter.Starred {
			values.Set("starred", "1")
		}

		if filter.Search != "" {
			values.Set("search", filter.Search)
		}

		path = fmt.Sprintf("%s?%s", path, values.Encode())
	}

	return path
}
