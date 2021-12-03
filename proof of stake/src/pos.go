package pos

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	DATE         = time.Now()
	GENESISBLOCK = map[string]interface{}{
		"Index":     0,
		"Timestamp": fmt.Sprintf("%v", DATE),
		"BPM":       0,
		"PrevHash":  "",
		"Validator": "", //address to receive the reward {validator, weight, age}
	}
)

var GENESISBLOCK2 = map[string]interface{}{
	"Index":     0,
	"Timestamp": fmt.Sprintf("%v", DATE),
	"BPM":       0,
	"PrevHash":  "",
	"Validator": "", //address to receive the reward {validator, weight, age}
}

var GENESISBLOCK3 = map[string]interface{}{
	"Index":     0,
	"Timestamp": fmt.Sprintf("%v", DATE),
	"BPM":       0,
	"PrevHash":  "",
	"Validator": "", //address to receive the reward {validator, weight, age}
}

var GENESISBLOCK4 = map[string]interface{}{
	"Index":     0,
	"Timestamp": fmt.Sprintf("%v", DATE),
	"BPM":       0,
	"PrevHash":  "",
	"Validator": "", //address to receive the reward {validator, weight, age}
}

type Blockchain struct {
	blockChain  []map[string]interface{}
	tempBlocks  []map[string]interface{}
	myCurrBlock map[string]interface{}
	validators  map[interface{}]*Blockchain //set
	nodes       map[interface{}]*Blockchain
	myAccount   map[string]interface{}
}

func NewBlockchain(_genesisBlock map[string]interface{}, account map[string]interface{}) (B *Blockchain) {
	B = new(Blockchain)
	//If the genesis block is valid, create chain
	B.blockChain = []map[string]interface{}{}
	B.tempBlocks = []map[string]interface{}{}
	B.myCurrBlock = map[string]interface{}{}
	B.validators = map[interface{}]*Blockchain{} //set
	B.nodes = map[interface{}]*Blockchain{}
	B.myAccount = map[string]interface{}{"Address": "", "Weight": 0, "Age": 0}
	B.myAccount["Address"] = account["Address"]
	B.myAccount["Weight"] = account["Weight"]
	func() {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					e := err
					fmt.Println("Invalid genesis block.\nOR\n" + fmt.Sprintf("%v", e))
					return
				}
				panic(r)
			}
		}()
		genesisBlock := B.generateGenesisBlock(_genesisBlock)
		if B.isBlockValid(genesisBlock.(map[string]interface{}), nil) {
			B.blockChain = append(B.blockChain, genesisBlock.(map[string]interface{}))
		} else {
			panic(fmt.Errorf("Exception: %v", "Unable to verify block"))
		}
	}()
	return
}

func (B *Blockchain) isBlockValid(block map[string]interface{}, prevBlock map[string]interface{}) bool {

	defer func() bool {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if strings.HasPrefix(err.Error(), "KeyError") {
					return false
				}
			}
			panic(r)
		}
		return false
	}()
	_hash := func(s *map[string]interface{}, h string) string {
		i := len(*s) - 1
		popped := (*s)[i]
		*s = append((*s)[:i], (*s)[i+1:]...)
		return popped
	}(&block, "Hash")
	defer func() bool {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				if strings.HasPrefix(err.Error(), "AssertionError") {
					return false
				}
			}
			panic(r)
		}
		return false
	}()
	hash2 := B.hasher(block)
	if !(_hash == string(hash2)) {
		panic(errors.New("AssertionError"))
	}
	prevHash := func() string {
		if prevBlock != nil {
			return prevBlock["Hash"].(string)
		}
		return ""
	}()
	block["Hash"] = _hash
	//hope this works
	if B.blockChain != nil {
		prevHash = func() string {
			if len(prevHash) == 0 {
				return B.blockChain[len(B.blockChain)-1]["Hash"].(string)
			}
			return prevHash
		}()
		func() {
			defer func() bool {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						if strings.HasPrefix(err.Error(), "AssertionError") {
							//e := err
							if prevHash == B.blockChain[0]["Hash"] {
								block["Hash"] = _hash
								return true
							}
							block["Hash"] = _hash
							return false

						}
					}
					panic(r)
				}
				return false
			}()
			if !reflect.DeepEqual(prevHash, block["PrevHash"]) {
				panic(errors.New("AssertionError"))
			}
		}()
	}
	block["Hash"] = _hash
	return true
}
func randint(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	min = 10
	max = 30
	return (rand.Intn(max-min+1) + min)
}

