# distributed-tracing-and-span

## Exemplo
```json
http://localhost:8081/zipcode
Body: { "cep": "10010100" }

```
```json
{
  "temp_C": 29.1,
  "temp_F": 84.4,
  "temp_K": 302.1,
  "city": "Belém"
}
```
## Dependências

[weather api key](https://www.weatherapi.com/)

## Execução

```bash
git clone https://github.com/crnvl96/distributed-tracing-and-span.git
cd distributed-tracing-and-span
cd service_b
go mod download
cp .env.example .env
cd ../service_a
go mod download
cd ..
docker compose up -d
```
