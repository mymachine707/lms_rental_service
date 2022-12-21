package postgres

import (
	"errors"
	"lms/protogen/rental_service"
)

/*
"id", "book_id", "book_name", "user_id", "rental_date_time" DEFAULT now(),
    "expected_return_date" TIMESTAMP,
	"return_date" TIMESTAMP,
	"book_fines" VARCHAR(255),
	"rental_status" my_type  DEFAULT 'RENTAL'
*/
// CreateRental ...
func (p Postgres) CreateRental(id string, req *rental_service.CreateRentalRequest) error {

	_, err := p.DB.Exec(`
	INSERT INTO
	 "rental_book"("id", "book_id", "book_name", "user_id", "rental_date_time", "expected_return_date", rental_status) 
	VALUES($1, $2, $3, $4, now(), $5, 'RENTAL')
	`, id, req.BookId, req.BookName, req.UserId, req.ExpectedReturnDate)
	if err != nil {
		return err
	}
	return nil
}

// GetRentalByID ...
func (p Postgres) GetRentalById(id string) (*rental_service.Rental, error) {
	res := &rental_service.Rental{}
	var returnDate *string
	var bookFines *string
	err := p.DB.QueryRow(`SELECT "id", "book_id", "book_name", "user_id", "rental_date_time", 
	"expected_return_date", "return_date", "book_fines", "rental_status"
    FROM "rental_book" WHERE "id" = $1`, id).Scan(
		&res.Id,
		&res.BookId,
		&res.BookName,
		&res.UserId,
		&res.RentalDateTime,
		&res.ExpectedReturnDate,
		&returnDate,
		&bookFines,
		&res.RentalStatus,
	)
	if err != nil {
		return res, err
	}

	if returnDate != nil {
		res.ReturnDate = *returnDate
	}

	if bookFines != nil {
		res.BookFines = *bookFines
	}

	return res, err
}

// GetRentalList ...
func (p Postgres) GetRentalList(offset, limit int, search string) (*rental_service.GetRentalListResponse, error) {
	resp := &rental_service.GetRentalListResponse{
		Rentals: make([]*rental_service.Rental, 0),
	}

	rows, err := p.DB.Queryx(`
	SELECT
	"id", "book_id", "book_name", "user_id", "rental_date_time", 
	"expected_return_date", "return_date", "book_fines", rental_status
	FROM 
		"rental_book" WHERE deleted_at IS NULL AND (book_name ILIKE '%' || $1 || '%')
	LIMIT $2
	OFFSET $3
	`, search, limit, offset)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		a := &rental_service.Rental{}
		var returnDate *string
		var bookFines *string
		err := rows.Scan(
			&a.Id, &a.BookId, &a.BookName, &a.UserId, &a.RentalDateTime, &a.ExpectedReturnDate,
			&returnDate, &bookFines, &a.RentalStatus,
		)
		if err != nil {
			return nil, err
		}

		if returnDate != nil {
			a.ReturnDate = *returnDate
		}

		if bookFines != nil {
			a.BookFines = *bookFines
		}

		resp.Rentals = append(resp.Rentals, a)
	}

	return resp, err
}

// UpdateRental ...
func (p Postgres) UpdateRental(entity *rental_service.UpdateRentalRequest) error {

	res, err := p.DB.NamedExec(`
	UPDATE 
		"rental_book"  
	SET 
		"id"=:id, "book_id"=:bi, "book_name"=:bn, "user_id"=:ui, "return_date"=now(), 
		"book_fines"=:bf, "rental_status"=:rs WHERE deleted_at IS NULL AND id=:id`, map[string]interface{}{
		"id": entity.Id,
		"bi": entity.BookId,
		"bn": entity.BookName,
		"ui": entity.UserId,
		"bf": entity.BookFines,
		"rs": entity.RentalStatus,
	})
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("rental_book not found")
}

// DeleteRental ...
func (p Postgres) DeleteRental(id string) error {
	res, err := p.DB.Exec(`UPDATE "rental_book" SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n > 0 {
		return nil
	}

	return errors.New("rental not found")
}
