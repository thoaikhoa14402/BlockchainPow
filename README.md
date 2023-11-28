# BlockchainPow
Building Blockchain PoW from scratch with Golang

## STEP 1: Run the node servers 
Run the node servers on ports 5000, 5001, 5002 and wait for them to peer with each other. 

**For example:**
**Node 1** (Node 1 running on port 5000 needs to peer with nodes running on ports 5001 and 5002). 

(wait for the following array to be printed in the terminal)

<img width="637" alt="image" src="https://github.com/thoaikhoa14402/BlockchainPow/assets/81000230/7478b59e-ddb4-424b-b800-cfec92227d95">

**Node 2**

<img width="641" alt="image" src="https://github.com/thoaikhoa14402/BlockchainPow/assets/81000230/a881599a-54da-4c2e-b8b9-17b738ac0311">

**Node 3**

<img width="642" alt="image" src="https://github.com/thoaikhoa14402/BlockchainPow/assets/81000230/52f6efe8-0f37-4d68-9ca6-ea73cc2b1924">

## STEP 2: Start wallet server
Afterwards, start the wallet server and access it at localhost:8080 to perform send/receive transaction 

## STEP 3: Access the API to query data
Access the APIs of each node server to query data

**Endpoint list:**
**1. /chain:** To get all blocks of blockchain

**2. /transactions:** To get all transactions of blockchain

**3. /mine:** To mine a new block from transaction pool

**4. /mine/start:** To automatically mine a new block from transaction pool after 20 seconds

**5. /consensus:** To reach consensus between nodes in blockchain (handle conflicts between blockchain)

**For example:** 
**Node 1:** localhost:5000/chain, localhost:5001/transactions, ...

**Node 2:** localhost:5001/chain, localhost:5001/transactions, ...

**Node 3:** localhost:5001/chain, localhost:5001/transactions, ...

