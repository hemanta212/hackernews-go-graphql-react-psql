package graph

import (
	"sync"

	"github.com/hemanta212/hackernews-go-graphql/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	// All posts since launching the graphql endpoint
	CreatedLink   *model.Link
	LinkObservers map[string]chan *model.Link
	CreatedVote   *model.Vote
	VoteObservers map[string]chan *model.Vote
	mu            sync.Mutex
}
