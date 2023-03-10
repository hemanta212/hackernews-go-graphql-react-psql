package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.24

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hemanta212/hackernews-go-graphql/graph/model"
	"github.com/hemanta212/hackernews-go-graphql/internal/auth"
	"github.com/hemanta212/hackernews-go-graphql/internal/links"
	"github.com/hemanta212/hackernews-go-graphql/internal/users"
	"github.com/hemanta212/hackernews-go-graphql/internal/utils"
	"github.com/hemanta212/hackernews-go-graphql/internal/votes"
	"github.com/hemanta212/hackernews-go-graphql/pkg/jwt"
)

// Post is the resolver for the post field.
func (r *mutationResolver) Post(ctx context.Context, input model.NewLink) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("Acces denied")
	}

	link := &links.Link{
		Description: input.Description,
		Url:         input.URL,
		PostedBy:    user,
	}
	link.Save()

	mLink := model.FromLink(link)
	r.CreatedLink = mLink
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, observer := range r.LinkObservers {
		fmt.Println("Sending link to ", k)
		fmt.Printf("Object %v\n", r.CreatedLink)
		observer <- r.CreatedLink
	}

	return model.FromLink(link), nil
}

// Signup is the resolver for the signup field.
func (r *mutationResolver) Signup(ctx context.Context, input model.NewUser) (*model.AuthPayload, error) {
	var payload model.AuthPayload
	user := &users.User{
		Username: input.Username,
		Password: input.Password,
		Email:    input.Email,
	}

	_, err := user.Save()
	if err != nil {
		return &payload, err
	}

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return &payload, err
	}

	return &model.AuthPayload{
		Token: &token,
		User:  model.FromUser(user),
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.Login) (*model.AuthPayload, error) {
	var payload model.AuthPayload
	user, err := users.GetUserByUsername(input.Username)
	if err != nil {
		return &payload, &users.WrongUsernameOrPasswordError{}
	}
	if !user.Authenticate(input.Password) {
		return &payload, &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return &payload, err
	}
	return &model.AuthPayload{
		Token: &token,
		User:  model.FromUser(user),
	}, nil
}

// Vote is the resolver for the vote field.
func (r *mutationResolver) Vote(ctx context.Context, linkID string) (*model.Vote, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Acces denied")
	}

	linkId, err := strconv.Atoi(linkID)
	if err != nil {
		return nil, err
	}
	link, err := links.GetLinkByID(linkId)

	vote := &votes.Vote{
		Link:    link,
		VotedBy: user,
	}
	vote.Save()

	mVote := model.FromVote(vote)
	r.CreatedVote = mVote
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, observer := range r.VoteObservers {
		fmt.Println("Sending vote to ", k)
		observer <- r.CreatedVote
	}

	return mVote, nil
}

// RefreshToken is the resolver for the refreshToken field.
func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access is denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

// Feed is the resolver for the feed field.
func (r *queryResolver) Feed(ctx context.Context, filter *string, offset *int, limit *int, orderBy *model.LinkOrderByInput) (*model.Feed, error) {
	linkMod := &links.LinkMod{
		Filter:  "%%",
		Offset:  0,
		Limit:   15,
		OrderBy: &links.LinkOrderByInput{},
	}
	if filter != nil {
		linkMod.Filter = fmt.Sprintf("%%%s%%", *filter)
	}
	if offset != nil {
		linkMod.Offset = *offset
	}
	if limit != nil {
		linkMod.Limit = *limit
	}
	if orderBy != nil && orderBy.CreatedAt != nil {
		linkMod.OrderBy = &links.LinkOrderByInput{
			CreatedAt: links.Sort(*orderBy.CreatedAt),
		}
	} else if orderBy != nil && orderBy.Description != nil {
		linkMod.OrderBy = &links.LinkOrderByInput{
			Description: links.Sort(*orderBy.Description),
		}
	}

	var dbLinks []*links.Link
	dbLinks = links.GetAll(linkMod)
	// perform conversion to model.Link
	var resultLinks []*model.Link
	for _, link := range dbLinks {
		custom_votes, err := votes.GetVotesByLinkId(link.ID)
		if err != nil {
			log.Fatal(err)
		}
		model_votes := []*model.Vote{}
		for _, vote := range custom_votes {
			model_votes = append(model_votes, model.FromVote(vote))
		}
		modelLink := model.FromLink(link)
		modelLink.Votes = model_votes

		resultLinks = append(resultLinks, modelLink)
	}

	feedId := fmt.Sprintf(`filter-%s-limit-%d-offset-%d-createdat-%s-description-%s}`,
		linkMod.Filter, linkMod.Limit, linkMod.Offset, linkMod.OrderBy.CreatedAt, linkMod.OrderBy.Description)
	return &model.Feed{ID: "main-feed::" + feedId, Links: resultLinks, Count: len(resultLinks)}, nil
}

// NewLink is the resolver for the newLink field.
func (r *subscriptionResolver) NewLink(ctx context.Context) (<-chan *model.Link, error) {
	if r.LinkObservers == nil {
		fmt.Println("Nill observers map creating")
		r.LinkObservers = make(map[string]chan *model.Link)
	}
	// create id and channel for each active subscription, we'll push changes into this channel
	// when a new subscription is created by the client, this resolver fires first
	id := utils.RandString(8)
	newLink := make(chan *model.Link, 1)
	fmt.Println("Subbed by", id)

	// start a goroutine to allow or cleaning up subscriptions that are disconnected.
	// this go routine will only get past Done() when a client terminates the subscription.
	// Allows us to only remove the ref from list of chat observers since its no longer needed
	go func() {
		<-ctx.Done()
		fmt.Println("Conn closed removing ", id)
		r.mu.Lock()
		delete(r.LinkObservers, id)
		r.mu.Unlock()
	}()
	r.mu.Lock()
	// keep a reference of the channel so that we can push changes into it when new messages are posted.
	r.LinkObservers[id] = newLink
	r.mu.Unlock()
	return newLink, nil
}

// NewVote is the resolver for the newVote field.
func (r *subscriptionResolver) NewVote(ctx context.Context) (<-chan *model.Vote, error) {
	if r.VoteObservers == nil {
		r.VoteObservers = make(map[string]chan *model.Vote)
	}
	fmt.Println("sub entered")
	id := utils.RandString(8)
	fmt.Println(id)
	newVote := make(chan *model.Vote, 1)
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		delete(r.VoteObservers, id)
		r.mu.Unlock()
	}()

	r.mu.Lock()
	r.VoteObservers[id] = newVote
	fmt.Println("Added channel to observers")
	r.mu.Unlock()
	return newVote, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
