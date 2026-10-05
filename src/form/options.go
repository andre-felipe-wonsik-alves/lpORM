package form

type Option func(*Form)

func WithAction(action string) Option {
	return func(f *Form) { f.action = action }
}

func WithMethod(method string) Option {
	return func(f *Form) { f.method = method }
}

func WithSubmitLabel(label string) Option {
	return func(f *Form) { f.submitLabel = label }
}
