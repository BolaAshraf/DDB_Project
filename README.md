# Building a Distributed Database Using Go

## 📚 Objective

Design and implement a basic distributed database system using the Go programming language. This project introduces core concepts of distributed systems, including data replication and fault tolerance.

## 💡 Key Features

* A distributed database system with a master-slave architecture.
* Supports dynamic creation of databases and tables.
* Master node can perform CRUD operations, and slaves can execute queries (excluding DROP).
* Data replication across multiple nodes for fault tolerance.
* Communication between nodes using TCP.

## 🛠 Architecture
                +-------------------------+
                |      Master Node        |
                |------------------------ |
                | - DB Write Access       |
                | - Broadcast to Slave    |
                +-------------------------+
                            |
                            v
            -----------------------------------
           |                                   |
           v                                   v
+-------------------------+     +-------------------------+
|     Slave Node 1        |     |     Slave Node 2        |
|-------------------------|     |-------------------------|
| - Read-only DB          |     | - Read-only DB          |
| - Listen for Replication|     | - Listen for Replication|
+-------------------------+     +-------------------------+
1. Master Node:

   * Manages the creation of databases and tables.
   * Executes queries and replicates data to slave nodes.
   * Uses TCP for communication with slaves.

2. Slave Nodes:

   * Independently store data and execute queries.
   * Receive replicated data from the master.
   * Restricted from executing database creation or deletion operations.

## ⚙️ Installation

1. Clone the repository:

   
   git clone https://github.com/BolaAshraf/DDB_Project.git
   cd distributed-mysql
   

2. Install Go and MySQL on your system.

3. Update the configuration (IP addresses, ports, credentials) in the master and slave files as needed.

4. Run the Master Node:

   
   go run master.go
   

5. Run the Slave Node:

   
   go run slave.go
   

## 🚀 Usage

* Start the master node first, followed by the slave nodes.
* Use the master console to execute queries.
* The master will automatically replicate data to all connected slave nodes.

### Example Commands:

* To create a new database:

  
  CREATE DATABASE testDB;
  
* To insert a record:

  
  INSERT INTO users (id, name) VALUES (1, 'Jessy');
  
* To view records from slaves:

  
  SELECT * FROM users;
  

## 📂 Project Structure

* master.go: Main file for the master node.

* slave.go: Main file for the slave nodes.

* pkg/db/database.go: Database management, handling SQL connections and query execution.

* pkg/network/client.go: Handles TCP communication to send messages.

* pkg/network/server.go: Handles TCP communication to receive messages.

* pkg/network: TCP communication package.

* master.go: Main file for the master node.

* slave.go: Main file for the slave nodes.

* pkg/db: Database management package.

* pkg/network: TCP communication package.

## 📝 License

This project is licensed under the MIT License.

## 👥 Authors

* Jessy - Project Lead and Developer

Feel free to contribute or report issues on the GitHub repository.
