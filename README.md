# Implementação do Algoritmo de Kruskal em Go

[📄 Ler Artigo IEEE (PDF)](./article.pdf)

## Testes

O projeto inclui uma suíte de testes robusta para garantir a corretude e a performance do algoritmo.

### Testes Unitários e de Cenário

Verificam a funcionalidade básica e cenários específicos, como grafos desconexos, arestas paralelas e ciclos.

```bash
# Executa todos os testes
go test

# Executa os testes com mais detalhes (verbose)
go test -v
```

### Fuzz Testing

Para descobrir casos extremos e inesperados, o projeto utiliza Fuzz Testing. Ele gera automaticamente grafos aleatórios para testar os invariantes do algoritmo (ex: a MST não pode ter ciclos).

```bash
# Executa os testes de fuzzing
go test -fuzz=FuzzKruskalMST
```

### Benchmarks

Para medir o desempenho do algoritmo em grafos maiores.

```bash
# Executa os benchmarks
go test -bench=.
```
