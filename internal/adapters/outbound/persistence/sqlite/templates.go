package sqlite

import (
	"database/sql"
	"errors"
	"tarmo/internal/core/templates"
	"tarmo/internal/core/templates/domain"
)

type templateRepository struct {
	db *sql.DB
}

func NewTemplateRepository(db *SQLiteDB) *templateRepository {
	return &templateRepository{db: db.db}
}

func (r *templateRepository) FindAll() ([]*domain.Template, error) {
	// Get all templates
	rows, err := r.db.Query(`
        SELECT id, name, description, quantity, unit, difficulty, created_at, updated_at
        FROM templates
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect template data first, then close rows immediately
	type templateData struct {
		id           int
		name         string
		description  string
		quantity     float64
		unit         string
		difficulty   int
		createdAtStr string
		updatedAtStr string
	}

	var templatesData []templateData
	for rows.Next() {
		var td templateData
		if err := rows.Scan(&td.id, &td.name, &td.description, &td.quantity, &td.unit, &td.difficulty, &td.createdAtStr, &td.updatedAtStr); err != nil {
			return nil, err
		}
		templatesData = append(templatesData, td)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Now build full templates by loading steps and resources for each
	var result []*domain.Template
	for _, td := range templatesData {
		// Get steps for this template
		stepRows, err := r.db.Query(`
            SELECT step_order, name, instructions
            FROM steps WHERE template_id = ?
            ORDER BY step_order
        `, td.id)
		if err != nil {
			return nil, err
		}

		steps := []domain.Step{}
		for stepRows.Next() {
			var (
				order        int
				stepName     string
				instructions string
			)
			if err := stepRows.Scan(&order, &stepName, &instructions); err != nil {
				stepRows.Close()
				return nil, err
			}

			step, err := domain.NewStep(stepName, instructions, order)
			if err != nil {
				stepRows.Close()
				return nil, err
			}
			steps = append(steps, step)
		}
		stepRows.Close()

		if err := stepRows.Err(); err != nil {
			return nil, err
		}

		// Get resources for this template
		resourceRefs := []domain.ResourceRef{}
		resourceRows, err := r.db.Query(`
            SELECT resource_id, quantity, unit
            FROM template_resources WHERE template_id = ?
        `, td.id)
		if err != nil {
			return nil, err
		}

		for resourceRows.Next() {
			var (
				resourceID int
				quantity   float64
				unit       string
			)
			if err := resourceRows.Scan(&resourceID, &quantity, &unit); err != nil {
				resourceRows.Close()
				return nil, err
			}

			resourceRef, err := domain.NewResourceRef(resourceID, quantity, unit)
			if err != nil {
				resourceRows.Close()
				return nil, err
			}
			resourceRefs = append(resourceRefs, resourceRef)
		}
		resourceRows.Close()

		if err := resourceRows.Err(); err != nil {
			return nil, err
		}

		createdAt, err := parseTime(td.createdAtStr)
		if err != nil {
			return nil, err
		}
		updatedAt, err := parseTime(td.updatedAtStr)
		if err != nil {
			return nil, err
		}

		// Reconstruct template with steps
		tmpl, err := domain.ReconstructTemplate(
			td.id,
			td.name,
			td.quantity,
			td.unit,
			td.difficulty,
			steps,
			td.description,
			resourceRefs,
			createdAt,
			updatedAt,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, tmpl)
	}

	return result, nil
}

func (r *templateRepository) FindByID(id int) (*domain.Template, error) {
	// Search template
	row := r.db.QueryRow(`
            SELECT id, name, description, quantity, unit, difficulty, created_at, updated_at
            FROM templates WHERE id=?
    `, id)

	// Temp variables
	var (
		templateID   int
		name         string
		description  string
		quantity     float64
		unit         string
		difficulty   int
		createdAtStr string
		updatedAtStr string
	)

	err := row.Scan(&templateID, &name, &description, &quantity, &unit, &difficulty, &createdAtStr, &updatedAtStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, templates.ErrTemplateNotFound // Not found
		}
		return nil, err
	}

	stepRows, err := r.db.Query(`
        SELECT step_order, name, instructions
        FROM steps WHERE template_id = ?
        ORDER BY step_order
    `, templateID)
	if err != nil {
		return nil, err
	}

	steps := []domain.Step{}
	for stepRows.Next() {
		var (
			order        int
			stepName     string
			instructions string
		)
		if err := stepRows.Scan(&order, &stepName, &instructions); err != nil {
			stepRows.Close()
			return nil, err
		}

		step, err := domain.NewStep(stepName, instructions, order)
		if err != nil {
			stepRows.Close()
			return nil, err // Invalid step
		}
		steps = append(steps, step)
	}
	stepRows.Close()

	if err := stepRows.Err(); err != nil {
		return nil, err
	}

	resourceRefs := []domain.ResourceRef{}
	resourceRows, err := r.db.Query(`
        SELECT resource_id, quantity, unit
        FROM template_resources WHERE template_id = ?
    `, templateID)
	if err != nil {
		return nil, err
	}

	for resourceRows.Next() {
		var (
			resourceID int
			quantity   float64
			unit       string
		)
		if err := resourceRows.Scan(&resourceID, &quantity, &unit); err != nil {
			resourceRows.Close()
			return nil, err
		}

		resourceRef, err := domain.NewResourceRef(resourceID, quantity, unit)
		if err != nil {
			resourceRows.Close()
			return nil, err // Invalid resource ref
		}
		resourceRefs = append(resourceRefs, resourceRef)
	}
	resourceRows.Close()

	if err := resourceRows.Err(); err != nil {
		return nil, err
	}

	createdAt, err := parseTime(createdAtStr)
	if err != nil {
		return nil, err
	}
	updatedAt, err := parseTime(updatedAtStr)
	if err != nil {
		return nil, err
	}

	// Reconstruct
	return domain.ReconstructTemplate(
		templateID,
		name,
		quantity,
		unit,
		difficulty,
		steps,
		description,
		resourceRefs,
		createdAt,
		updatedAt,
	)
}

func (r *templateRepository) Save(tmpl *domain.Template) (id int, err error) {
	now := nowUTC()

	tx, err := r.db.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			id = 0
			return
		}
		if err = tx.Commit(); err != nil {
			id = 0
		}
	}()

	query := `
        INSERT INTO templates (name, description, quantity, unit, difficulty, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	result, err := tx.Exec(query,
		tmpl.Name(),
		tmpl.Description(),
		tmpl.QuantityValue(),
		tmpl.QuantityUnitName(),
		tmpl.Difficulty(),
		now,
		now,
	)
	if err != nil {
		return
	}

	templateID, err := result.LastInsertId()
	if err != nil {
		return
	}

	stepQuery := `INSERT INTO steps (template_id, step_order, name, instructions) VALUES (?, ?, ?, ?)`
	for _, step := range tmpl.Steps() {
		_, err = tx.Exec(stepQuery, templateID, step.Order(), step.Name(), step.Instructions())
		if err != nil {
			return
		}
	}

	resourceQuery := `INSERT INTO template_resources (template_id, resource_id, quantity, unit) VALUES (?, ?, ?, ?)`
	for _, resource := range tmpl.Resources() {
		_, err = tx.Exec(resourceQuery, templateID, resource.ResourceID(), resource.QuantityValue(), resource.QuantityUnitName())
		if err != nil {
			return
		}
	}

	id = int(templateID)
	return
}

func (r *templateRepository) Update(tmpl *domain.Template) (err error) {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
        UPDATE templates 
        SET name = ?, description = ?, quantity = ?, unit = ?, difficulty = ?, updated_at = ?
        WHERE id = ?
    `
	result, err := tx.Exec(query,
		tmpl.Name(),
		tmpl.Description(),
		tmpl.QuantityValue(),
		tmpl.QuantityUnitName(),
		tmpl.Difficulty(),
		nowUTC(),
		tmpl.ID(),
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		err = templates.ErrTemplateNotFound
		return err
	}

	_, err = tx.Exec("DELETE FROM steps WHERE template_id = ?", tmpl.ID())
	if err != nil {
		return err
	}

	stepQuery := `INSERT INTO steps (template_id, step_order, name, instructions) VALUES (?, ?, ?, ?)`
	for _, step := range tmpl.Steps() {
		_, err = tx.Exec(stepQuery, tmpl.ID(), step.Order(), step.Name(), step.Instructions())
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("DELETE FROM template_resources WHERE template_id = ?", tmpl.ID())
	if err != nil {
		return err
	}

	resourceQuery := `INSERT INTO template_resources (template_id, resource_id, quantity, unit) VALUES (?, ?, ?, ?)`
	for _, resource := range tmpl.Resources() {
		_, err = tx.Exec(resourceQuery, tmpl.ID(), resource.ResourceID(), resource.QuantityValue(), resource.QuantityUnitName())
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *templateRepository) Remove(id int) error {
	result, err := r.db.Exec(`
        DELETE FROM templates WHERE id = ?
    `, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return templates.ErrTemplateNotFound
	}

	return nil
}
