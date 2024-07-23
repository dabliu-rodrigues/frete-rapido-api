# ☁️🚚📦 Frete Rápido - Teste técnico!

<p align="center">
  <img alt="Frete Rapido Logo"src="https://i.imgur.com/5DNdxeP.png" width="40%" height="auto">
</p>

<p align="center">
  <a href="#-projeto">Projeto</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#-tecnologias">Tecnologias</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#-fundamentos-e-estratégias-abordadas">Fundamentos e estratégias</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#-documentação">Documentação</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#-clonando-e-executando">Clonando e executando</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#️-rotas">Rotas</a>&nbsp;&nbsp;&nbsp;|&nbsp;&nbsp;&nbsp;
  <a href="#-licença">Licença</a>
</p>

## 📌 Projeto
<b>Esta aplicação foi desenvolvida entre os dias 16 à 23 de julho de 2024, com base no teste técnico fornecido pela</b> [@Frete Rapido](https://freterapido.com.br/)!<br>

Para ficar dentro de todas especificações e requisitos, [clique aqui](requisites.md)

O objetivo desta aplicação é construir uma API HTTPs que utiliza a API Frete Rapido para realizar simulações de cota de frete, além de persistir os dados no banco de dados para a retidada de métricas.

## 👩‍💻 Tecnologias

Esse projeto foi desenvolvido utilizando as seguintes tecnologias:

- [Golang](https://go.dev/) - Linguagem escolhida pela sua robustez no processamento de dados de forma concorrente.
- [Go-chi](https://github.com/go-chi/chi) - Biblioteca super leve e performática para construção de servidores HTTP utilizando handlers nativos do Go!
- [Docker](https://www.docker.com/) - Criação de imagens e conteiners para melhor orquestração e execução do aplicativo em especificos contextos.
- [MongoDB](https://www.mongodb.com/) - Banco da dados não relacional para persistência dos dados.
- [Swagger](https://swagger.io/) - Para construção da documentação da API.

## 👨‍🏫 Fundamentos e estratégias abordadas

Esta API foi desenvolvida com o propósito de <b>entrega performática e escalável</b>.<br>

<b>Neste projeto, foram abordados patterns de código, entre eles, estão:</b>
- Single responsibility principle (SPR);
- Dependecy inversion principle (DIP);
- Keep it simple, Silly (KISS);
- You aren’t gonna need it (YAGNI);
- Interface segregation principle (ISP);
- Distroless Docker Image to a maximum compressed size

## 📚 Documentação
A documentação dessa API foi construída utilizando um toolset famoso para a construção de documentações, chamado de [Swagger](https://swagger.io/).<br>
<b>Para conseguir acessar a mesma, basta se redirecionar para a rota `/docs` depois de iniciar o servidor.</b>

## 📥 Clonando e executando
Para conseguir executar o projeto sem nenhuma interferência, certifique-se de ter os requisitos mínimos:<br>

- [Golang](https://go.dev/)
- [Git](https://git-scm.com/)
- [Docker](https://www.docker.com/)
- [Docker compose](https://docs.docker.com/compose/)
<br>

<b>Passo a passo:</b>

1. Clone o repositório localmente usando o seguinte comando no seu terminal de preferência:
```shell
    git clone https://github.com/jsGolden/frete-rapido-api    # Clonar repositório
    cd frete-rapido-api                                       # Entrar no repositório clonado
```

2. Para que seja possível a configuração das variavéis de ambiente, será necessário criar um arquivo .env utilizando o arquivo [.env.example](.env.example) como base
```shell
  cp .env.example .env    # caso esteja no Linux
            -- ou --
  copy .env.example .env  # caso esteja no Windows

```

3. Suba os containers (API e MongoDB) para a integração total dos serviços
```shell
  docker-compose up

  # é possível também parar os containers assim que quiser usando:
  # docker-compose down
```

4. Pronto! Por padrão, o seu servidor estará rodando na URI http://localhost:8080

## 🛣️ Rotas

- **[GET]** /docs
  - Responsável por renderizar a documentação Swagger
- **[POST]** /quote
  - Responsável por simular a cotação com a FreteRapido e salvar o dado no banco
  Exemplo de requisição cURL:
  ```shell
    curl --request POST \
      --url http://localhost:8080/quote \
      --header 'Content-Type: application/json' \
      --header 'User-Agent: insomnia/9.3.2' \
      --data '{
      "recipient": {
        "address": {
          "zipcode": "01311000"
        }
      },
      "volumes": [
        {
          "category": -7,
          "amount": 1,
          "unitary_weight": 5,
          "price": 349,
          "sku": "abc-teste-123",
          "height": 0.2,
          "width": 0.2,
          "length": 0.2
        },
        {
          "category": 7,
          "amount": 2,
          "unitary_weight": 4,
          "price": 556,
          "sku": "abc-teste-527",
          "height": 0.4,
          "width": 0.6,
          "length": 0.15
        }
      ]
    }'
  ```
  Exemplo de resposta:
  ```json
  [
    {
      "name": "JADLOG",
      "service": ".PACKAGE",
      "deadline": 3,
      "price": 35.99
    },
    {
      "name": "AZUL CARGO",
      "service": "Convencional",
      "deadline": 2,
      "price": 43.56
    },
    {
      "name": "PRESSA FR (TESTE)",
      "service": "Normal",
      "deadline": 0,
      "price": 60.74
    },
    {
      "name": "BTU BRASPRESS",
      "service": "Normal",
      "deadline": 5,
      "price": 93.35
    }
  ]
  ```
- **[POST]** /metrics
  - Responsável por utilizar os dados persistidos para calcular métricas
  Exemplo de requisição cURL:
  ```shell
  curl --request GET \
    --url http://localhost:8080/metrics \
    --header 'Content-Type: application/json' \
    --header 'User-Agent: insomnia/9.3.2'
  ```
  Exemplo de resposta:
  ```json
  {
    "cheapest_quote": 251.93,
    "most_expensive_quote": 1010.58,
    "services": [
      {
        "average_price": 35.99,
        "carrier": "JADLOG",
        "count": 7,
        "total_price": 251.93
      },
      {
        "average_price": 65.94,
        "carrier": "PRESSA FR (TESTE)",
        "count": 9,
        "total_price": 593.5
      },
      {
        "average_price": 86.47,
        "carrier": "AZUL CARGO",
        "count": 9,
        "total_price": 778.31
      },
      {
        "average_price": 112.28,
        "carrier": "BTU BRASPRESS",
        "count": 9,
        "total_price": 1010.58
      }
    ]
  }
  ```

## 📑 Licença
Este projeto está sobre a licença MIT.

<hr>
<p align="center">Desenvolvido com 💜 por Wagner Rodrigues</p>
