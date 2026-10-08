package input

type Job struct {
	ID         int64  `lporm:"pk,auto"`
	Name       string `lporm:"size=100" form:"type=text;label=Nome completo" validate:"required,min=3"`
	DaysOfWork int
}
