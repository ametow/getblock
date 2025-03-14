package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"

	"github.com/ametow/getblock/config"
)

const contentType = "application/json"

func NewService(c *config.Config) *Service {
	return &Service{
		addresses: make(map[string]*big.Int),
		config:    c,
	}
}

type Service struct {
	config    *config.Config
	addresses map[string]*big.Int
	mutex     sync.Mutex
}

func (s *Service) Run(ctx context.Context) (*MaxTransaction, error) {
	latestBlockNumber, err := s.getLatestBlockNumber(ctx)
	if err != nil {
		return nil, err
	}

	blockNumberInt, _ := strconv.ParseInt(latestBlockNumber[2:], 16, 64)

	blocksChan, errc := s.walkBlocks(ctx, blockNumberInt)

	var wg sync.WaitGroup
	for i := 0; i < s.config.GoroutineCount; i++ {
		wg.Add(1)
		go s.worker(ctx, i, blocksChan, &wg)
	}
	wg.Wait()

	if err := <-errc; err != nil {
		return nil, err
	}

	return s.getMaxTransaction(), nil
}

func (s *Service) walkBlocks(ctx context.Context, blockNumber int64) (<-chan *BlockResponse, <-chan error) {
	ch := make(chan *BlockResponse)
	errc := make(chan error, 1)
	go func() {
		defer func() {
			close(ch)
			close(errc)
		}()
		for i := 0; i < s.config.BlockCount; i++ {
			blockNumberHex := fmt.Sprintf("0x%x", blockNumber-int64(i))
			block, err := s.getBlockByNumber(ctx, blockNumberHex)
			if err != nil {
				errc <- err
				return
			}
			select {
			case ch <- block:
			case <-ctx.Done():
				errc <- ctx.Err()
				return
			}
		}
	}()
	return ch, errc
}

func (s *Service) worker(ctx context.Context, id int, dataChan <-chan *BlockResponse, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Print("context Done worker", ctx.Err())
			return
		case block, ok := <-dataChan:
			if !ok {
				log.Print("worker closed of channel close")
				return
			}
			log.Print("processing block: ", block.Result.Number, "on worker: ", id)
			for _, tx := range block.Result.Transactions {
				tx := tx
				value, ok := new(big.Int).SetString(tx.Value, 0)
				if !ok {
					log.Printf("Failed to decode transaction value %s: worder: %d", tx.Value, id)
					continue
				}
				s.updateAddresses(&tx, value)
			}
		}
	}
}

func (s *Service) updateAddresses(tx *Transaction, value *big.Int) {
	var receiverBalance, senderBalance *big.Int
	var exists bool
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if senderBalance, exists = s.addresses[tx.From]; !exists {
		senderBalance = new(big.Int)
		s.addresses[tx.From] = senderBalance
	}
	if receiverBalance, exists = s.addresses[tx.To]; !exists {
		receiverBalance = new(big.Int)
		s.addresses[tx.To] = receiverBalance
	}
	s.addresses[tx.From] = senderBalance.Sub(senderBalance, value)
	s.addresses[tx.To] = receiverBalance.Add(receiverBalance, value)
}

func (s *Service) getMaxTransaction() *MaxTransaction {
	res := new(MaxTransaction)
	res.Wei = new(big.Int)
	for key, val := range s.addresses {
		if val.CmpAbs(res.Wei) > 0 {
			res.Address = key
			res.Wei = val
		}
	}
	res.Eth = weiToEth(res.Wei)
	oneEth := new(big.Float).SetFloat64(s.config.EthPrice)
	dollars := new(big.Float).Mul(res.Eth, oneEth)
	res.Dollars = fmt.Sprintf("%.2f $", dollars)
	return res
}

func (s *Service) getLatestBlockNumber(ctx context.Context) (string, error) {
	payload := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":"getblock.io"}`
	body, err := s.makeRequest(ctx, payload)
	if err != nil {
		return "", err
	}

	var result struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Result, nil
}

func (s *Service) getBlockByNumber(ctx context.Context, blockNumber string) (*BlockResponse, error) {
	payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":"getblock.io"}`, blockNumber)
	body, err := s.makeRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	var blockResponse BlockResponse
	if err := json.Unmarshal(body, &blockResponse); err != nil {
		return nil, err
	}

	return &blockResponse, nil
}

func (s *Service) makeRequest(ctx context.Context, payload string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.BaseApiUrl+s.config.ApiKey, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	defer client.CloseIdleConnections()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func weiToEth(wei *big.Int) *big.Float {
	ether := new(big.Float).SetInt(wei)
	base := new(big.Float).SetInt(big.NewInt(1e18))
	ether = ether.Quo(ether, base)
	return ether
}
