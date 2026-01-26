package server

import "github/SirVoly/chirpy/internal/database"

type JSON_Chirp struct {
	Id         string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Body       string `json:"body"`
	User_ID    string `json:"user_id"`
}

func createJSONChirp(db_chirp database.Chirp) JSON_Chirp {
	return JSON_Chirp{
		Id:         db_chirp.ID.String(),
		Created_at: db_chirp.CreatedAt.String(),
		Updated_at: db_chirp.UpdatedAt.String(),
		Body:       db_chirp.Body,
		User_ID:    db_chirp.UserID.String(),
	}
}

type JSON_User struct {
	Id         string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Email      string `json:"email"`
}

func createJSONUser(db_user database.User) JSON_User {
	return JSON_User{
		Id:         db_user.ID.String(),
		Created_at: db_user.CreatedAt.String(),
		Updated_at: db_user.UpdatedAt.String(),
		Email:      db_user.Email,
	}
}