package model

import (
	"github.com/hemanta212/hackernews-go-graphql/internal/links"
	"github.com/hemanta212/hackernews-go-graphql/internal/users"
	"github.com/hemanta212/hackernews-go-graphql/internal/votes"
)

func FromLink(link *links.Link) *Link {
	return &Link{
		ID:          link.ID,
		Description: link.Description,
		URL:         link.Url,
		CreatedAt:   link.CreatedAt,
		PostedBy:    FromUser(link.PostedBy),
		Votes:       []*Vote{},
	}
}

func FromUser(user *users.User) *User {
	return &User{ID: user.ID, Username: user.Username, Email: user.Email}
}

func FromVote(vote *votes.Vote) *Vote {
	m_vote := &Vote{
		ID:   vote.ID,
		User: FromUser(vote.VotedBy),
		Link: FromLink(vote.Link),
	}
	m_vote.Link.Votes = []*Vote{}
	return m_vote
}
