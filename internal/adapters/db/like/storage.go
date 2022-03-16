package like

import (
	"database/sql"
	"fmt"
	"forum/internal/domain/like"
	"forum/internal/model"
	"forum/pkg"
	"log"
	"sort"
)

type likeStorage struct {
	db *sql.DB
}

func NewLikeStorage(db *sql.DB) like.LikeStorage {
	return &likeStorage{
		db: db,
	}
}

func (l *likeStorage) GetAllPostLike(likeAndDislike *model.PostLikeAndDislike) (*model.PostLikeAndDislike, error) {
	log.Println("GetAllPostLike POST-ID---------->", likeAndDislike.PostId)
	fmt.Println("Start GetAll Post Like storage  method")
	rowLike := l.db.QueryRow(`SELECT COUNT(*) 
		FROM post_likes
		WHERE post_id = $1 and reaction = 'like'`, likeAndDislike.PostId)
	err := rowLike.Scan(&likeAndDislike.Like)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ERROR GetAll Post Like storage method ErrNoRows 'LIKE POST' :--> %v\n", err)
			likeAndDislike.Like = 0
		} else {
			log.Printf("ERROR GetAll Post Like storage method 'LIKE POST':--> %v\n", err)
			return likeAndDislike, err
		}
	}
	rowDislike := l.db.QueryRow(`SELECT COUNT(*) 
		FROM post_likes
		WHERE post_id = $1 and reaction = 'dislike'`, likeAndDislike.PostId)
	err = rowDislike.Scan(&likeAndDislike.DisLike)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ERROR GetAll Post Like storage method ErrNoRows 'DISLIKE POST' :--> %v\n", err)
			likeAndDislike.DisLike = 0
		} else {
			log.Printf("ERROR GetAll Post Like storage method 'DISLIKE POST' :--> %v\n", err)
			return likeAndDislike, err
		}
	}
	rowMyReaction := l.db.QueryRow(`SELECT reaction FROM post_likes
		WHERE post_id = $1 and user_uuid = $2`, likeAndDislike.PostId, likeAndDislike.AuthorUUID)
	err = rowMyReaction.Scan(&likeAndDislike.VoitStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ERROR GetAll Post Like storage method ErrNoRows 'My user left no response'  :--> %v\n", err)
			likeAndDislike.VoitStatus = ""
		} else {
			log.Printf("ERROR GetAll Post Like storage method 'rowMyReaction' :--> %v\n", err)
			return likeAndDislike, err
		}
	}
	if likeAndDislike.VoitStatus != "" {
		likeAndDislike.MyReaction = true
	}
	fmt.Println("END GetAll Post Like storage  method")
	return likeAndDislike, nil
}

func (l *likeStorage) GetAllCommentLike(likeAndDislike *model.CommentLikeAndDislike) (*model.CommentLikeAndDislike, error) {
	fmt.Println("Start GetAll Comment Like storage  method")
	rowLike := l.db.QueryRow(`SELECT COUNT(*) 
		FROM comment_likes
		WHERE comment_id = $1 and reaction = 'like'`, likeAndDislike.CommentId)
	err := rowLike.Scan(&likeAndDislike.Like)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ERROR GetAll Comment Like storage method ErrNoRows 'LIKE Comment' :--> %v\n", err)
		} else {
			log.Printf("ERROR GetAll Comment Like storage method 'LIKE Comment':--> %v\n", err)
			return likeAndDislike, err
		}
	}
	rowDislike := l.db.QueryRow(`SELECT COUNT(*) 
		FROM comment_likes
		WHERE comment_id = $1 and reaction = 'dislike'`, likeAndDislike.CommentId)
	err = rowDislike.Scan(&likeAndDislike.DisLike)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ERROR GetAll Comment Like storage method ErrNoRows 'DISLIKE Comment' :--> %v\n", err)
		} else {
			log.Printf("ERROR GetAll Comment Like storage method 'DISLIKE Comment' :--> %v\n", err)
			return likeAndDislike, err
		}
	}
	rowMyReaction := l.db.QueryRow(`SELECT reaction FROM comment_likes
		WHERE comment_id = $1 and user_uuid = $2`, likeAndDislike.CommentId, likeAndDislike.AuthorUUID)
	err = rowMyReaction.Scan(&likeAndDislike.VoitStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("ERROR GetAll Comment Like storage method ErrNoRows 'My user left no response'  :--> %v\n", err)
		} else {
			log.Printf("ERROR GetAll Comment Like storage method 'rowMyReaction' :--> %v\n", err)
			return likeAndDislike, err
		}
	}
	if likeAndDislike.VoitStatus != "" {
		likeAndDislike.MyReaction = true
	}
	fmt.Println("END GetAll Comment Like storage  method")
	return likeAndDislike, nil
}

