# main.go e routes/crud.go

Esse arquivo basicamente define as principais rotas HTTP do servidor.

Em geral é um CRUD para cada item da sessão, as rotas desse crud são geradas automaticamente da seguinte forma:

### Exemplo de Add

criasse uma função add para cada tipo que permite adicionar um desses tipo à um objeto de sessão:

```go
func AddDisciplina(s *Session, d Disciplina) int {
	d.Dispo = NewDisponibilidade()

	d.Id = int(binary.BigEndian.Uint32([]byte(uuid.NewString())[:4]))
	s.Disciplinas[d.Id] = d
	return d.Id
}

```
Para esta função em geral precisa-se apenas do Id do item nesse caso `d.Id` e o atributo de `Session` que armazena itens desse tipo, no caso `Session.Disciplinas` ou `s.Disciplinas`. Esata função simplesmente define um Id aleatório para o item o adiciona à Sessão.

Essa função simples possíbilita o desenvolvimente da função `AddBuilder`

```go
func AddBuilder[T SessionItem] (slice_adder func(sess *Session, en T) int) func (http.ResponseWriter, *http.Request) {
    return func (w http.ResponseWriter, r * http.Request) {
        s := Session_or_nil(r)
        if s == nil {
            w.WriteHeader(http.StatusUnauthorized)
            w.Write([]byte("No session or Session expired"))
            return 
        }

        var e T 

        if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte("malformed body 1: " + err.Error()))
            return
        }
        fmt.Printf("Addded  %v\n", e)

        w.Write([]byte(fmt.Sprint(slice_adder(s, e))))
    }
}

```

Essa função recebe como parametro a função usada para adicionar
o um item de tipo `T` e retorna uma função handler http que processa pedidos HTTP que adiciona itens do tipo `T` à Sessão `s`

Isso permite que para adcionar uma nova rota de Adição para um novo tipo de entidade, basta definir o tipo em código assim como sua função "adicionadora" para que seguindo o modelo criar uma nova rota.

Os outrs item do CRUD (além do Create) são implementados de forma semelhante.

Alguns tipo de entidade como por exemplo o `Turma` ou `Etapa` não tem o próprio CRUD ou não tem CRUD completo. Isso se deve pois esses dois elementes específicos são contidos em `Curso`, logo, para por exemplo adcionar uma nova eta `Etapa` para a sessão deve-se obter um objeto `JSON` do curso (`GET`), editá-lo para conter uma nova  etapa e então salvar a edição com um pedido `SET`, da mesma forma para editar ou remover uma etapa existente.
