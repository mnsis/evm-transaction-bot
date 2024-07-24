### Step-by-Step Guide for Setting Up and Running the Script

#### Prerequisites
Ensure that you have the following installed on your machine:
- Go (Golang) programming language
- Git

#### Installing Go
1. Download and install Go:
   ```sh
   wget https://dl.google.com/go/go1.16.5.linux-amd64.tar.gz
   sudo tar -xvf go1.16.5.linux-amd64.tar.gz
   sudo mv go /usr/local
   ```

2. Update your PATH environment variable:
   ```sh
   echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
   source ~/.profile
   ```

3. Verify the installation:
   ```sh
   go version
   ```

#### Cloning the Repository
1. Clone the repository from GitHub:
   ```sh
   git clone https://github.com/yourusername/evm-transactions-monitor.git
   cd evm-transactions-monitor
   ```

#### Setting Up the Project
1. Initialize a new Go module:
   ```sh
   go mod init evm-transactions-monitor
   ```

2. Install the required Go packages:
   ```sh
   go get github.com/ethereum/go-ethereum
   ```

3. Create the main.go file and copy the provided script into it:
   ```sh
   nano main.go
   ```

#### Running the Script
1. Build and run the Go program:
   ```sh
   go run main.go
   ```

#### Script Explanation

**Features:**
- The script monitors logs from a specified service (`stationd`) and sends transactions based on specific log entries.
- The script sends initial 25 transactions and then sends additional 56 transactions whenever a "Generating New unverified pods" log entry is detected.

**Usage:**
1. Run the script and enter the required details:
   ```sh
   go run main.go
   ```
   - Enter your EVM Private Key when prompted.
   - Enter the RPC URL when prompted.

**Functional Breakdown:**
- **`createWalletFromPrivateKey`:** This function creates a wallet from a given private key.
- **`sendTransactionsSequentially`:** This function sends a specified number of transactions sequentially.
- **`sendTx`:** This function sends a single transaction.
- **`monitorStationLogs`:** This function monitors logs from the `stationd` service and triggers transactions based on specific log entries.

**Dependencies:**
- **`github.com/ethereum/go-ethereum`:** This package is used to interact with Ethereum nodes.

**Advantages:**
- Automates the process of sending transactions based on specific log entries.
- Uses Go's concurrency features to monitor logs and send transactions efficiently.
- Includes error handling to manage common issues like insufficient funds and network errors.

**Example Usage:**
1. Ensure the service `stationd` is running and generating logs.
2. Run the script:
   ```sh
   go run main.go
   ```
3. Enter the required details when prompted.

4. Monitor the console output to see the transaction statuses.

5. The script will keep running, monitoring the logs and sending transactions as needed.

### Uploading to GitHub

1. Initialize a new Git repository:
   ```sh
   git init
   ```

2. Add the files to the repository:
   ```sh
   git add main.go
   ```

3. Commit the changes:
   ```sh
   git commit -m "Initial commit"
   ```

4. Create a new repository on GitHub.

5. Add the remote repository and push the changes:
   ```sh
   git remote add origin https://github.com/mnsis/evm-transactions-monitor.git
   git push -u origin master
   ```
