package persist

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib" //Import for side-effects, need pgx driver
	"gotraining/exercise8/model"
)

var (
	ddl = `DROP TABLE IF EXISTS phone_numbers; 
CREATE TABLE IF NOT EXISTS phone_numbers (id SERIAL, "number" VARCHAR(50));`
	insert = `INSERT INTO phone_numbers (number) VALUES
                                           ('1234567890'),('123 456 7891'),('(123) 456 7892'),
                                           ('(123) 456-7893'),('123-456-7894'),('123-456-7890'),
                                           ('1234567892'),('(123)456-7892)');`
)

type postgres struct {
	cxn *sql.DB
}

func New(driver string, dsn string) (model.PhoneRepo, error) {
	cxn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connection: %w", err)
	}
	if err := cxn.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}
	return &postgres{cxn: cxn}, nil
}

func (p *postgres) Setup() error {
	_, err := p.cxn.Exec(ddl)
	if err != nil {
		return fmt.Errorf("unable to perform define table: %w", err)
	}
	_, err = p.cxn.Exec(insert)
	if err != nil {
		return fmt.Errorf("unable to populate table data: %w", err)
	}
	return nil
}

func (p *postgres) ListAll() ([]model.PhoneNumber, error) {
	rows, err := p.cxn.Query("SELECT * FROM phone_numbers")
	if err != nil {
		return nil, fmt.Errorf("ListAll(): %w", err)
	}
	var ret []model.PhoneNumber
	for rows.Next() {
		var p model.PhoneNumber
		if err := rows.Scan(&p.ID, &p.Number); err != nil {
			return nil, fmt.Errorf("ListAll(): %w", err)
		}
		ret = append(ret, p)
	}
	return ret, rows.Close()
}

func (p *postgres) Find(number string) (*model.PhoneNumber, error) {
	var ret model.PhoneNumber
	err := p.cxn.QueryRow("SELECT * FROM phone_numbers WHERE number=$1", number).Scan(&ret.ID, &ret.Number)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNoPhoneNumber
		}
		return nil, err
	}
	return &ret, nil
}

func (p *postgres) Delete(phone model.PhoneNumber) error {
	_, err := p.cxn.Exec(`DELETE FROM phone_numbers WHERE id=$1`, phone.ID)
	if err != nil {
		return fmt.Errorf("Delete(): %w", err)
	}
	return nil
}

func (p *postgres) Update(phone model.PhoneNumber) error {
	_, err := p.cxn.Exec(`UPDATE phone_numbers SET number=$1 WHERE id=$2`, phone.Number, phone.ID)
	if err != nil {
		return fmt.Errorf("Update(): %w", err)
	}
	return nil
}
