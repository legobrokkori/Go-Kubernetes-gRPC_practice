package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/legobrokkori/go-kubernetes-grpc_practice/util"
	"github.com/stretchr/testify/require"
)

func createRondomTransfer(t *testing.T) Transfer {
	account1 := createRondomAccount(t)
	account2 := createRondomAccount(t)
	arg := CreateTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, account1.ID, transfer.FromAccountID)
	require.Equal(t, account2.ID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	createRondomTransfer(t)
}

func TestGetTransfer(t *testing.T) {
	transfer1 := createRondomTransfer(t)
	tansfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tansfer2)

	require.Equal(t, transfer1.ID, tansfer2.ID)
	require.Equal(t, transfer1.FromAccountID, tansfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, tansfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, tansfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, tansfer2.CreatedAt, time.Second)
}

func TestUpdatesTransfer(t *testing.T) {
	transfer1 := createRondomTransfer(t)

	arg := UpdateTransferParams{
		ID:     transfer1.ID,
		Amount: util.RandomMoney(),
	}
	tansfer2, err := testQueries.UpdateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tansfer2)

	require.Equal(t, transfer1.ID, tansfer2.ID)
	require.Equal(t, transfer1.FromAccountID, tansfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, tansfer2.ToAccountID)
	require.NotEqual(t, transfer1.Amount, tansfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, tansfer2.CreatedAt, time.Second)
}

func TestDeleteTransfer(t *testing.T) {
	transfer1 := createRondomTransfer(t)

	err := testQueries.DeleteEntry(context.Background(), transfer1.ID)
	require.NoError(t, err)
	entry2, err := testQueries.GetEntry(context.Background(), transfer1.ID)
	require.Error(t, err)
	require.Error(t, err, sql.ErrNoRows.Error())
	require.Empty(t, entry2)
}

func TestListTransfers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRondomTransfer(t)
	}

	arg := ListTransfersParams{
		Limit:  5,
		Offset: 5,
	}
	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 5)

	for _, transer := range transfers {
		require.NotEmpty(t, transer)
	}
}
