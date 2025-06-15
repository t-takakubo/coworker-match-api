package repositories

import (
	"database/sql"

	models "github.com/coworker-match-api/gen/go"
)

type IUserRepo interface {
	CreateUser(userId string, req *models.CreateUserReq) (*models.CreateUserRes, error)
	GetUserById(userId string) (*models.GetUserRes, error)
	UpdateUser(userId string, req *models.UpdateUserReq) (*models.UpdateUserRes, error)
	IsUserExist(userId string) (bool, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) IUserRepo {
	return &userRepo{db: db}
}

func (ur *userRepo) CreateUser(userId string, req *models.CreateUserReq) (*models.CreateUserRes, error) {
	const query = `
				INSERT INTO 
					users (user_id, user_name, email, avatar_url)
				VALUES 
					($1, $2, $3, $4)
				RETURNING
					user_id, user_name, email, avatar_url
				`
	var response models.CreateUserRes
	err := ur.db.QueryRow(query, userId, req.UserName, req.Email, req.AvatarUrl).Scan(
		&response.UserId, &response.UserName, &response.Email, &response.AvatarUrl,
	)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (ur *userRepo) GetUserById(userId string) (*models.GetUserRes, error) {
	query := `
			SELECT
				user_id, user_name, email, avatar_url
			FROM 
				users
			WHERE
				user_id = $1
			`

	var response models.GetUserRes
	ur.db.QueryRow(query, userId).Scan(
		&response.UserId, &response.UserName, &response.Email, &response.AvatarUrl,
	)

	return &response, nil
}

func (r *userRepo) UpdateUser(userId string, req *models.UpdateUserReq) (*models.UpdateUserRes, error) {
	query := `
                UPDATE
                        users
                SET
                        user_name = $1,
                        email = $2,
                        avatar_url = $3,
                        updated_at = NOW()
                WHERE
                        user_id = $4
                RETURNING user_id, user_name, email, avatar_url
        `

	var res models.UpdateUserRes
	err := r.db.QueryRow(query, req.UserName, req.Email, req.AvatarUrl, userId).Scan(
		&res.UserId,
		&res.UserName,
		&res.Email,
		&res.AvatarUrl,
	)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *userRepo) IsUserExist(userId string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE user_id = $1)`
	row := r.db.QueryRow(query, userId)
	var isExist bool
	err := row.Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist, nil
}
