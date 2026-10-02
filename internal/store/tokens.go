package store

import (
	"database/sql"
	"time"

	"github.com/aymen/GoProject/internal/tokens"

)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore{
	return &PostgresTokenStore{
		db: db,
	}
}


type TokenStore interface{
	CreateNewToken(userID int,ttl time.Duration,scope string) (*tokens.Token,error)
	Insert(*tokens.Token) error
	DeleteAllTokensForUser(userID int,scope string) error
}

func (p *PostgresTokenStore) CreateNewToken(userID int,ttl time.Duration,scope string) (*tokens.Token,error){
	token, err := tokens.GenerateToken(userID,ttl,scope)
	if err != nil {
		return nil,err
	}
	err = p.Insert(token)

	
    return token,err
}
func (p *PostgresTokenStore) Insert(token *tokens.Token) error {
	query:=	`
    INSERT INTO tokens (hash, user_id,expiry,scope)
	VALUES ($1,$2,$3,$4)
`
_,err:= p.db.Exec(query,token.Hash,token.UserID,token.Expiry,token.Scope)

	
	return err
}

func (p *PostgresTokenStore) DeleteAllTokensForUser(userID int,scope string) error{
	query:= `
	DELETE FROM tokens
	WHERE scope =$1 AND usera_id = $2`

	_,err:= p.db.Exec(query,scope,userID)

	return err
}