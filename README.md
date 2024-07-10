# bear_robot test
계좌 생성
```bash
// 계좌 생성
curl -X POST http://localhost:8080/create -d '{"account":"myaccount", "password":"1234"}'

// 계좌 생성과 동시에 입금
curl -X POST http://localhost:8080/create -d '{"account":"myaccount", "password":"1234", "amount":1000}'
```

계좌 조회
```bash
curl -X POST http://localhost:8080/history -d '{"account":"myaccount", "password":"1234"}'
```

입금
```bash
curl -X POST http://localhost:8080/deposit -d '{"account":"myaccount", "password":"1234", "amount":1000}'
```

출금
```bash
curl -X POST http://localhost:8080/withdraw -d '{"account":"myaccount", "password":"1234", "amount":1000}'
```
