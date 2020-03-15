// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package worker // import "miniflux.app/worker"

import (
	"miniflux.app/logger"
	"miniflux.app/model"
	"miniflux.app/reader/feed"
)

// Worker refreshes a feed in the background.
type Worker struct {
	id          int
	feedHandler *feed.Handler
}

// Run wait for a job and refresh the given feed.
func (w *Worker) Run(c chan model.Job) {
	logger.Debug("[Worker] #%d started", w.id)

	for {
		job := <-c
		logger.Debug("[Worker #%d] got userID=%d, feedID=%d", w.id, job.UserID, job.FeedID)

		err := w.feedHandler.RefreshFeed(job.UserID, job.FeedID)
		if err != nil {
			logger.Error("[Worker] %v", err)
		}
	}
}
