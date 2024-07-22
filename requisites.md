<h1 align="center"> Desafio Back-end #2</h1>
<p align="center">
  <img alt="Frete Rapido Logo"src="https://i.imgur.com/5DNdxeP.png" width="30%" height="auto">
</p>

**Tecnologias que você pode utilizar:**

*   Linguagem: **Golang**
    *   _caso não se sinta segura em fazer em Golang, aceitamos NodeJS e/ou PHP._ 
        
*   Banco de dados: a sua escolha

**Requisitos:**

*   Aplicação conteinerizada: sua resolução (e todas as suas dependências também) devem ser executadas por meio de um container do Docker.
    
*   Validação dos dados de entrada;
    
*   Mensagens de erros claras;  
    
*   Documentação clara de como executar o código.

**O que seria muito legal se você utilizar/implementar:**

*   Boas práticas de programação (código limpo e bem estruturado);
*   Aplicação de TDD;

**Objetivo geral:**

*   Desenvolver uma API Rest para consultas externas e devolver apenas valores esperados.  
    

* * *

  

**Rota 1: [POST] .../quote**

* Objetivo: Criar uma rota para receber dados de entrada e realizar uma cotação fictícia com a API da Frete Rápido (os valores e transportadoras retornadas não são reais);

* Entrada: Esperar receber um JSON de entrada conforme exemplo abaixo:
```json
{
    "recipient":{
      "address":{
          "zipcode":"01311000"
      }
    },
    "volumes":[
      {
          "category":7,
          "amount":1,
          "unitary_weight":5,
          "price":349,
          "sku":"abc-teste-123",
          "height":0.2,
          "width":0.2,
          "length":0.2
      },
      {
          "category":7,
          "amount":2,
          "unitary_weight":4,
          "price":556,
          "sku":"abc-teste-527",
          "height":0.4,
          "width":0.6,
          "length":0.15
      }
    ]
}
``` 

* Processo: Utilizar os dados de entrada para complementar os dados necessários para consumir a API da Frete Rápido no método [“Cotações de Frete v3”](https://dev.freterapido.com/ecommerce/cotacao_v3/). Ao complementar a estrutura padrão obrigatória para a requisição na Frete Rápido, realizar a cotação.

  
Os retornos das cotações devem ser gravadas em um banco de dados para que sejam consumidas na rota 2 desse desafio.  

* Retorno esperado:
```json
{
    "carrier":[
      {
          "name":"EXPRESSO FR",
          "service":"Rodoviário",
          "deadline":"3",
          "price":17
      },
      {
          "name":"Correios",
          "service":"SEDEX",
          "deadline":1,
          "price":20.99
      }
    ]
}
```


* Observação: Para consumir a API da Frete Rápido, você vai precisar dos dados obrigatórios:
  *   CNPJ Remetente: **25.438.296/0001-58** (apenas números) 
      *   Usar o mesmo CNPJ para "shipper.registered_number" e "dispatchers.registered_number"
      *   Token de autenticação: **1d52a9b6b78cf07b08586152459a5c90**
  *    Código Plataforma: **5AKVkHqCn**
  *   Cep: **29161-376** _(dispatchers[*].zipcode)_

_"unitary\_price" deve ser informado._ 

  

**Rota 2: [GET] .../metrics?last\_quotes={?}**

  

● Objetivo: consultar métricas das cotações armazenadas no seu banco de dados, permitindo receber um parâmetro “last\_quotes”, não obrigatório, informando a quantidade de cotações (ordem decrescente);

  

● Processo: ao consultar os retornos gravados na base de dados, retornar as

seguintes métricas dos resultados:

○ Quantidade de resultados por transportadora;

○ Total de “preco\_frete” por transportadora; (final\_price na API)

○ Média de “preco\_frete” por transportadora; (final\_price na API)

○ O frete mais barato geral;

○ O frete mais caro geral;


Dúvidas sobre a nossa API, leia nossa documentação: [https://dev.freterapido.com/](https://dev.freterapido.com/)

> Seja um parâmetro de qualidade. Algumas pessoas não estão acostumadas a um ambiente onde a excelência é esperada. - Steve Jobs.

Mãos à obra e sucesso!