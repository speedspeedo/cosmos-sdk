# This is not executable, but helper functions for the other scripts

# XXX XXX XXX XXX XXX
# The following global variables must be defined before calling common functions: 
# SERVER_EXE=foobar       # Server binary name
# CLIENT_EXE=foobarcli    # Client binary name
# ACCOUNTS=(foo bar)      # List of accounts for initialization 
# RICH=${ACCOUNTS[0]}     # Account to assign genesis balance

prepareClient() {
  echo "Preparing client keys..."
  ${CLIENT_EXE} reset_all
  assertTrue $?

  for i in "${!ACCOUNTS[@]}"; do
      newKey ${ACCOUNTS[$i]}
  done
}

# XXX Ex Usage1: initServer $ROOTDIR $CHAINID
# XXX Ex Usage2: initServer $ROOTDIR $CHAINID $PORTPREFIX
# Desc: Grabs the Rich account and gives it all genesis money
#       port-prefix default is 4665{6,7,8}
initServer() {
  echo "Setting up genesis..."
  SERVE_DIR=$1/server
  assertNotNull "no chain" $2
  CHAIN=$2
  SERVER_LOG=$1/${SERVER_EXE}.log

  GENKEY=$(${CLIENT_EXE} keys get ${RICH} | awk '{print $2}')
  ${SERVER_EXE} init --chain-id $CHAIN $GENKEY --home=$SERVE_DIR >>$SERVER_LOG

  # optionally set the port
  if [ -n "$3" ]; then
    echo "setting port $3"
    sed -ie "s/4665/$3/" $SERVE_DIR/config.toml
  fi

  echo "Starting ${SERVER_EXE} server..."
  ${SERVER_EXE} start --home=$SERVE_DIR >>$SERVER_LOG 2>&1 &
  sleep 5
  PID_SERVER=$!
  disown
  if ! ps $PID_SERVER >/dev/null; then
    echo "**FAILED**"
    # cat $SERVER_LOG
    # return 1
  fi
}

# XXX Ex Usage1: initClient $CHAINID
# XXX Ex Usage2: initClient $CHAINID $PORTPREFIX
# Desc: Initialize the client program
#       port-prefix default is 46657
initClient() {
  echo "Attaching ${CLIENT_EXE} client..."
  PORT=${2:-46657}
  # hard-code the expected validator hash
  ${CLIENT_EXE} init --chain-id=$1 --node=tcp://localhost:${PORT} --valhash=EB168E17E45BAEB194D4C79067FFECF345C64DE6
  assertTrue "initialized light-client" $?
}

# XXX Ex Usage1: newKey $NAME
# XXX Ex Usage2: newKey $NAME $PASSWORD
# Desc: Generates key for given username and password
newKey(){
  assertNotNull "keyname required" "$1"
  KEYPASS=${2:-qwertyuiop}
  (echo $KEYPASS; echo $KEYPASS) | ${CLIENT_EXE} keys new $1 >/dev/null 2>/dev/null
  assertTrue "created $1" $?
  assertTrue "$1 doesn't exist" "${CLIENT_EXE} keys get $1"
}

# XXX Ex Usage: getAddr $NAME
# Desc: Gets the address for a key name
getAddr() {
  assertNotNull "keyname required" "$1"
  RAW=$(${CLIENT_EXE} keys get $1)
  assertTrue "no key for $1" $?
  # print the addr
  echo $RAW | cut -d' ' -f2
}

# Desc: Assumes just one coin, checks the balance of first coin in any case
checkAccount() {
  # make sure sender goes down
  ACCT=$(${CLIENT_EXE} query account $1)
  assertTrue "must have genesis account" $?
  assertEquals "proper sequence" "$2" $(echo $ACCT | jq .data.sequence)
  assertEquals "proper money" "$3" $(echo $ACCT | jq .data.coins[0].amount)
  return $?
}

# XXX Ex Usage: txSucceeded $? "$RES"
# Desc: Must be called right after the `tx` command, makes sure it got a success response
txSucceeded() {
  if (assertTrue "sent tx: $2" $1); then
    TX=`echo $2 | cut -d: -f2-` # strip off first line asking for password
    assertEquals "good check: $TX" "0" $(echo $TX | jq .check_tx.code)
    assertEquals "good deliver: $TX" "0" $(echo $TX | jq .deliver_tx.code)
  else
    return 1
  fi
}

# XXX Ex Usage: checkSendTx $HASH $HEIGHT $SENDER $AMOUNT
# Desc: This looks up the tx by hash, and makes sure the height and type match
#       and that the first input was from this sender for this amount
checkSendTx() {
  TX=$(${CLIENT_EXE} query tx $1)
  assertTrue "found tx" $?
  assertEquals "proper height" $2 $(echo $TX | jq .height)
  assertEquals "type=send" '"send"' $(echo $TX | jq .data.type)
  assertEquals "proper sender" "\"$3\"" $(echo $TX | jq .data.data.inputs[0].address)
  assertEquals "proper out amount" "$4" $(echo $TX | jq .data.data.outputs[0].coins[0].amount)
  return $?
}

# XXX Ex Usage: waitForBlock $port
# Desc: Waits until the block height on that node increases by one
waitForBlock() {
  addr=http://localhost:$1
  b1=`curl -s $addr/status | jq .result.latest_block_height`
  b2=$b1
  while [ "$b2" == "$b1" ]; do
                echo "Waiting for node $addr to commit a block ..."
                sleep 1
    b2=`curl -s $addr/status | jq .result.latest_block_height`
  done
}


