package input

type Cliente struct {
	ID          int64  `lporm:"pk,auto"`
	Nome        string `lporm:"size=100" form:"type=text;label=Nome completo" validate:"required,min=3"`
	Email       string `lporm:"unique"   form:"type=email"                   validate:"required,email"`
	Genero      string `form:"type=radio;options=M:Masculino,F:Feminino;label=Gênero"`
	Ativo       bool   `form:"type=checkbox;label=Ativo"`
	Estado      string `form:"type=select;options=PR,SP,SC;label=Estado"`
	Observacoes string `form:"type=textarea;label=Observações"`
}
