package repository

import (
	"database/sql"
	"encoding/json"
	"github.com/lib/pq"
	"github.kostyadodich/demo/pkg/model"
)

type Dialer struct {
	db *sql.DB
}

func NewDialer(db *sql.DB) *Dialer {
	return &Dialer{
		db: db,
	}
}

func (d *Dialer) Create(dialer model.Dialer) error {
	tx, err := d.db.Begin()

	err = tx.QueryRow("INSERT INTO dialers (name, phone_number) VALUES ($1, $2) RETURNING id",
		dialer.Name,
		dialer.PhoneNumber,
	).Scan(&dialer.ID)
	if err != nil {
		tx.Rollback()
	}

	stm, err := tx.Prepare(pq.CopyIn("dialer_item", "dialer_id", "item_id"))
	if err != nil {
		tx.Rollback()
	}

	for _, item := range dialer.Items {
		_, err := stm.Exec(dialer.ID, item.ID)
		if err != nil {
			tx.Rollback()
		}
	}
	_, err = stm.Exec()
	if err != nil {
		tx.Rollback()
	}

	err = stm.Close()
	if err != nil {
		tx.Rollback()
	}

	return tx.Commit()
}

func (d *Dialer) DialerList() ([]model.Dialer, error) {
	rows, err := d.db.Query(
		`SELECT d.id, d.name, d.phone_number, json_agg(i)
      			FROM dialers d JOIN dialer_item di ON d.id = di.dialer_id 
      			JOIN items i ON i.id = di.item_id
      		   GROUP BY d.id, d.name, d.phone_number`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []model.Dialer
	for rows.Next() {
		dialer := model.Dialer{}
		itemsJson := ""

		err := rows.Scan(&dialer.ID, &dialer.Name, &dialer.PhoneNumber, &itemsJson)
		if err != nil {
			return nil, err
		}

		var items []model.Item
		if err := json.Unmarshal([]byte(itemsJson), &items); err != nil {
			return nil, err
		}
		dialer.Items = items

		result = append(result, dialer)
	}

	return result, nil
}

func (d *Dialer) GetDialerByID(id int64) (*model.Dialer, error) {
	row, err := d.db.Query(`
	SELECT d.id, d.name, d.phone_number, i.id, i.name, i.is_legal
	FROM dialers d
    JOIN dialer_item di ON d.id = di.dialer_id
    JOIN items i ON di.item_id = i.id
    WHERE d.id=$1`, id)

	if err != nil {
		return nil, err
	}

	dialer := model.Dialer{}
	item := model.Item{}
	for row.Next() {
		err := row.Scan(&dialer.ID, &dialer.Name, &dialer.PhoneNumber, &item.ID, &item.Name, &item.IsLegal)
		if err != nil {
			return nil, err
		}

		dialer.Items = append(dialer.Items, item)
	}

	return &dialer, nil
}

func (d *Dialer) UpdateDialer(dialer model.Dialer) error {
	_, err := d.db.Exec("UPDATE dialers SET name=$1, phone_number=$2 WHERE id=$3",
		dialer.Name,
		dialer.PhoneNumber,
		dialer.ID,
	)

	return err
}

func (d *Dialer) DeleteDialer(id int64) error {
	_, err := d.db.Exec("DELETE FROM dialer WHERE id = $1", id)

	return err
}

func (d *Dialer) GetByItem(itemId int64) ([]model.Dialer, error) {
	row, err := d.db.Query(
		`SELECT d.id, d.name, d.phone_number, json_agg(i) 
				FROM dialers d JOIN dialer_item di ON d.id = di.dialer_id
				JOIN items i ON i.id = di.item_id
				WHERE di.item_id = 1
				GROUP BY d.id, d.name, d.phone_number`, itemId)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	resutl := []model.Dialer{}
	for row.Next() {
		dialer := model.Dialer{}
		itemJson := ""

		err := row.Scan(dialer.ID, &dialer.Name, &dialer.PhoneNumber, &itemJson)
		if err != nil {
			return nil, err
		}

		item := []model.Item{}
		err = json.Unmarshal([]byte(itemJson), &item)
		if err != nil {
			return nil, err
		}
		dialer.Items = item
		resutl = append(resutl, dialer)
	}

	return resutl, nil
}
