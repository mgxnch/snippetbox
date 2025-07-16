package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet holds the data from the snippets table.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel interacts with the database.
type SnippetModel struct {
	DB *sql.DB
}

// Custom errors
var ErrNoRecord = errors.New("models: no matching record found")

// Insert inserts the snippet into the database.
func (m *SnippetModel) Insert(title, content string, expires int) (int, error) {
	// Prepared statement
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Get fetches the snippet with the specified id.
func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires from snippets where 
	expires > UTC_TIMESTAMP() and id = ?`

	row := m.DB.QueryRow(stmt, id)

	var snippet Snippet
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return &snippet, nil
}

// Latest returns the 10 most recently snippets.
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires from snippets 
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Only call this after you are sure that sql.DB.Query didnt fail,
	// or else you will get a panic trying to close a nil resultset
	defer rows.Close()

	var snippets []*Snippet
	for rows.Next() {
		var snippet Snippet
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, &snippet)
	}

	// There could still be an error after iterating through the entire resultset,
	// and we must handle them like this
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
