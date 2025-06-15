package repositories

import (
	"database/sql"

	models "github.com/coworker-match-api/gen/go"
)

type IUserHobbyRepo interface {
	CreateUserHobby(*models.CreateUserHobbyReq, string) (*models.CreateUserHobbyRes, error)
	GetAllUserHobby(string) (*models.GetUserHobbyRes, error)
	UpdateUserHobby(*models.UpdateUserHobbyReq, string) (*models.UpdateUserHobbyRes, error)
	GetUserHobbiesByCategory(userId string) (map[string]struct {
		Name  string
		Count int
	}, error)
	GetTotalHobbiesByCategory() (map[string]struct {
		Name  string
		Count int
	}, error)
}

type userHobbyRepo struct {
	db *sql.DB
}

func NewUserHobbyRepo(db *sql.DB) IUserHobbyRepo {
	return &userHobbyRepo{db: db}
}

func (uhr *userHobbyRepo) CreateUserHobby(req *models.CreateUserHobbyReq, userId string) (*models.CreateUserHobbyRes, error) {
	// トランザクションを開始
	tx, err := uhr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// user_hobbies テーブルから古いエントリを削除
	_, err = tx.Exec(`DELETE FROM user_hobbies WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}

	// user_hobbies テーブルに新しいエントリを挿入
	for _, hobbyID := range req.HobbyIds {
		_, err := tx.Exec(`
			INSERT INTO user_hobbies (user_id, hobby_id)
			VALUES ($1, $2)
		`, userId, hobbyID)
		if err != nil {
			return nil, err
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var result models.CreateUserHobbyRes = models.CreateUserHobbyRes(*req)
	return &result, nil
}

func (uhr *userHobbyRepo) GetAllUserHobby(userId string) (*models.GetUserHobbyRes, error) {
	query := `
		SELECT
			h.hobby_id,
			h.hobby_name
		FROM
			users u
		LEFT JOIN
			user_hobbies uh ON u.user_id = uh.user_id
		LEFT JOIN
			hobbies h ON uh.hobby_id = h.hobby_id
		WHERE
			u.user_id = $1;
	`

	rows, err := uhr.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &models.GetUserHobbyRes{Hobbies: []models.Hobby{}}

	for rows.Next() {
		var hobbyId, hobbyName *string
		if err := rows.Scan(&hobbyId, &hobbyName); err != nil {
			return nil, err
		}
		if hobbyId != nil && hobbyName != nil {
			response.Hobbies = append(response.Hobbies, models.Hobby{
				HobbyId:   *hobbyId,
				HobbyName: *hobbyName,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return response, nil
}

func (uhr *userHobbyRepo) UpdateUserHobby(req *models.UpdateUserHobbyReq, userId string) (*models.UpdateUserHobbyRes, error) {
	// トランザクションを開始
	tx, err := uhr.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// user_hobbies テーブルから古いエントリを削除
	_, err = tx.Exec(`DELETE FROM user_hobbies WHERE user_id = $1`, userId)
	if err != nil {
		return nil, err
	}

	// user_hobbies テーブルに新しいエントリを挿入
	for _, hobbyId := range req.HobbyIds {
		_, err := tx.Exec(`
			INSERT INTO user_hobbies(
				user_id,
				hobby_id
			)
			VALUES(
				$1,
				$2
			)
		`, userId, hobbyId)
		if err != nil {
			return nil, err
		}
	}

	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	var result models.UpdateUserHobbyRes = models.UpdateUserHobbyRes(*req)
	return &result, nil
}

func (uhr *userHobbyRepo) GetUserHobbiesByCategory(userId string) (map[string]struct {
	Name  string
	Count int
}, error) {
        query := `
        SELECT
            c.category_id,
            c.category_name,
            COUNT(h.category_id) as hobby_count
        FROM
            user_hobbies uh
        JOIN
            hobbies h ON uh.hobby_id = h.hobby_id
        JOIN
            categories c ON h.category_id = c.category_id
        WHERE
            uh.user_id = $1
        GROUP BY
            c.category_id;
    `

	rows, err := uhr.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userHobbies := make(map[string]struct {
		Name  string
		Count int
	})
	for rows.Next() {
		var categoryId, categoryName string
		var hobbyCount int
		if err := rows.Scan(&categoryId, &categoryName, &hobbyCount); err != nil {
			return nil, err
		}
		userHobbies[categoryId] = struct {
			Name  string
			Count int
		}{
			Name:  categoryName,
			Count: hobbyCount,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userHobbies, nil
}

func (uhr *userHobbyRepo) GetTotalHobbiesByCategory() (map[string]struct {
	Name  string
	Count int
}, error) {
	query := `
        SELECT
            c.category_id,
			c.category_name,
            COUNT(h.hobby_id) as hobby_count
        FROM
            hobbies h
		JOIN
			categories c ON h.category_id = c.category_id
        GROUP BY
            c.category_id;
    `

	rows, err := uhr.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	totalHobbies := make(map[string]struct {
		Name  string
		Count int
	})
	for rows.Next() {
		var categoryId string
		var categoryName string
		var hobbyCount int

		if err := rows.Scan(&categoryId, &categoryName, &hobbyCount); err != nil {
			return nil, err
		}
		totalHobbies[categoryId] = struct {
			Name  string
			Count int
		}{
			Name:  categoryName,
			Count: hobbyCount,
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return totalHobbies, nil
}
