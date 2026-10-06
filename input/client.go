package input

type Cliente struct {
	ID     int64  `lporm:"pk,auto"`
	Nome   string `lporm:"size=100" form:"type=text;label=Nome completo" validate:"required,min=3"`
	Email  string `lporm:"unique"   form:"type=email"                   validate:"required,email"`
	Ativo  bool   `form:"type=checkbox;label=Ativo"`
	Estado string `form:"type=select;options=PR,SP,SC"`
}
