# LP-ORM
*Framework* criado em GO para funcionar como um `ORM` e um gerador de *widgets* `HTML` de maneira automatizada.

## Como usar
O LPORM funciona com base em marcações em Model de entrada. Assim, para cada campo/tabela criada precisamos de uma marcação específica. Exemplo de uso:

```go
type Cliente struct {
	ID     int64  `lporm:"pk,auto"`
	Nome   string `lporm:"size=100" form:"type=text;label=Nome completo" validate:"required,min=3"`
	Email  string `lporm:"unique"   form:"type=email"                   validate:"required,email"`
	Ativo  bool   `form:"type=checkbox;label=Ativo"`
	Estado string `form:"type=select;options=PR,SP,SC"`
}
```

### Marcadores existentes
- `lporm`: opções de configuração para o Banco de Dados
  - Argumentos: `pk`, `auto`, `size`, `unique`
- `form`: opções de geração do HTML
  - Argumentos: `type`, `label`, `options`

## Referências
O principal guia da implementação foi o [GORM](https://gorm.io/index.html); uma biblioteca ORM para o GO. Nos inspiramos na maneira com que ele lida com as marcações no model.