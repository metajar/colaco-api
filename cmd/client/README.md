# ColaCo API CLI Tool

## About the Project

The ColaCo API CLI Tool is a command-line interface application designed to interact with the ColaCo vending machine service. Built with Go and utilizing the Cobra library, it offers a comprehensive set of commands for managing vending machine operations, including inventory management, pricing updates, and transaction processing.

## Features

- **Inventory Management**: Add, restock, and remove soda items from the vending machine's inventory.
- **Pricing Updates**: Update the price of existing soda items.
- **Transaction Processing**: Process purchases, handling inventory deductions and sales tracking.
- **Authentication**: Secure access with username and password authentication.
- **Easy Configuration**: Customize server URL, authentication details, and more through command-line flags.

## Getting Started

### Prerequisites

- Go 1.15 or later
- Git

### Installation

1. **Clone the GitHub repository**:
   ```bash
   git clone https://github.com/metajar/colaco-api.git
   ```
2. **Navigate to the project directory**:
   ```bash
   cd colaco-api
   ```
3. **Build the application**:
   ```bash
   go build -o colaco-cli ./cmd/client
   ```

### Accessing Help

```bash
(base) ➜  colaco-api git:(main) ✗ ./colaco-cli                       
Simple CLI client used to interact with the Vending Machine server.

Usage:
  client [command]

Available Commands:
  add-soda      Adds a new soda to the vending machine
  completion    Generate the autocompletion script for the specified shell
  delete-soda   deletes soda from the vending machine by removing the vending slot
  get-sodas     Gathers all the sodas that are in the vending slots.
  get-token     gets token from the server that can be used with other tooling such as postman.
  help          Help about any command
  purchase-soda Purchases a soda from the vending machine
  restock-soda  Restocks a specific soda in the vending machine
  update-price  updates the price of a soda

Flags:
  -h, --help              help for client
  -p, --password string   Password to use to communicate with the vending machine.
  -s, --server string     Server URL of the vending machine service. (default "http://localhost:8080")
  -t, --toggle            Help message for toggle
  -u, --username string   Username to use to communicate with the vending machine. (default "admin")

Use "client [command] --help" for more information about a command.

```




### Configuration

Configure the tool using command-line flags:

- `--server` (`-s`): Specify the server URL. Default: `http://localhost:8080`.
- `--username` (`-u`): Authentication username. Default: `admin`.
- `--password` (`-p`): Authentication password.

## Usage

Ensure that the vending machine server is up and running. Once this is done
you will be able to use the client to interact with the server.

Utilize the CLI tool to manage the vending machine:

- **View Inventory**:
  ```bash
  ./colaco-cli  get-sodas -u admin -p password
  ```
- **Add New Soda**:
  ```bash
  ./colaco-cli add-soda -u admin -p password --name "Dre.Pepper" --description "Another One" --price 1.23 --quantity 100 --calories 133 --ounces 15
  ```
- **Restock Soda**:
  
  ```bash
  ./colaco-cli restock-soda -u admin -p password --soda Pop --qty 11
  ```
- **Update Soda Price**:
  ```bash
  ./colaco-cli update-price -u admin -p password --soda Pop --price 9.93 
  ```
- **Delete Soda**:
  ```bash
  ./colaco-cli delete-soda -u admin -p password --soda "Fizz"

  ```
- **Process Purchase**:
  ```bash
  ./colaco-cli purchase-soda -u admin -p password --soda Pop --payment 1.44

  ```

## API Endpoints

The CLI tool interfaces with the following API endpoints:

- `GET /vending`: Retrieve vending machine inventory.
- `POST /soda/new`: Add a new soda item.
- `PUT /soda/restock`: Restock an existing soda item.
- `PUT /soda/price`: Update the price of a soda item.
- `DELETE /soda/{name}`: Remove a soda item from inventory.
- `POST /purchase`: Process a soda purchase.


## Contact

Project Link: [https://github.com/metajar/colaco-api](https://github.com/metajar/colaco-api)

---

