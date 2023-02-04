package links

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/postgresql"
	"github.com/hemanta212/hackernews-go-graphql/internal/users"
)

type Link struct {
	ID          string
	Description string
	Url         string
	PostedBy    *users.User
	CreatedAt   time.Time
}

type LinkMod struct {
	Filter  string
	Limit   int
	Offset  int
	OrderBy *LinkOrderByInput
}

type LinkOrderByInput struct {
	CreatedAt   Sort `json:"createdAt"`
	Description Sort `json:"description"`
}

type Sort string

const (
	SortAsc  Sort = "asc"
	SortDesc Sort = "desc"
)

func (link *Link) Save() int {
	var lastInsertId int
	var createdAt time.Time
	err := database.Db.QueryRow("INSERT INTO Links(Description, Url, UserID) VALUES($1,$2,$3) returning id,CreatedAt;",
		link.Description, link.Url, link.PostedBy.ID,
	).Scan(&lastInsertId, &createdAt)

	link.CreatedAt = createdAt
	link.ID = strconv.Itoa(lastInsertId)

	if err != nil {
		log.Fatal(err)
	}
	return lastInsertId
}

func GetLinkByID(id int) (*Link, error) {
	stmt, err := database.Db.Prepare(
		"SELECT L.Description, L.Url, L.CreatedAt, L.UserID, U.Username, U.Email from Links L inner join Users U on L.UserID = U.ID WHERE L.ID=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	link := &Link{ID: strconv.Itoa(id)}
	user := &users.User{}
	err = stmt.QueryRow(id).Scan(&link.Description, &link.Url, &link.CreatedAt, &user.ID, &user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		} else {
			log.Fatal(err)
		}
	}
	link.PostedBy = user

	return link, nil
}

func GetAll(mods *LinkMod) []*Link {
	query := fmt.Sprintf(`SELECT
                       L.ID, L.Description, L.Url, L.UserID, L.CreatedAt,
                       U.Username, U.Email
                       FROM Links L
                       INNER JOIN Users U ON L.UserID = U.ID
                       WHERE L.Description LIKE $1
                       OR L.Url LIKE $1
                       OR to_char(L.CreatedAt, 'YYYY') LIKE $1
                       ORDER BY %s
                       LIMIT $2
                       OFFSET $3`, sanitizedOrderByClause(mods.OrderBy))

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(mods.Filter, mods.Limit, mods.Offset)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var links []*Link
	for rows.Next() {
		link, user := &Link{}, &users.User{}
		err := rows.Scan(&link.ID, &link.Description, &link.Url, &user.ID, &link.CreatedAt, &user.Username, &user.Email)
		if err != nil {
			log.Fatal(err)
		}
		link.PostedBy = user
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}

func sanitizedOrderByClause(orderBy *LinkOrderByInput) string {
	// To protect from arbitary/malicious inputs
	// TODO: DRY, but its fine since its 2 atm
	createdAtClause := ""
	switch orderBy.CreatedAt {
	case "asc", "desc":
		createdAtClause = fmt.Sprintf("CreatedAt %s", orderBy.CreatedAt)
	default:
		createdAtClause = ""
	}

	descriptionClause := ""
	switch orderBy.Description {
	case "asc", "desc":
		descriptionClause = fmt.Sprintf("Description %s", orderBy.Description)
	default:
		descriptionClause = ""
	}

	var clause string
	if createdAtClause == "" && descriptionClause == "" {
		// No sort specified, provide default sorting
		clause = "CreatedAt Desc"
	} else if createdAtClause != "" && descriptionClause != "" {
		// Both sort specified: Doesnt makes sense to do createdAt sort first then desc sorts
		clause = createdAtClause + ", " + descriptionClause
	} else if descriptionClause != "" {
		clause = descriptionClause
	} else {
		clause = createdAtClause
	}

	return clause
}
