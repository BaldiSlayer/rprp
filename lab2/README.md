# Что пришлось сделать чтобы заработало

```bash
brew install pkg-config
brew install openmpi
```

### Запуск обычной версии:
```bash
go run ./cmd/normal
```

### Запуск версии с mpi:
```bash
go build -tags mpi ./cmd/parallel && mpirun -np 2 ./parallel
```