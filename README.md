# handshake

### Installation
##### Macbook
 - Install homebrew: `/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`
 - Install bitcore: `brew install bitcoin`
 - Clone repo locally using: `git clone git@github.com:niharika88/handshake.git`

Similar steps for Linux/Windows Machines.
### Testing
##### Macbook
 - Since bitcoin core is already installed, just run it directly
 - Run bitcore node using: `/usr/local/opt/bitcoin/bin/bitcoind`
 - Navigate to cloned repo directory (_should have Go installed locally_)
   - Run `go build` to build the code
   - Execute the test: `./handshake`

##### Testing screenshots
 - ![one](https://github.com/niharika88/handshake/blob/main/payloadstring.png)
 - ![one](https://github.com/niharika88/handshake/blob/main/payloadbytes.png)
