# LP-ORM
O **LP-ORM** é um mini-framework web desenvolvido em Go com o objetivo de automatizar a criação de tabelas de banco de dados (por meio de um `ORM`) e a geração de formulários e *widgets* `HTML`. 

O projeto foi desenvolvido para a disciplina **Aspectos de Linguagens de Programação (BCC6002)** da Universidade Tecnológica Federal do Paraná (UTFPR), pelos alunos:

- André Felipe Wonsik Alves
- Eber Felipe Barrotti Louback
- Erick Molina Gehring
- Gustavo Martins França

## Como usar
O LPORM funciona com base em marcações em Model de entrada. Assim, para cada campo/tabela criada precisamos de uma marcação específica. Exemplo de uso:

```go
type Cliente struct {
	ID          int64  `lporm:"pk,auto"`
	Nome        string `lporm:"size=100" form:"type=text;label=Nome completo" validate:"required,min=3"`
	Email       string `lporm:"unique"   form:"type=email"                   validate:"required,email"`
	Genero      string `form:"type=radio;options=M:Masculino,F:Feminino;label=Gênero"`
	Ativo       bool   `form:"type=checkbox;label=Ativo"`
	Estado      string `form:"type=select;options=PR,SP,SC;label=Estado"`
	Observacoes string `form:"type=textarea;label=Observações"`
}
```

### Marcadores existentes
- `lporm`: opções de configuração para o Banco de Dados
  - Argumentos: `pk`, `auto`, `size`, `unique`, `fk`, `n-n`
- `form`: opções de geração do HTML
  - Argumentos: `type`, `label`, `options`

## Como funciona

Os `structs` Go representam as relações que serão cadastradas no banco de dados. O lpORM lê os campos e suas marcações (tags) para obter os metadados do modelo e definir a estrutura da tabela, incluindo chaves primárias, chaves estrangeiras, valores únicos e campos obrigatórios. Uma chave estrangeira é declarada com `fk=tabela.coluna`, que cria a constraint correspondente no banco. A marcação `n-n=tabela.coluna` declara uma associação muitos-para-muitos e gera uma tabela de associação com chaves estrangeiras para os dois modelos. As opções `form` e `validate` são usadas na configuração e validação dos campos dos formulários HTML.

Para integrar o lpORM a uma aplicação, abra uma conexão com o banco de dados e passe-a para `assembler.New`. Em seguida, registre cada `struct` com `assembler.Register`, como no exemplo em [`example/api`](example/api). O `Assembler` reúne as rotas HTTP e os modelos registrados.

Ao registrar um modelo, o lpORM cria a tabela correspondente no banco, prepara um formulário HTML para inserir novos registros e associa a ele uma rota. Ao enviar o formulário, a requisição é tratada pelo handler dessa rota e os valores são inseridos na tabela correspondente. A página inicial (`/`) apresenta links para os formulários registrados. O exemplo inicia o servidor HTTP e disponibiliza essa página e os formulários por meio dele.

## Escopo e limitações

- O lpORM ainda não oferece uma CLI para interação pelo terminal, a integração demonstrada é feita por uma aplicação Go com servidor HTTP.
- Atualmente, apenas PostgreSQL é suportado.
- Não há suporte dedicado para declarar e gerenciar associações um-para-um (1:1). É possível definir uma chave estrangeira com `fk` e uma constraint de unicidade com `unique`, mas o lpORM não implementa o tratamento automático dessa relação.
- A inserção não usa `RETURNING` nem retorna a tupla criada, a operação executa o `INSERT` e informa se houve erro.

## Referências
O principal guia da implementação foi o [GORM](https://gorm.io/index.html); uma biblioteca ORM para o GO. Nos inspiramos na maneira com que ele lida com as marcações no model.