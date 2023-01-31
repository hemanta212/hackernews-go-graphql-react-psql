package votes

import (
	"log"

	"github.com/hemanta212/hackernews-go-graphql/internal/links"
	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/mysql"
	"github.com/hemanta212/hackernews-go-graphql/internal/users"
)

type Vote struct {
	ID      string
	Link    *links.Link
	VotedBy *users.User
}

func (vote Vote) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Votes(LinkID, UserID) VALUES(?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(vote.Link.ID, vote.VotedBy.ID)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	return id
}

func GetVotesByLinkId(linkID string) ([]*Vote, error) {
	stmt, err := database.Db.Prepare(`SELECT V.ID, L.Description, L.Url, L.CreatedAt,
                                                 LU.ID, LU.username, LU.email,
                                                 U.ID, U.username, U.email
                                                 from Votes V inner join Users U on V.UserID = U.ID
                                                 inner join Links L on V.LinkID=L.ID
                                                 inner join Users LU on L.UserID=LU.ID
                                                 WHERE L.ID=?`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(linkID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var votes []*Vote
	for rows.Next() {
		vote, link, linkAuthor, VotedBy := &Vote{}, &links.Link{ID: linkID}, &users.User{}, &users.User{}
		err := rows.Scan(&vote.ID, &link.Description, &link.Url, &link.CreatedAt,
			&linkAuthor.ID, &linkAuthor.Username, &linkAuthor.Email,
			&VotedBy.ID, &VotedBy.Username, &VotedBy.Email)
		if err != nil {
			log.Fatal(err)
		}
		link.PostedBy = linkAuthor
		vote.VotedBy = VotedBy
		vote.Link = link
		votes = append(votes, vote)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return votes, nil
}
