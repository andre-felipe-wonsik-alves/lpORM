package input

type Car struct {
	ID       int64  `lporm:"pk,auto"`
	ClientID int64  `lporm:"fk=clientes.id"`
	Name     string `lporm:"size=100" form:"type=text;label=Nome completo" validate:"required,min=3"`
}
