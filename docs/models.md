# models/*.go

Os arquivos contidos em `models` definem os Entidades 
participantes de uma `Session`, principal objeto da para geração 
de instâncias. A `Session` armazena em memória todas as variáveis
necessárias para a construção da instância e é definida da seguinte forma

```go
type Session struct {
	Id          int
	Cursos      map[int]Curso
	Professores map[int]Professor
	Disciplinas map[int]Disciplina
	Contratos   map[int]Contrato
	Recursos    map[int]Recurso
}
 
```

Os objetos são armazenados em mapas ao invés de arrays dinâmicos uma vez que 
possibilita remover elementos de forma mais eficiente e com seus `Id`s não mais referenciando 
posições em vetor.

Os elementos mais importantes da `Session` são os `Contrato`
que são definidos dessa forma:

```go
type Contrato struct {
    Id int
    Participantes []ParticipanteContrato
    Tipo_por_participante []string      // Tipo_por_participante[i] é o tipo do participante Participantes[i]
    Dispo       Disponibilidade
}

type ParticipanteContrato interface {
    GetId() int
    GetNome() string
    GetDispo() Disponibilidade
}
```

A interface `ParticipantecContrato` possúi os metodos essencias para o funcionamento do contrato.
`GetDispo` Possibilita a construção da `Disponibilidade` do contrato, sendo a interseção das disponibilidades de seus participantes. `GetId` quando usado em conjunto com `Contrato.Tipo_por_participante` permite encontrar o exato elemento concreto da `Session` que este participante representa.
Por exemplo tendo um participante 
```go
{
    Id: 12345,
    Nome: "artur cs",
    Dispo: [...],
}
```
E conferindo através de `Contrato.Tipo_por_participante`
que esse participante é de tipo `Professor`, basta acessar `Session.Professores[12345]` para obeter informações
mais completas. Isso é importante na serialização JSON do contrato, uma vez
que não se pode serializar um objeto acessado por polimorfismo
então ao serializar um contrato seus participantes, independente do tipo, ficam na forma do ultimo _snippet_ de código.


#### Map_to_slice

Gostaria de adiconar algumas considerações sobre os metodos `map_to_slice` e `slice_to_map` contidos em `models/session.go`.

Como mencionado antes, encontrei que é mais prático armazenar os elementos de uma `Session` em mapas e não arrays dinâmicos (em Go chamados de slices). Porém isso não é verdade 
no contexto de um cliente recebendo esses elementos para, por exemplo, exibilos em 
uma UI web. Nesse tipo de situação é mais fácil lidar com um grupo fácilmente iterável de elementos como em um array. Logo ao serializar uma `Session` para JSON 
seus mapas são convertidos para arrays (`map_to_slice`), e ao desseriizar o JSON de volta para uma `Sesison` o oposto ocorre (`slice_to_map`).