func (l *likeStorage) GetOneCommentLike(likeAndDislike *model.CommentLikeAndDislike) ([]*model.CommentLikeAndDislike, error) {
	fmt.Println("Start GetOneCommentLike storage method")
	rows, err := l.db.Query(`SELECT * FROM comment_likes
		WHERE user_uuid=$1 and comment_id=$2`, likeAndDislike.AuthorUUID, likeAndDislike.CommentId)
	if err != nil {
		log.Println(err)
		// log.Fatal("GetOneCommentLike storage QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	commentLike := []*model.CommentLikeAndDislike{}

	for rows.Next() {
		oneCommentLike := &model.CommentLikeAndDislike{}
		err := rows.Scan(&oneCommentLike.Id, &oneCommentLike.Reaction, &oneCommentLike.AuthorUUID, &oneCommentLike.CommentId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR GetOneCommentLike storage  method ErrNoRows:--> %v\n", err)
				return commentLike, err
			}
			log.Printf("ERROR GetOneCommentLike storage  method :--> %v\n", err)
			return commentLike, err
		}
		commentLike = append(commentLike, oneCommentLike)
	}
	fmt.Println("END GetOneCommentLike storage  method")
	// if len(commentLike) != 0 {
	// 	return true, nil
	// }
	return commentLike, nil
}

func (l *likeStorage) CreateCommentLike(likeAndDislike *model.CommentLikeAndDislike) error {
	fmt.Println("Start Comment Post Like storage method")
	result, err := l.db.Exec(`INSERT INTO comment_likes
	(reaction, user_uuid, comment_id)
	VALUES ($1,$2,$3)`,
		likeAndDislike.Reaction, likeAndDislike.AuthorUUID, likeAndDislike.CommentId)
	if err != nil {
		log.Printf("ERROR Comment Post Like storage method db.EXEC func:--> %v\n", err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR Comment Post Like method LastInsert func:--> %v\n", err)
		return err
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
	fmt.Println(id)
	fmt.Println("END Comment Post Like storage method")
	return nil
}

func (l *likeStorage) DeleteCommentLike(likeAndDislike *model.CommentLikeAndDislike) error {
	_, err := l.db.Exec(`DELETE FROM comment_likes
	WHERE comment_id= $2 and user_uuid = $3`, likeAndDislike.CommentId, likeAndDislike.AuthorUUID)
	if err != nil {
		log.Printf("ERROR Delete Comment Like storage method: %v\n", err)
		return err
	}
	return nil
}

func (l *likeStorage) UpdateCommentLike(likeAndDislike *model.CommentLikeAndDislike) error {
	_, err := l.db.Exec(`UPDATE comment_likes
	SET reaction=$1
	WHERE comment_id= $2 and user_uuid = $3`, likeAndDislike.Reaction, likeAndDislike.CommentId, likeAndDislike.AuthorUUID)
	if err != nil {
		log.Printf("ERROR Update Comment Like storage func :--->%v\n", err)
		return err
	}
	return nil
}

func (l *likeStorage) CreatePostLike(likeAndDislike *model.PostLikeAndDislike) error {
	fmt.Println("Start Create Post Like storage method")
	result, err := l.db.Exec(`INSERT INTO post_likes
	(reaction, user_uuid, post_id)
	VALUES ($1,$2,$3)`,
		likeAndDislike.Reaction, likeAndDislike.AuthorUUID, likeAndDislike.PostId)
	if err != nil {
		log.Printf("ERROR Create Post Like storage method db.EXEC func:--> %v\n", err)
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("ERROR Create Post Like method LastInsert func:--> %v\n", err)
		return err
	}
	fmt.Println(result.LastInsertId())
	fmt.Println(result.RowsAffected())
	fmt.Println(id)
	fmt.Println("END Create Post Like storage method")
	return nil
}

func (l *likeStorage) DeletePostLike(likeAndDislike *model.PostLikeAndDislike) error {
	_, err := l.db.Exec(`DELETE FROM post_likes
	WHERE post_id= $2 and user_uuid = $3`, likeAndDislike.PostId, likeAndDislike.AuthorUUID)
	if err != nil {
		log.Printf("ERROR Delete Post Like storage method: %v\n", err)
		return err
	}
	return nil
}

func (l *likeStorage) UpdatePostLike(likeAndDislike *model.PostLikeAndDislike) error {
	_, err := l.db.Exec(`UPDATE post_likes
	SET reaction=$1
	WHERE post_id= $2 and user_uuid = $3`, likeAndDislike.Reaction, likeAndDislike.PostId, likeAndDislike.AuthorUUID)
	if err != nil {
		log.Printf("ERROR UpdatePostLike storage func :--->%v\n", err)
		return err
	}
	return nil
}

func (l *likeStorage) GetOnePostLike(likeAndDislike *model.PostLikeAndDislike) ([]*model.PostLikeAndDislike, error) {
	log.Println("Start GetOnePostLike storage method")
	rows, err := l.db.Query(`SELECT * FROM post_likes
		WHERE user_uuid=$1 and post_id=$2`, likeAndDislike.AuthorUUID, likeAndDislike.PostId)
	if err != nil {
		log.Println(err)
		// log.Fatal("GetOnePostLike storage QUERY ERROR ")
		return nil, err
	}
	defer rows.Close()
	postLike := []*model.PostLikeAndDislike{}

	for rows.Next() {
		onePostLike := &model.PostLikeAndDislike{}
		err := rows.Scan(&onePostLike.Id, &onePostLike.Reaction, &onePostLike.AuthorUUID, &onePostLike.PostId)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR GetOnePostLike storage  method ErrNoRows :--> %v\n", err)
				return postLike, err
			}
			log.Printf("ERROR GetOnePostLike storage  method :--> %v\n", err)
			return postLike, err
		}
		postLike = append(postLike, onePostLike)
	}
	log.Println("END GetOnePostLike storage  method")
	// if len(postLike) != 0 {
	// 	return true, nil
	// }
	// return false, nil
	return postLike, nil
}

