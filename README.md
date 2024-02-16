# VIDEO DEMO
**Link: https://youtu.be/I3eV-X-P5gA**

# ABOUT
Building Blockchain Proof Of Work from scratch with Golang

# Build the program
Firstly, we need to build two executable files: one for the blockchain node server and another for the client.

To build these two files, follow these steps:

From the BlockchainPow directory, enter the following command to build the executable file for the blockchain_server package:

**go build -o blockchain_node ./blockchain_server**

Still in the BlockchainPow directory, enter the following command to build the executable file for the blockchain_service package:

**go build -o client ./blockchain_service**

Now, we have two executable files to run the demo program.

In this context, the team has developed a user interface program that runs in the terminal, connecting to the nearest node to query blockchain-related information.

To ensure a smooth running process, please do following these steps to initialize the program correctly (by default, only allows peers to connect to nodes with IP ranges from 127.0.0.1 to 127.0.0.2 and ports from 5000 to 5003).

## Run the node servers 

Initialize two blockchain node servers running on different ports (using two separate terminals):

The default port for the first blockchain node server is 5000. Enter the following command: **./blockchain_node**

For the second blockchain node server running on port 5001, use the following command: **./blockchain_node -port 5001**

Wait for the blockchain node servers to connect with each other. Each blockchain node server, upon successfully connecting to another, will add the connected node to its peers node array (which is an attribute of the blockchain).

In the example above, we have two nodes running on ports 5000 and 5001. Therefore, if the connection is successful, the connectivity status will look like the following (observe the terminal interface):

**Node 1** 

For the blockchain node server running on port 5000, it needs to peer with the blockchain node server running on port 5001 (wait for the following array to be printed in the terminal)

<img width="468" alt="image" src="https://github.com/thoaikhoa14402/BlockchainPow/assets/81000230/4e759378-b49a-4f6a-9946-a9065e548115">

**Node 2**

For the blockchain node server running on port 5001, it needs to peer with the blockchain node server running on port 5000 (wait for the following array to be printed in the terminal)

<img width="468" alt="image" src="https://github.com/thoaikhoa14402/BlockchainPow/assets/81000230/d2ff2a63-f8e0-4ca1-b474-3b8ef286d187">

If both terminals print the array as shown above, it means we have successfully connected. **If the connection is not established, please reinitialize step 1**

## Run the client application

The first client program, by default, runs on port 8080 and connects to the node with the IP address **http://127.0.0.1:5000** . Therefore, you can simply enter the following command:

**./client**

The second client program is allowed to specify the port and IP address of the blockchain node it wants to connect to using the flags 'port' and 'gateway'. For example:

**./client -port 8081 -gateway http://127.0.0.1:5001**

If the connection is successful, both client programs will have interfaces like the following:

<img width="468" alt="image" src="https://github.com/thoaikhoa14402/BlockchainPow/assets/81000230/4b53dcdc-215b-4d03-a29e-b66025395c99">




## Note

- To run the program, users need to first execute function number 9 to request the blockchain node server to automatically mine new blocks.

- Next, users must create a wallet using function number 1 and wait for a moment for the system to transfer 10 coins to the user's account, this initial transaction will be included in a block (verify using function number 4)

- Once the user's balance has been updated to 10 coins (use function number 10 to reload the program), you can start utilizing the system's functionalities.


## Endpoint list
Access the APIs of each node server to query data

**1. /chain:** To get all blocks of blockchain

**2. /transactions:** To get all transactions of blockchain

**3. /mine:** To mine a new block from transaction pool

**4. /mine/start:** To automatically mine a new block from transaction pool after 20 seconds

**5. /consensus:** To reach consensus between nodes in blockchain (handle conflicts between blockchain). This endpoint will be automaticall hit after the mining process done -> Check it by review the blockchain (/chain) after mining process done

**For example:** 

**Node 1:** localhost:5000/chain, localhost:5000/transactions, ...

**Node 2:** localhost:5001/chain, localhost:5001/transactions, ...

**Node 3:** localhost:5002/chain, localhost:5002/transactions, ...

