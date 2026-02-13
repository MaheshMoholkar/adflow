package repository

import (
	"context"
	"errors"

	"callflow/internal/domain/template"
	db "callflow/internal/sql/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TemplateRepository implements template.Repository
type TemplateRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewTemplateRepository creates a new template repository
func NewTemplateRepository(pool *pgxpool.Pool) *TemplateRepository {
	return &TemplateRepository{
		pool:    pool,
		queries: db.New(pool),
	}
}

func (r *TemplateRepository) GetByUserID(ctx context.Context, userID int64) ([]*template.Template, error) {
	rows, err := r.queries.GetTemplateByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, template.ErrTemplateNotFound
	}
	templates := make([]*template.Template, len(rows))
	for i, row := range rows {
		templates[i] = dbTemplateToModel(row)
	}
	return templates, nil
}

func (r *TemplateRepository) GetByID(ctx context.Context, id int64, userID int64) (*template.Template, error) {
	row, err := r.queries.GetTemplateByID(ctx, db.GetTemplateByIDParams{
		ID:     id,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, template.ErrTemplateNotFound
		}
		return nil, err
	}
	return dbTemplateToModel(row), nil
}

func (r *TemplateRepository) Create(ctx context.Context, userID int64, data template.TemplateCreate) (*template.Template, error) {
	lang := data.Language
	if lang == "" {
		lang = "en"
	}

	row, err := r.queries.CreateTemplate(ctx, db.CreateTemplateParams{
		UserID:    userID,
		Name:      data.Name,
		Body:      data.Body,
		Type:      data.Type,
		Channel:   data.Channel,
		Language:  pgtype.Text{String: lang, Valid: true},
		IsDefault: data.IsDefault,
	})
	if err != nil {
		return nil, err
	}
	return dbTemplateToModel(row), nil
}

func (r *TemplateRepository) Update(ctx context.Context, id int64, userID int64, data template.TemplateUpdate) (*template.Template, error) {
	lang := data.Language
	if lang == "" {
		lang = "en"
	}

	row, err := r.queries.UpdateTemplate(ctx, db.UpdateTemplateParams{
		ID:        id,
		UserID:    userID,
		Name:      data.Name,
		Body:      data.Body,
		Type:      data.Type,
		Channel:   data.Channel,
		Language:  pgtype.Text{String: lang, Valid: true},
		IsDefault: data.IsDefault,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, template.ErrTemplateNotFound
		}
		return nil, err
	}
	return dbTemplateToModel(row), nil
}

func (r *TemplateRepository) Delete(ctx context.Context, id int64, userID int64) error {
	return r.queries.DeleteTemplate(ctx, db.DeleteTemplateParams{
		ID:     id,
		UserID: userID,
	})
}

func dbTemplateToModel(row db.Template) *template.Template {
	var lang string
	if row.Language.Valid {
		lang = row.Language.String
	}
	return &template.Template{
		ID:        row.ID,
		UserID:    row.UserID,
		Name:      row.Name,
		Body:      row.Body,
		Type:      row.Type,
		Channel:   row.Channel,
		Language:  lang,
		IsDefault: row.IsDefault,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