func (B *Blockchain) generateNewBlock(bpm int) map[string]interface{} {
	if bpm != 0 {
		bpm = bpm
	} else {
		bpm = randint(53, 63)
	}
	oldBlock := map[string]interface{}{}
	var address interface{}
	//hope this is ok
	if B.myCurrBlock != nil {
		return B.myCurrBlock
	}
	prevHash := B.blockChain[len(B.blockChain)-1]["Hash"].(string)
	index := func() int {
		//hope this is ok too
		if oldBlock != nil {
			return len(B.blockChain)
		}
		return int(oldBlock["Index"].(int)) + 1
	}()
	address = func() interface{} {
		if address != nil {
			return B.getValidator(B.myAccount)
		}
		return address
	}()
	newBlock := map[string]interface{}{
		"Index":     index,
		"Timestamp": fmt.Sprintf("%v", time.Now()),
		"BPM":       bpm,
		"PrevHash":  prevHash,
		"Validator": address,
	}
	newBlock["Hash"] = B.hasher(newBlock)
	if !B.isBlockValid(newBlock, nil) {
		panic(errors.New("AssertionError"))
	}
	B.myCurrBlock = newBlock
	return newBlock
}

func (B *Blockchain) getBlocksFromNodes() {
	if B.nodes != nil {
		for _, node := range B.nodes {
			node.addAnotherBlock(B.myCurrBlock)
			resp := node.generateNewBlock(0)
			if B.isBlockValid(resp, nil) {
				if !(func() int {
					for i, v := range B.validators {
						if v == resp["Validator"] {
							return i.(int)
						}
					}
					return -1
				}() != -1) {
					B.tempBlocks = append(B.tempBlocks, resp)
					B.validators[resp["Validator"]] = B
				}
			}
		}
	}
}

func (B *Blockchain) addAnotherBlock(anotherBlock map[string]interface{}) {
	if B.isBlockValid(anotherBlock, nil) {
		//if the validator of this block is not in the validator map[interface{}]struct{} or set
		if !(func() int {
			for i, v := range B.validators {
				if reflect.DeepEqual(v, anotherBlock["Validator"]) {
					return i.(int)
				}
			}
			return -1
		}() != -1) {
			B.tempBlocks = append(B.tempBlocks, anotherBlock)
			B.validators[anotherBlock["Validator"]] = B
		}
	}
}

func (B *Blockchain) pickWinner() []string {
	var acct []string
	//Creates a lottery pool of validators and choose the validator
	//who gets to forge the next block. Random selection weighted by amount of token staked
	//Do this every 30 seconds
	winner := []string{}
	B.tempBlocks = append(B.tempBlocks, B.myCurrBlock)
	B.validators[B.myCurrBlock["Validator"]] = B
	for _, validator := range B.validators {
		acct = validator.rsplit(", ")
		s1, _ := strconv.Atoi(acct[1])
		s2, _ := strconv.Atoi(acct[2])
		acct = append(acct, s1*s2)
		if winner && acct[len(acct)-1] {
			winner = func() []string {
				if winner[len(winner)-1] < acct[len(acct)-1] {
					return acct
				}
				return winner
			}()
		} else {
			winner = func() []string {
				//hope
				if acct[len(acct)-1] != "" {
					return acct
				}
				return winner
			}()
		}
	}
	if len(winner) != 0 {
		return winner
	}
	for _, validator := range B.validators {
		acct := validator.rsplit(", ")
		acct = append(acct, float64(int(acct[1])+int(acct[2]))/len(acct[0]))
		if len(winner) != 0 {
			winner = func() []string {
				if winner[len(winner)-1] < acct[len(acct)-1] {
					return acct
				}
				return winner
			}()
		} else {
			winner = acct
		}
	}
	return winner
}

func (B *Blockchain) pos() {
	var newBlock map[string]interface{}
	//get other's stakes,add owns claim,pick winner
	fmt.Println(
		fmt.Sprintf("%v", B.myAccount) + " =======================> Getting Valid chain\n",
	)
	B.resolveConflict()
	time.Sleep(1 * time.Second)
	B._pos()
	fmt.Println("***Calling other nodes to announce theirs***" + "\n")
	time.Sleep(1 * time.Second)
	for _, node := range B.nodes {
		node._pos()
	}
	time.Sleep(1 * time.Second)
	for _, block := range B.tempBlocks {
		validator := strings.Fields(block["Validator"].(string))
		if validator[0] == B.pickWinner()[0] {
			newBlock = block
			break
		} else {
		}
	}
	fmt.Println("New Block ====> " + fmt.Sprintf("%v", newBlock) + "\n")
	time.Sleep(1 * time.Second)
	B.addNewBlock(newBlock)
	for _, node := range B.nodes {
		node.addNewBlock(newBlock)
	}
	fmt.Println("Process ends" + "\n")
}

func (B *Blockchain) announceWinner() {
	B.blockChain = append(B.blockChain, B.myCurrBlock)
}

func (B *Blockchain) addNewBlock(block map[string]interface{}) {
	if B.isBlockValid(block, nil) {
		B.blockChain = append(B.blockChain, block)
		acct := strings.Fields(block["Validator"].(string))
		if B.myAccount["Address"] != acct[0] {
			if i, ok := B.myAccount["Age"].(int); ok {
				i += 1
			}
		} else {
			if x, ok := B.myAccount["Weight"].(int); ok {
				x += randint(1, 10) * B.myAccount["Age"].(int)
			}
			B.myAccount["Age"] = 0
		}
	}
	B.tempBlocks = []map[string]interface{}{}
	B.myCurrBlock = map[string]interface{}{}
	B.validators = map[interface{}]*Blockchain{}
}