func (l *likeStorage) GetIdPostByLike(userUUID string) ([]int, error) {
	log.Println("STARTING READ ID POSTS ----> GetIdPostByLike ")
	listId, err := l.getIdPostByPostLike(userUUID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("ERROR GET ID POSTS BY POST LIKES :--->%v\n", err)
		return nil, err
	}
	fmt.Println("slice listId------------>", listId)
	listIdComment, err := l.getIdPostByCommentLike(userUUID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("ERROR GET ID POSTS BY POST LIKES :--->%v\n", err)
		return nil, err
	}
	listId = append(listId, listIdComment...)
	fmt.Println("slice listId before sort ------------>", listId)
	sort.Ints(listId)
	fmt.Println("slice listId after sort before remove ------------>", listId)
	list := pkg.RemoveDuplicates(listId)
	fmt.Println("slice listId after sort and remove------------>", list)
	log.Println("STARTING READ ID POSTS ----> GetIdPostByLike ")
	return list, nil
}

func (l *likeStorage) getIdPostByPostLike(userUUID string) ([]int, error) {
	log.Println("Start Like storage get Id Posts method")
	rowsPostLike, err := l.db.Query(`SELECT post_id FROM post_likes
	WHERE user_uuid = $1
	ORDER BY post_id DESC`, userUUID)
	if err != nil {
		log.Println(err)
		log.Fatal("Like STORAGE get Id Posts By (POST LIKE) QUERY ERROR ")
	}
	defer rowsPostLike.Close()
	var listId []int

	for rowsPostLike.Next() {
		var id int
		err := rowsPostLike.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR Like storage get Id Posts By (POST LIKE) method :--> %v\n", err)
				return listId, err
			}
			log.Printf("ERROR Like storage get Id Posts By (POST LIKE) method:--> %v\n", err)
			return nil, err
		}
		listId = append(listId, id)
	}
	log.Println("END Like storage get Id Posts By (POST LIKE) method")
	return listId, nil
}

func (l *likeStorage) getIdPostByCommentLike(userUUID string) ([]int, error) {
	log.Println("Start Like storage get Id Posts method")
	rowsCommentLike, err := l.db.Query(`SELECT DISTINCT comments.post_id
		FROM comment_likes LEFT JOIN comments
		WHERE comment_likes.user_uuid = $1
		ORDER BY comments.post_id DESC`, userUUID)
	if err != nil {
		log.Println(err)
		// log.Fatal("Like STORAGE get Id Posts By (COMMENTS LIKE) QUERY ERROR ")
		return nil, err
	}
	defer rowsCommentLike.Close()
	var listId []int
	for rowsCommentLike.Next() {
		var id int
		err := rowsCommentLike.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("ERROR Like storage get Id Posts By (COMMENTS LIKE) method :--> %v\n", err)
				return listId, err
			}
			log.Printf("ERROR Like storage get Id Posts By (COMMENTS LIKE) method:--> %v\n", err)
			return nil, err
		}
		listId = append(listId, id)
	}
	log.Println("END Like storage get Id Posts By (COMMENTS LIKE) method")

	return listId, nil
}
