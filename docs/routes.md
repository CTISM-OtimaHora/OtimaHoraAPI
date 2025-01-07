# contratos.go e session.go

Ambos estes arquivos servem para tipos que não se enquadram em no sistema definido
em `docs/main.md`. `Session` pois é o tipo central desta API onde todos os outros são armazenados
tem apenas funcionalidades para ser criada (uma do zero e outra a partir de um arquivo JSON)
assim como uma para obter uma sessão existente a partir de um cookie

`Contrato` também contem copias de itens de outros tipos e logo não pode ser facilmente serializado em JSON, uma vez que sua estrutura emprega polimorfismo dos participantes que devem se adequar à interface `ParticipanteContrato`. Isso faz com que para serializar um `Contrato` para JSON é necessário um certo malabarismo com seus participantes de forma que seja possível reconstruir o mesmo contrato a partir do resultado da serialização (importanto para salvar as instâncias em arquivos JSON). Mais informações dessa implementação estão em `models/contrato.go`
