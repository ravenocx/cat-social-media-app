package repositories

import (
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/ravenocx/cat-socialx/internal/models"
)

type CatMatchQueries struct {
	*sqlx.DB
}

func (q *CatMatchQueries) CreateCatMatch(cm *models.CatMatch) error {
	query := `INSERT INTO cat_matches VALUES($1, $2, $3, $4, $5)`

	_, err := q.Exec(query, cm.ID, cm.CatIssuerID, cm.CatMatchID, cm.Message, cm.Status)

	if err != nil {
		return err
	}

	return nil
}

func (q *CatMatchQueries) GetCatMatchByCatIds(match_catId uuid.UUID, issuer_catId uuid.UUID) ([]models.CatMatch, error) {
	cat_matches := []models.CatMatch{}

	query := `SELECT * FROM cat_matches WHERE cat_issuer_id = $1 AND cat_match_id = $2`

	if err := q.Select(&cat_matches, query, match_catId, issuer_catId); err != nil {
		return nil, err
	}

	return cat_matches, nil
}

func (q *CatMatchQueries) GetCatMatchRequests(cat_id uuid.UUID) ([]models.CatMatchDetail, error) {
	cat_matches := []models.CatMatchDetail{}

	query := `SELECT 
	cat_matches.id,
    cat_matches.message,
	cat_matches.created_at,
	u.name AS "issueruser.name",
    u.email AS "issueruser.email",
	u.created_at AS "issueruser.created_at",
	ci.id AS "issuercat.id",
    ci.name AS "issuercat.name",
    ci.race AS "issuercat.race",
    ci.sex AS "issuercat.sex",
    ci.description AS "issuercat.description",
    ci.ageinmonth AS "issuercat.ageinmonth",
    ci.imageurls AS "issuercat.imageurls",
    ci.hasmatched AS "issuercat.hasmatched",
    ci.created_at AS "issuercat.created_at",
    cm.id AS "matchcat.id",
    cm.name AS "matchcat.name",
    cm.race AS "matchcat.race",
    cm.sex AS "matchcat.sex",
    cm.description AS "matchcat.description",
    cm.ageinmonth AS "matchcat.ageinmonth",
    cm.imageurls AS "matchcat.imageurls",
    cm.hasmatched AS "matchcat.hasmatched",
    cm.created_at AS "matchcat.created_at"
	FROM cat_matches
	LEFT JOIN cats ci ON cat_matches.cat_issuer_id = ci.id
	LEFT JOIN cats cm ON cat_matches.cat_match_id = cm.id
	LEFT JOIN users u ON ci.user_id = u.id
	WHERE cat_issuer_id = $1 OR cat_match_id = $1
	ORDER BY created_at DESC
	`
	rows, err := q.Query(query, cat_id)
	// log.Printf("Query : %+v", query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		cm := models.CatMatchDetail{}
		var imgUrlMatchCat string
		var imgUrlUserCat string
		err := rows.Scan(
			&cm.ID,
			&cm.Message,
			&cm.CreatedAt,
			&cm.IssuedBy.Name,
			&cm.IssuedBy.Email,
			&cm.IssuedBy.CreatedAt,
			&cm.MatchCatDetail.ID,
			&cm.MatchCatDetail.Name,
			&cm.MatchCatDetail.Race,
			&cm.MatchCatDetail.Sex,
			&cm.MatchCatDetail.Description,
			&cm.MatchCatDetail.AgeInMonth,
			&imgUrlMatchCat,
			&cm.MatchCatDetail.HasMatched,
			&cm.MatchCatDetail.CreatedAt,
			&cm.UserCatDetail.ID,
			&cm.UserCatDetail.Name,
			&cm.UserCatDetail.Race,
			&cm.UserCatDetail.Sex,
			&cm.UserCatDetail.Description,
			&cm.UserCatDetail.AgeInMonth,
			&imgUrlUserCat,
			&cm.UserCatDetail.HasMatched,
			&cm.UserCatDetail.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		cm.MatchCatDetail.ImageUrls = strings.Split(imgUrlMatchCat , ",")
		cm.UserCatDetail.ImageUrls = strings.Split(imgUrlUserCat , ",")
		cat_matches = append(cat_matches, cm)
	}

	return cat_matches, nil
}

func (q *CatMatchQueries) UpdateCatMatch(id uuid.UUID, status string) error {
	query := `UPDATE cat_matches SET status = $2 WHERE id = $1`

	_, err := q.Exec(query, id, status)

	if err != nil {
		return err
	}

	return nil
}

func (q *CatMatchQueries) GetCatMatchById(id uuid.UUID) ([]models.CatMatch, error) {
	catmatch := []models.CatMatch{}

	query := `SELECT * FROM cat_matches WHERE id = $1`

	err := q.Select(&catmatch, query, id)
	if err != nil {
		return nil, err
	}

	return catmatch, nil
}

func (q *CatMatchQueries) DeleteCatMatchExceptNotPending(id uuid.UUID) error {
	query := `DELETE FROM cat_matches WHERE id = $1 AND status = 'pending'`

	_, err := q.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (q *CatMatchQueries) DeleteCatMatchById(id uuid.UUID) error {
	query := `DELETE FROM cat_matches WHERE id = $1`

	_, err := q.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}