package service

type Balance interface {
	AddMoney(idUser int, moneyAmount float32) (float32, error)
	Reserve(idUser, idService, idOrder int, orderCost float32) (float32, error)
	DealSuccess(idUser, idService, idOrder int, moneyAmount float32) error
	GetBalance(idUser int) (float32, error)
	CancelOrder(idUser, idService, idOrder int, moneyAmount float32) error
}
