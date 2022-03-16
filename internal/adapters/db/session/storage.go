package session

import (
	"context"
	"database/sql"
	"errors"
	"forum/internal/domain/session"
	"forum/internal/model"
	"log"
)

type sessionStorage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) session.SessionStorage {
	return &sessionStorage{
		db: db,
	}
}

func (ss *sessionStorage) Create(ctx context.Context, session *model.Session) error {
	_, err := ss.db.Exec(`INSERT INTO sessions 
		(
			user_uuid, key, date
		) 
		VALUES (
			$1, $2, $3
		)`, session.UserUUID, session.SessionUUID, session.Date)
	if err != nil {
		log.Println("Error in storage Create session for user")
		return err
	}
	return nil
}

func (ss *sessionStorage) Check(ctx context.Context, session *model.SessionDto) error {
	log.Println("STARTING CHECK SESSION STORAGE METHOD ------->")
	rows, err := ss.db.Query(`SELECT * FROM sessions 
		WHERE user_uuid = $1 and key = $2
		ORDER BY id DESC`, session.MyUUID, session.Value)
	if err != nil {
		log.Printf("ERROR Session STORAGE QUERY ERROR :--->%v\n", err)
		return err
	}
	defer rows.Close()
	sessions := []*model.Session{}

	for rows.Next() {
		oneSession := &model.Session{}
		err := rows.Scan(&oneSession.Id, &oneSession.UserUUID, &oneSession.SessionUUID, &oneSession.Date)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR Session STORAGE QUERY NO ROWS SCAN method :--->%v\n", err)
				// break
				return err
			}
			log.Printf("ERROR Session STORAGE QUERY ROWS SCAN method :--->%v\n", err)
			// continue
			return err
		}
		log.Println("CHECK SESSION STORAGE METHOD ------->", session.Value)
		log.Println("CHECK SESSION STORAGE METHOD ------->", oneSession.SessionUUID)
		sessions = append(sessions, oneSession)
	}
	if len(sessions) == 0 {
		err := errors.New("there is no such entry")
		return err
	}
	log.Println("END CHECK SESSION STORAGE METHOD ------->")
	return nil
}

func (ss *sessionStorage) Delete(ctx context.Context, session *model.SessionDto) error {
	_, err := ss.db.Exec(`DELETE FROM sessions
	WHERE user_uuid= $1`, session.MyUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ERROR i didn't find the entry in the database:--> %v\n", err)
			return nil
		}
		log.Printf("Error in storage Delete method session in DB: --> %v\n", err)
		return err
	}
	return nil
}
