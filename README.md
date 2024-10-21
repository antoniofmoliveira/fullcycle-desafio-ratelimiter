# Desafio FullCycle Rate Limiter

Desenvolver habilidades práticas em Go através da criação de um rate limiter eficiente, capaz de controlar o volume de requisições a um serviço web com base em endereços IP ou tokens de acesso.

## Solução

Implementação de rate limiter do tipo Token Bucket.

- a partir da definição da quantidade de tokens e do tempo de refresh o limiter irá se atualizar independente da quantidade de tokens usados

- possui implementação de middleware para http.Handler

- irá logar a cada refresh a quantidade de tokens remanescentes

Implementação de servidor de teste

- aceita flag de porta a ser utilizada

- aceita flags com quantidade de tokens e refresh do Token Bucket

- aceita flags para simular erros no servidor

- aceita flags para simular latência no servidor


Sugestão para utilizar juntamente com o stress tester em várias configurações para obter a melhor configuração:

exemplo:

```bash
$ go run cmd/main.go 
2024/10/21 08:40:06 Tokens remaining: 10000
2024/10/21 08:40:06 Reset bucket at 2024-10-21 08:40:06.544855658 -0300 -03 m=+5.005055072
2024/10/21 08:40:11 Tokens remaining: 0
2024/10/21 08:40:11 Reset bucket at 2024-10-21 08:40:11.546082605 -0300 -03 m=+10.006281983
2024/10/21 08:40:16 Tokens remaining: 0
2024/10/21 08:40:16 Reset bucket at 2024-10-21 08:40:16.547630689 -0300 -03 m=+15.007830116
2024/10/21 08:40:21 Tokens remaining: 10000
2024/10/21 08:40:21 Reset bucket at 2024-10-21 08:40:21.551544029 -0300 -03 m=+20.011743460
2024/10/21 08:40:26 Tokens remaining: 10000
2024/10/21 08:40:26 Reset bucket at 2024-10-21 08:40:26.555798352 -0300 -03 m=+25.015997761
2024/10/21 08:40:31 Tokens remaining: 10000
2024/10/21 08:40:31 Reset bucket at 2024-10-21 08:40:31.560010025 -0300 -03 m=+30.020209437
2024/10/21 08:40:36 Tokens remaining: 10000
2024/10/21 08:40:36 Reset bucket at 2024-10-21 08:40:36.563334265 -0300 -03 m=+35.023533671
2024/10/21 08:40:41 Tokens remaining: 10000
2024/10/21 08:40:41 Reset bucket at 2024-10-21 08:40:41.567372185 -0300 -03 m=+40.027571604
2024/10/21 08:40:46 Tokens remaining: 10000
2024/10/21 08:40:46 Reset bucket at 2024-10-21 08:40:46.57163632 -0300 -03 m=+45.031835753
^C2024/10/21 08:40:48 server: shutting down
Server stopped
```

no stress tester:

```bash
$ go run main.go --numtests 400000 
2024/10/21 08:40:08 Tests starting...
 Running  400000  tests with interval  1  microseconds and endpoint  http://localhost:8080/hello
Min/Seg       Rate           Error        Avg Time       Net Error
4008        73,785          63,849       750.032µs               0
4009       108,537         108,537      18.999011ms              0
4010        56,399          51,349      872.80571ms              0
4011        24,524          19,661      1.725696512s             0
4012        21,793          21,793      586.862088ms             0
4013        86,177          86,177      51.574782ms              0
4014        28,785          28,785      1.417574ms               0
2024/10/21 08:40:14 Tests finished.

```

observar a quantidade na coluna Error para verificar a efetividade do rate limiter ao atingir os limites estabelecidos.

baixando o refresh rate para 1 segundo

```bash
$ go run cmd/main.go -time-frame-seconds 1
2024/10/21 08:49:07 Tokens remaining: 10000
2024/10/21 08:49:07 Reset bucket at 2024-10-21 08:49:07.682646014 -0300 -03 m=+1.001003578
2024/10/21 08:49:08 Tokens remaining: 10000
2024/10/21 08:49:08 Reset bucket at 2024-10-21 08:49:08.683524514 -0300 -03 m=+2.001882080
2024/10/21 08:49:09 Tokens remaining: 10000
2024/10/21 08:49:09 Reset bucket at 2024-10-21 08:49:09.683727759 -0300 -03 m=+3.002085328
2024/10/21 08:49:10 Tokens remaining: 10000
2024/10/21 08:49:10 Reset bucket at 2024-10-21 08:49:10.684436569 -0300 -03 m=+4.002794106
2024/10/21 08:49:11 Tokens remaining: 10000
2024/10/21 08:49:11 Reset bucket at 2024-10-21 08:49:11.684650904 -0300 -03 m=+5.003008470
2024/10/21 08:49:12 Tokens remaining: 10000
2024/10/21 08:49:12 Reset bucket at 2024-10-21 08:49:12.685408407 -0300 -03 m=+6.003765970
2024/10/21 08:49:13 Tokens remaining: 10000
2024/10/21 08:49:13 Reset bucket at 2024-10-21 08:49:13.685806363 -0300 -03 m=+7.004163953
2024/10/21 08:49:14 Tokens remaining: 10000
2024/10/21 08:49:14 Reset bucket at 2024-10-21 08:49:14.686454489 -0300 -03 m=+8.004812050
2024/10/21 08:49:15 Tokens remaining: 0
2024/10/21 08:49:15 Reset bucket at 2024-10-21 08:49:15.686681591 -0300 -03 m=+9.005039133
2024/10/21 08:49:16 Tokens remaining: 0
2024/10/21 08:49:16 Reset bucket at 2024-10-21 08:49:16.68673631 -0300 -03 m=+10.005093853
2024/10/21 08:49:17 Tokens remaining: 0
2024/10/21 08:49:17 Reset bucket at 2024-10-21 08:49:17.68740388 -0300 -03 m=+11.005761412
2024/10/21 08:49:18 Tokens remaining: 0
2024/10/21 08:49:18 Reset bucket at 2024-10-21 08:49:18.687864941 -0300 -03 m=+12.006222471
2024/10/21 08:49:19 Tokens remaining: 10000
2024/10/21 08:49:19 Reset bucket at 2024-10-21 08:49:19.688722265 -0300 -03 m=+13.007079827
2024/10/21 08:49:20 Tokens remaining: 10000
2024/10/21 08:49:20 Reset bucket at 2024-10-21 08:49:20.688827733 -0300 -03 m=+14.007185266
2024/10/21 08:49:21 Tokens remaining: 10000
2024/10/21 08:49:21 Reset bucket at 2024-10-21 08:49:21.689555312 -0300 -03 m=+15.007912892
2024/10/21 08:49:22 Tokens remaining: 10000
2024/10/21 08:49:22 Reset bucket at 2024-10-21 08:49:22.689704875 -0300 -03 m=+16.008062444
^C2024/10/21 08:49:23 server: shutting down
Server stopped
```

o stress tester agora teve 10.000 sucessos por segundo

```bash
$ go run main.go --numtests 400000 
2024/10/21 08:49:14 Tests starting...
 Running  400000  tests with interval  1  microseconds and endpoint  http://localhost:8080/hello
Min/Seg       Rate           Error        Avg Time       Net Error
4914        11,919           1,968       807.233µs               0
4915        97,725          87,772      15.82696ms               0
4916       112,907         102,969      4.373471ms               0
4917       116,336         106,396       256.732µs               0
4918        61,113          61,113        251.45µs               0
2024/10/21 08:49:18 Tests finished.
```