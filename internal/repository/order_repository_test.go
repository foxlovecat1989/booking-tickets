package repository

import (
	"testing"

	models "tickets/internal/models/domain"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrderRepository(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)
	assert.NotNil(t, repo)
	assert.Equal(t, baseRepo, repo.BaseRepository)
}

func TestOrderRepository_CreateOrder(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)

	// Test creating a valid order
	order := &models.Order{
		Status:     "pending",
		TotalPrice: decimal.NewFromFloat(99.99),
	}

	err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
		return repo.CreateOrder(tx, order)
	})

	require.NoError(t, err)
	assert.NotZero(t, order.ID)
	assert.NotZero(t, order.CreatedAt)
	assert.Equal(t, "pending", order.Status)
	assert.True(t, order.TotalPrice.Equal(decimal.NewFromFloat(99.99)))

	// Verify the order was actually created in the database
	var dbOrder struct {
		ID         int             `db:"id"`
		CreatedAt  int64           `db:"created_at"`
		Status     string          `db:"status"`
		TotalPrice decimal.Decimal `db:"total_price"`
	}
	err = baseRepo.db.Get(&dbOrder, "SELECT id, created_at, status, total_price FROM orders WHERE id = $1", order.ID)
	require.NoError(t, err)
	assert.Equal(t, order.ID, dbOrder.ID)
	assert.Equal(t, order.CreatedAt, dbOrder.CreatedAt)
	assert.Equal(t, order.Status, dbOrder.Status)
	assert.True(t, order.TotalPrice.Equal(dbOrder.TotalPrice))
}

func TestOrderRepository_CreateOrder_ZeroPrice(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)

	order := &models.Order{
		Status:     "pending",
		TotalPrice: decimal.Zero,
	}

	err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
		return repo.CreateOrder(tx, order)
	})

	require.NoError(t, err)
	assert.NotZero(t, order.ID)
	assert.NotZero(t, order.CreatedAt)
	assert.Equal(t, "pending", order.Status)
	assert.True(t, order.TotalPrice.Equal(decimal.Zero))
}

func TestOrderRepository_CreateOrder_LargePrice(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)

	order := &models.Order{
		Status:     "pending",
		TotalPrice: decimal.NewFromFloat(999999.99),
	}

	err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
		return repo.CreateOrder(tx, order)
	})

	require.NoError(t, err)
	assert.NotZero(t, order.ID)
	assert.NotZero(t, order.CreatedAt)
	assert.Equal(t, "pending", order.Status)
	assert.True(t, order.TotalPrice.Equal(decimal.NewFromFloat(999999.99)))
}

func TestOrderRepository_CreateOrder_DifferentStatuses(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)

	testCases := []string{"pending", "confirmed", "cancelled", "completed"}

	for _, status := range testCases {
		t.Run("status_"+status, func(t *testing.T) {
			order := &models.Order{
				Status:     status,
				TotalPrice: decimal.NewFromFloat(50.00),
			}

			err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
				return repo.CreateOrder(tx, order)
			})

			require.NoError(t, err)
			assert.NotZero(t, order.ID)
			assert.NotZero(t, order.CreatedAt)
			assert.Equal(t, status, order.Status)
		})
	}
}

func TestOrderRepository_CreateOrder_TransactionRollback(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)

	// Create an order that should be rolled back
	order := &models.Order{
		Status:     "pending",
		TotalPrice: decimal.NewFromFloat(100.00),
	}

	// Simulate a transaction that will be rolled back
	err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
		// Create the order
		err := repo.CreateOrder(tx, order)
		if err != nil {
			return err
		}

		// Simulate an error that causes rollback
		return assert.AnError
	})

	require.Error(t, err)

	// Verify the order was not persisted
	var count int
	err = baseRepo.db.Get(&count, "SELECT COUNT(*) FROM orders WHERE id = $1", order.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}

func TestOrderRepository_CreateOrder_ConcurrentAccess(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)

	// Test concurrent order creation
	const numGoroutines = 10
	errors := make(chan error, numGoroutines)
	orders := make(chan *models.Order, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			order := &models.Order{
				Status:     "pending",
				TotalPrice: decimal.NewFromFloat(float64(index) + 1.00),
			}

			err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
				return repo.CreateOrder(tx, order)
			})

			if err != nil {
				errors <- err
			} else {
				orders <- order
			}
		}(i)
	}

	// Collect results
	createdOrders := make([]*models.Order, 0, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errors:
			require.NoError(t, err)
		case order := <-orders:
			createdOrders = append(createdOrders, order)
		}
	}

	// Verify all orders were created successfully
	assert.Len(t, createdOrders, numGoroutines)

	// Verify all orders have unique IDs
	orderIDs := make(map[int]bool)
	for _, order := range createdOrders {
		assert.NotZero(t, order.ID)
		assert.False(t, orderIDs[order.ID], "Duplicate order ID found: %d", order.ID)
		orderIDs[order.ID] = true
	}
}

func TestOrderRepository_CreateOrder_Integration(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)

	// Test that we can connect to the database
	err := baseRepo.db.Ping()
	require.NoError(t, err)

	// Test that we can execute a simple query
	var count int
	err = baseRepo.db.Get(&count, "SELECT COUNT(*) FROM orders")
	require.NoError(t, err)
	assert.GreaterOrEqual(t, count, 0)

	// Test that the repository is properly initialized
	assert.NotNil(t, repo)
	assert.Equal(t, baseRepo, repo.BaseRepository)
}

func TestOrderRepository_CreateOrder_DataConsistency(t *testing.T) {
	baseRepo, cleanup := SetupTestDB(t)
	defer cleanup()

	repo := NewOrderRepository(baseRepo)

	// Create multiple orders and verify data consistency
	orders := []*models.Order{
		{Status: "pending", TotalPrice: decimal.NewFromFloat(10.00)},
		{Status: "confirmed", TotalPrice: decimal.NewFromFloat(20.00)},
		{Status: "cancelled", TotalPrice: decimal.NewFromFloat(30.00)},
	}

	for i, order := range orders {
		err := baseRepo.WithTransaction(func(tx *sqlx.Tx) error {
			return repo.CreateOrder(tx, order)
		})
		require.NoError(t, err, "Failed to create order %d", i)
	}

	// Verify all orders were created with correct data
	for i, order := range orders {
		var dbOrder struct {
			ID         int             `db:"id"`
			CreatedAt  int64           `db:"created_at"`
			Status     string          `db:"status"`
			TotalPrice decimal.Decimal `db:"total_price"`
		}
		err := baseRepo.db.Get(&dbOrder, "SELECT id, created_at, status, total_price FROM orders WHERE id = $1", order.ID)
		require.NoError(t, err, "Failed to retrieve order %d", i)

		assert.Equal(t, order.ID, dbOrder.ID)
		assert.Equal(t, order.CreatedAt, dbOrder.CreatedAt)
		assert.Equal(t, order.Status, dbOrder.Status)
		assert.True(t, order.TotalPrice.Equal(dbOrder.TotalPrice))
	}
}
