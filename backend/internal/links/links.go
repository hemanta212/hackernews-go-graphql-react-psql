package links

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	database "github.com/hemanta212/hackernews-go-graphql/internal/pkg/db/mysql"
	"github.com/hemanta212/hackernews-go-graphql/internal/users"
)

type Link struct {
	ID          string
	Description string
	Url         string
	PostedBy    *users.User
	CreatedAt   string
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

func (link *Link) Save() int64 {
	stmt, err := database.Db.Prepare("INSERT INTO Links(Description, Url, UserID, CreatedAt) VALUES(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(link.Description, link.Url, link.PostedBy.ID, link.CreatedAt)
	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error: ", err.Error())
	}
	return id
}

func GetLinkByID(id int) (*Link, error) {
	stmt, err := database.Db.Prepare(
		"SELECT L.Description, L.Url, L.CreatedAt, L.UserID, U.Username, U.Email from Links L inner join Users U on L.UserID = U.ID WHERE L.ID=?")
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
                       WHERE L.Description LIKE ?
                       OR L.Url LIKE ?
                       OR L.CreatedAt LIKE ?
                       ORDER BY %s
                       LIMIT ?
                       OFFSET ?`, sanitizedOrderByClause(mods.OrderBy))

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(mods.Filter, mods.Filter, mods.Filter, mods.Limit, mods.Offset)
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
	case "asc":
	case "desc":
		createdAtClause = fmt.Sprintf("CreatedAt %s", orderBy.CreatedAt)
	default:
		createdAtClause = ""
	}

	descriptionClause := ""
	switch orderBy.Description {
	case "asc":
	case "desc":
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
