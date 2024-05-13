package repositories

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ravenocx/cat-socialx/internal/models"
)

type CatQueries struct {
	*sqlx.DB
}

func (q *CatQueries) GetCats() ([]models.Cats, error) {
	cats := []models.Cats{}

	query := `SELECT * FROM cats AND deleted_at IS NULL`

 	rows, err := q.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var cat = models.Cats{}
		var imgUrlStr string
		var err = rows.Scan(
			&cat.ID,
			&cat.UserID,
			&cat.Name,
			&cat.Race,
			&cat.Sex,
			&cat.AgeInMonth,
			&cat.Description,
			&cat.HasMatched,
			&imgUrlStr,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&cat.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		cat.ImageUrls = strings.Split(imgUrlStr, ",")
		cats = append(cats, cat)
	}

	return cats, nil
}

func (q *CatQueries) GetCatsData(appendQuery string) ([]models.CatData, error) {
	query := "SELECT id,name,race,sex,ageinmonth,imageurls,description,hasmatched,created_at FROM cats WHERE deleted_at IS NULL"
	rows, err := q.Query(query + appendQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.CatData
	for rows.Next() {
		var each = models.CatData{}
		var imgUrlStr string
		var err = rows.Scan(
			&each.ID,
			&each.Name,
			&each.Race,
			&each.Sex,
			&each.AgeInMonth,
			&imgUrlStr,
			&each.Description,
			&each.HasMatched,
			&each.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		each.ImageUrls = strings.Split(imgUrlStr, ",")
		result = append(result, each)
	}

	return result, nil
}

func (q *CatQueries) GetCatById(id uuid.UUID) ([]models.Cats, error) {
	cats := []models.Cats{}

	query := `SELECT * FROM cats WHERE id = $1 AND deleted_at IS NULL`

	rows, err := q.Query(query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var cat = models.Cats{}
		var imgUrlStr string
		var err = rows.Scan(
			&cat.ID,
			&cat.UserID,
			&cat.Name,
			&cat.Race,
			&cat.Sex,
			&cat.AgeInMonth,
			&cat.Description,
			&cat.HasMatched,
			&imgUrlStr,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&cat.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		cat.ImageUrls = strings.Split(imgUrlStr, ",")
		cats = append(cats, cat)
	}

	return cats, nil
}

func (q *CatQueries) GetCatsByUserId(userId uuid.UUID) ([]models.Cats, error) {
	cats := []models.Cats{}

	query := `SELECT * FROM cats WHERE user_id = $1 AND deleted_at IS NULL`

	rows, err := q.Query(query, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var cat = models.Cats{}
		var imgUrlStr string
		var err = rows.Scan(
			&cat.ID,
			&cat.UserID,
			&cat.Name,
			&cat.Race,
			&cat.Sex,
			&cat.AgeInMonth,
			&cat.Description,
			&cat.HasMatched,
			&imgUrlStr,
			&cat.CreatedAt,
			&cat.UpdatedAt,
			&cat.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		cat.ImageUrls = strings.Split(imgUrlStr, ",")
		cats = append(cats, cat)
	}

	return cats, nil
}

func (q *CatMatchQueries) UpdateCatHasMatched(id uuid.UUID) error {
	query := `UPDATE cats SET hasmatched = true WHERE id = $1 AND deleted_at IS NULL`

	_, err := q.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (q *CatQueries) UpdateCat(id uuid.UUID, c *models.CatUpdateRequest) error {
	query := `UPDATE cats SET name = $2, race = $3, sex = $4, ageinmonth = $5, description = $6, imageurls = $7 WHERE id = $1`

	_, err := q.Exec(query, id, c.Name, c.Race, c.Sex, c.AgeInMonth, c.Description, c.ImageUrls)
	if err != nil {
		return err
	}

	return nil
}

func (q *CatQueries) CreateCat(c *models.Cat) error {
	query := `INSERT INTO cats (id, user_id, name, race, sex, ageinmonth, description, imageurls, hasmatched, created_at)
           VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := q.Exec(
		query,
		c.ID,
		c.UserID,
		c.Name,
		c.Race,
		c.Sex,
		c.AgeInMonth,
		c.Description,
		c.ImageUrls,
		c.HasMatched,
		c.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (q *CatQueries) DeleteCat(catId string, userID string) error {
	query := `UPDATE cats SET deleted_at = $1 WHERE id = $2 AND user_id = $3`

	deletedAt := time.Now()
	res, err := q.Exec(query, deletedAt, catId, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no id found")
	}

	return nil
}