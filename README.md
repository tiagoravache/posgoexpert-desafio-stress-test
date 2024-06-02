# goexpert-desafio-stress-test
Resposta para o desafio de Stress Test do curso de pós graduação Go Expert.

Para realizar o build da imagem docker do projeto, execute o seguinte comando:
```shell
docker build -t <nome_da_imagem> -f Dockerfile .
```

Para subir o container da imagem do projeto, execute o comando abaixo:
```shell
docker run <nome_da_imagem> -url=<url_de_preferencia> -concurrency=<quantidade_de_requests_simultaneos> -requests=<quantidade_total_de_requests>
``` 