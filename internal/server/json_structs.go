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

// No passwords passed here!
type JSON_User struct {
	Id         		string `json:"id"`
	Created_at 		string `json:"created_at"`
	Updated_at 		string `json:"updated_at"`
	Email      		string `json:"email"`
	Is_Chirpy_Red 	bool   `json:"is_chirpy_red"`
}

func createJSONUser(db_user database.User) JSON_User {
	return JSON_User{
		Id:         	db_user.ID.String(),
		Created_at: 	db_user.CreatedAt.String(),
		Updated_at: 	db_user.UpdatedAt.String(),
		Email:      	db_user.Email,
		Is_Chirpy_Red: 	db_user.IsChirpyRed,
	}
}

// Login Return
type JSON_LoginUser struct {
	Id         string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Email      string `json:"email"`
	Is_Chirpy_Red 	bool   `json:"is_chirpy_red"`
	Token      string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func createJSONLoginUser(db_user database.User, token string, refreshToken string) JSON_LoginUser {
	return JSON_LoginUser{
		Id:         		db_user.ID.String(),
		Created_at: 		db_user.CreatedAt.String(),
		Updated_at: 		db_user.UpdatedAt.String(),
		Email:      		db_user.Email,
		Is_Chirpy_Red: 	db_user.IsChirpyRed,
		Token:				token,
		RefreshToken:		refreshToken,
	}
}