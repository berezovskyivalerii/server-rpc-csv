package repository

import (
	"berezovskyivalerii/server-rpc-csv/pkg/domain"
	"context"
	"database/sql"
	"time"
)

type Product struct {
	db *sql.DB
}

func NewProduct(db *sql.DB) *Product {
	return &Product{
		db: db,
	}
}

func (r *Product) Fetch(ctx context.Context, req []domain.Product) error {
	// Начало транзакции
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// При ошибки откатываем транзакции обратно
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	//Перебираем все продукты в списке
	for _, p := range req {
		var currentPrice float64
		var priceChangeCount int32

		err = tx.QueryRowContext(ctx,
			"SELECT price, price_change_count FROM products WHERE id = ?", p.ID).Scan(&currentPrice, &priceChangeCount)

		// Если продукта нет, то создаем новый
		if err == sql.ErrNoRows {
			_, err = tx.ExecContext(ctx, `INSERT INTO products (id, price, last_request, last_updated, price_change_count)
				 VALUES (?, ?, ?, ?, ?)`, p.ID, p.Name, p.Price, time.Now(), time.Now(), 0)
			if err != nil {
				return err
			}
		} else if err != nil { // В случае если происходит другая ошибка (не обрабатываем)
			return nil
		} else {
			// Если продукт есть, то проверяем изменилась ли цена
			if p.Price != currentPrice { // Если цена изменилась, то обновляем цену, время последнего запроса и количество изменений цены
				_, err = tx.ExecContext(ctx, `UPDATE products
					 SET price = ?, last_updated = ?, price_change_count = ?
					 WHERE id = ?`, p.Price, time.Now(), priceChangeCount+1, p.ID)
				if err != nil {
					return err
				}
			} else { // Если цена не изменилась, то обновляем время последнего запроса
				_, err = tx.ExecContext(ctx,
					`UPDATE products
					 SET last_request = ?
					 WHERE id = ?`, time.Now(), p.ID,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	// Фиксируем транзакцию
	return tx.Commit()
}
