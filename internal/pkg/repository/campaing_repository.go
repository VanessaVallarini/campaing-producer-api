package repository

import (
	"campaing-producer-service/internal/model"
	"campaing-producer-service/internal/pkg/db"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lockp111/go-easyzap"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	TABLE_NAME        = "campaing"
	allCampaingFields = `id, user_id, slug_id, merchant_id, created_at, updated_at, active,lat , long, clicks, impressions`
)

type ICampaingRepository interface {
	GetById(ctx context.Context, campaingId uuid.UUID) (model.Campaing, error)
	List(ctx context.Context, filters model.CampaingListingFilters) ([]model.Campaing, model.Paging, error)
}

type CampaingRepository struct {
	conn db.IDb
}

func NewFormRepository(conn db.IDb) *CampaingRepository {
	return &CampaingRepository{
		conn: conn,
	}
}

func (repo *CampaingRepository) GetById(ctx context.Context, campaingId uuid.UUID) (model.Campaing, error) {
	var campaing model.Campaing

	query := `SELECT ` + allCampaingFields + ` from campaing WHERE id = $1`

	rows, err := repo.conn.Query(ctx, query, campaingId)
	if err != nil {
		easyzap.Error(ctx, err, "error on get query for campaing id",
			zap.String("campaingId", fmt.Sprintf("%v", campaingId)))
		return model.Campaing{}, err
	}

	defer rows.Close()

	err = rows.Scan(
		&campaing.Id, &campaing.UserId, &campaing.SlugId, &campaing.MerchantId,
		&campaing.CreatedAt, &campaing.UpdatedAt, &campaing.Active, &campaing.Lat,
		&campaing.Long, &campaing.Clicks, &campaing.Impressions)
	if err != nil {
		easyzap.Error(ctx, err, "error on scan for campaing id",
			zap.String("campaingId", fmt.Sprintf("%v", campaingId)))
		return model.Campaing{}, err
	}

	// Verificar se houve algum erro durante a iteração
	if err := rows.Err(); err != nil {
		return model.Campaing{}, err
	}

	return model.Campaing{}, nil
}

func (repo *CampaingRepository) List(ctx context.Context, filters model.CampaingListingFilters) ([]model.Campaing, model.Paging, error) {
	dataQuery := `SELECT ` + allCampaingFields + ` from campaing WHERE 1 = 1`
	countQuery := `SELECT count(id) as "total" from campaing WHERE 1 = 1`

	paramCount := 1
	extraParams := make([]interface{}, 0, 6)

	// Add filters as needed
	if len(filters.Ids) > 0 {
		idsQuery := ""
		for _, id := range filters.Ids {
			idsQuery += fmt.Sprintf("$%d, ", paramCount)
			paramCount++
			extraParams = append(extraParams, id.String())
		}
		idsQuery = idsQuery[0 : len(idsQuery)-2]

		extra := fmt.Sprintf(" and id in (%s) ", idsQuery)
		dataQuery, countQuery = dataQuery+extra, countQuery+extra
	}

	// Add filters as needed
	if filters.Active != nil {
		extra := fmt.Sprintf(" and active = $%d ", paramCount)
		paramCount++
		dataQuery, countQuery = dataQuery+extra, countQuery+extra
		extraParams = append(extraParams, filters.Active)
	}

	// default limit size
	if filters.Size == 0 {
		filters.Size = 50
	}

	// Select total count
	type counting struct {
		Total int `db:"total"`
	}
	countingResult := counting{}
	countRow := repo.conn.QueryRow(countQuery, extraParams...)
	err := countRow.Scan(&countingResult.Total)
	if err != nil {
		return nil, model.Paging{}, errors.Wrap(err, "Failed to get total items count for campaing table")
	}

	// pagination params
	dataQuery = dataQuery + fmt.Sprintf(" limit $%d offset $%d", paramCount, paramCount+1)
	paramCount += 2
	extraParams = append(extraParams, filters.Size, filters.Size*filters.Page)

	// Select paged data
	result := []model.Campaing{}
	rows, err := repo.conn.Query(ctx, dataQuery, extraParams...)
	if err != nil {
		return nil, model.Paging{}, errors.Wrap(err, "Failed to get campaing list data")
	}
	defer rows.Close()

	for rows.Next() {
		var campaing model.Campaing
		err := rows.Scan(
			&campaing.Id,
			&campaing.UserId,
			&campaing.SlugId,
			&campaing.MerchantId,
			&campaing.CreatedAt,
			&campaing.UpdatedAt,
			&campaing.Active,
			&campaing.Lat,
			&campaing.Long,
			&campaing.Clicks,
			&campaing.Impressions,
		)
		if err != nil {
			return nil, model.Paging{}, err
		}
		result = append(result, campaing)
	}

	if err := rows.Err(); err != nil {
		return nil, model.Paging{}, errors.Wrap(err, "Failed to read row when listing campaings data")
	}

	return result, model.Paging{Page: filters.Page, Size: len(result), TotalItems: countingResult.Total}, nil
}