func (B *Blockchain) _pos() {
	fmt.Println(
		"Coming from ==========================> " + fmt.Sprintf("%v", B.myAccount) + "\n",
	)
	time.Sleep(1 * time.Second)
	fmt.Println("***Generating new stake block***" + "\n")
	time.Sleep(1 * time.Second)
	B.generateNewBlock(0)
	fmt.Println("***Exchanging temporary blocks with other nodes***" + "\n")
	time.Sleep(1 * time.Second)
	B.getBlocksFromNodes()
	fmt.Println("***Picking a winner***" + "\n")
	time.Sleep(1 * time.Second)
	fmt.Println(
		"Winner is =======================> " + fmt.Sprintf("%v", B.pickWinner()) + "\n",
	)
}

func (B *Blockchain) resolveConflict() {
	for _, node := range B.nodes {
		if len(node.blockChain) > len(B.blockChain) {
			if B.isChainValid(node.blockChain) {
				fmt.Println("***Replacing node***" + "\n")
				B.blockChain = node.blockChain
				return
			}
		}
	}
	fmt.Println("***My chain is authoritative***" + "\n")
	return
}

func (B *Blockchain) isChainValid(chain []map[string]interface{}) bool {
	_prevBlock := map[string]interface{}{}
	for _, block := range chain {
		if B.isBlockValid(block, _prevBlock) {
			_prevBlock = block
		} else {
			return false
		}
	}
	return true
}

func (B *Blockchain) addNewNode(newNode interface{}) {
	B.nodes[newNode] = B
	newNode.(*Blockchain).addAnotherNode(B)
}

func (B *Blockchain) addAnotherNode(anotherNode interface{}) {
	B.nodes[anotherNode] = B
}

func (B *Blockchain) hasher(block map[string]interface{}) []byte {
	blockString, _ := json.MarshalIndent(block, "", " ")
	hash := sha256.Sum256(blockString)
	return hash[:]
}

func (B *Blockchain) getValidator(address map[string]interface{}) interface{} {
	return strings.Join(
		[]string{
			fmt.Sprintf("%v", address["Address"]),
			fmt.Sprintf("%v", address["Weight"]),
			fmt.Sprintf("%v", address["Age"]),
		},
		", ",
	)
}

func (B *Blockchain) generateGenesisBlock(genesisblock map[string]interface{}) interface{} {
	address := map[string]interface{}{"Address": "eltneg", "Weight": 50, "Age": 0}
	address = B.getValidator(address).(map[string]interface{})
	genesisblock["Index"] = func() int {
		if genesisblock["Index"].(string) != "" {
			return 0
		}
		return genesisblock["Index"].(int)
	}()
	genesisblock["Timestamp"] = func() string {
		if genesisblock["Timestamp"] == "" {
			return fmt.Sprintf("%v", time.Now())
		}
		return genesisblock["Timestamp"].(string)
	}()
	genesisblock["BPM"] = func() int {
		if genesisblock["BPM"] == 0 {
			return 0
		}
		return genesisblock["BPM"].(int)
	}()
	genesisblock["PrevHash"] = "0000000000000000"
	genesisblock["Validator"] = func() interface{} {
		if genesisblock["Validator"] == "" {
			return address
		}
		return genesisblock["Validator"]
	}()
	genesisblock["Hash"] = B.hasher(genesisblock)
	return genesisblock
}

func main() {
	//"Run test"
	account := map[string]interface{}{"Address": "eltneg", "Weight": 50}
	account2 := map[string]interface{}{"Address": "account2", "Weight": 55}
	account3 := map[string]interface{}{"Address": "account3", "Weight": 43}
	account4 := map[string]interface{}{"Address": "account4", "Weight": 16}
	blockchain := NewBlockchain(GENESISBLOCK, account)
	blockchain.generateNewBlock(52)

	blockchain2 := NewBlockchain(GENESISBLOCK2, account2)
	blockchain3 := NewBlockchain(GENESISBLOCK3, account3)

	clients := []*Blockchain{blockchain, blockchain2, blockchain3}
	blockchain.addNewNode(blockchain2)

	blockchain.addNewNode(blockchain3)
	blockchain2.addNewNode(blockchain)
	blockchain2.addNewNode(blockchain3)

	blockchain.getBlocksFromNodes()
	blockchain2.getBlocksFromNodes()
	blockchain.pickWinner()
	//check if temp blocks are same

	blockchain.pos()
	blockchain2.pos()
	blockchain3.pos()

	blockchain4 := NewBlockchain(GENESISBLOCK4, account4)
	blockchain4.addNewNode(blockchain)
	blockchain4.addNewNode(blockchain2)
	blockchain4.addNewNode(blockchain3)
	blockchain4.pos()
	clients = append(clients, blockchain4)
	for {
		fmt.Println("============================================ \n\n")
		client := clients[randint(0, 3)]
		client.pos()
	}
}
