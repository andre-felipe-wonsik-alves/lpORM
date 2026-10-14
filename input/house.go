package input

type House struct {
	ID       int64  `lporm:"pk,auto"`
	ClientID int64  `lporm:"n-n=clientes.id"`
	Address  string `lporm:"size=100" form:"type=text;label=Endereço" validate:"required"`
}
