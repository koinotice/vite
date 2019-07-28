package api

import (
	"github.com/pkg/errors"
	"github.com/koinotice/vite/chain"
	"github.com/koinotice/vite/common/types"
	"github.com/koinotice/vite/log15"
	apidex "github.com/koinotice/vite/rpcapi/api/dex"
	"github.com/koinotice/vite/vite"
	"github.com/koinotice/vite/vm/contracts/dex"
	"github.com/koinotice/vite/vm/util"
	"math/big"
)

type DexFundApi struct {
	vite  *vite.Vite
	chain chain.Chain
	log   log15.Logger
}

func NewDexFundApi(vite *vite.Vite) *DexFundApi {
	return &DexFundApi{
		vite:  vite,
		chain: vite.Chain(),
		log:   log15.New("module", "rpc_api/dexfund_api"),
	}
}

func (f DexFundApi) String() string {
	return "DexFundApi"
}

type AccountFundInfo struct {
	TokenInfo *RpcTokenInfo `json:"tokenInfo,omitempty"`
	Available string        `json:"available"`
	Locked    string        `json:"locked"`
}

func (f DexFundApi) GetAccountFundInfo(addr types.Address, tokenId *types.TokenTypeId) (map[types.TokenTypeId]*AccountFundInfo, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	dexFund, _ := dex.GetUserFund(db, addr)
	fundInfo, err := dex.GetAccountFundInfo(dexFund, tokenId)
	if err != nil {
		return nil, err
	}

	accFundInfo := make(map[types.TokenTypeId]*AccountFundInfo, 0)
	for _, v := range fundInfo {
		tokenInfo, err := f.chain.GetTokenInfoById(v.Token)
		if err != nil {
			return nil, err
		}
		info := &AccountFundInfo{TokenInfo: RawTokenInfoToRpc(tokenInfo, v.Token)}
		a := "0"
		if v.Available != nil {
			a = v.Available.String()
		}
		info.Available = a

		l := "0"
		if v.Locked != nil {
			l = v.Locked.String()
		}
		info.Locked = l

		accFundInfo[v.Token] = info
	}
	return accFundInfo, nil
}

func (f DexFundApi) GetAccountFundInfoByStatus(addr types.Address, tokenId *types.TokenTypeId, status byte) (map[types.TokenTypeId]string, error) {
	if status != 0 && status != 1 && status != 2 {
		return nil, errors.New("args's status error, 1 for available, 2 for locked, 0 for total")
	}

	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	dexFund, _ := dex.GetUserFund(db, addr)
	fundInfo, err := dex.GetAccountFundInfo(dexFund, tokenId)
	if err != nil {
		return nil, err
	}

	fundInfoMap := make(map[types.TokenTypeId]string, 0)
	for _, v := range fundInfo {
		amount := big.NewInt(0)
		if a := v.Available; a != nil {
			if status == 0 || status == 1 {
				amount.Add(amount, a)
			}
		}
		if l := v.Locked; l != nil {
			if status == 0 || status == 2 {
				amount.Add(amount, l)
			}
		}
		fundInfoMap[v.Token] = amount.String()
	}
	return fundInfoMap, nil
}

func (f DexFundApi) GetOwner() (*types.Address, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	return dex.GetOwner(db)
}

func (f DexFundApi) GetTime() (int64, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return -1, err
	}
	return dex.GetTimestampInt64(db), nil
}

func (f DexFundApi) GetPeriodId() (uint64, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return 0, err
	}
	return dex.GetCurrentPeriodId(db, getConsensusReader(f.vite)), nil
}

func (f DexFundApi) GetTokenInfo(token types.TokenTypeId) (*apidex.RpcDexTokenInfo, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if tokenInfo, ok := dex.GetTokenInfo(db, token); ok {
		return apidex.TokenInfoToRpc(tokenInfo, token), nil
	} else {
		return nil, dex.InvalidTokenErr
	}
}

func (f DexFundApi) GetMarketInfo(tradeToken, quoteToken types.TokenTypeId) (*apidex.RpcMarketInfo, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if marketInfo, ok := dex.GetMarketInfo(db, tradeToken, quoteToken); ok {
		return apidex.MarketInfoToRpc(marketInfo), nil
	} else {
		return nil, dex.TradeMarketNotExistsErr
	}
}

func (f DexFundApi) GetCurrentDividendPools() (map[types.TokenTypeId]*apidex.DividendPoolInfo, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	var pools map[types.TokenTypeId]*apidex.DividendPoolInfo
	if feeSumByPeriod, ok := dex.GetCurrentFeeSum(db, getConsensusReader(f.vite)); !ok {
		return nil, nil
	} else {
		pools = make(map[types.TokenTypeId]*apidex.DividendPoolInfo)
		for _, pool := range feeSumByPeriod.FeesForDividend {
			tk, _ := types.BytesToTokenTypeId(pool.Token)
			if tokenInfo, ok := dex.GetTokenInfo(db, tk); !ok {
				return nil, dex.InvalidTokenErr
			} else {
				amt := new(big.Int).SetBytes(pool.DividendPoolAmount)
				pool := &apidex.DividendPoolInfo{amt.String(), tokenInfo.QuoteTokenType, apidex.TokenInfoToRpc(tokenInfo, tk)}
				pools[tk] = pool
			}
		}
	}
	return pools, nil
}

func (f DexFundApi) GetCurrentFeeSum() (*apidex.RpcFeeSumByPeriod, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if feeSum, ok := dex.GetCurrentFeeSum(db, getConsensusReader(f.vite)); !ok {
		return nil, nil
	} else {
		return apidex.FeeSumByPeriodToRpc(feeSum), nil
	}
}

func (f DexFundApi) GetFeeSumByPeriod(periodId uint64) (*apidex.RpcFeeSumByPeriod, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if feeSum, ok := dex.GetFeeSumByPeriodId(db, periodId); !ok {
		return nil, nil
	} else {
		return apidex.FeeSumByPeriodToRpc(feeSum), nil
	}
}

func (f DexFundApi) GetCurrentBrokerFeeSum(broker types.Address) (*apidex.RpcBrokerFeeSumByPeriod, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if brokerFeeSum, ok := dex.GetCurrentBrokerFeeSum(db, getConsensusReader(f.vite), broker.Bytes()); !ok {
		return nil, nil
	} else {
		return apidex.BrokerFeeSumByPeriodToRpc(brokerFeeSum), nil
	}
}

func (f DexFundApi) GetBrokerFeeSumByPeriod(periodId uint64, broker types.Address) (*apidex.RpcBrokerFeeSumByPeriod, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if brokerFeeSum, ok := dex.GetBrokerFeeSumByPeriodId(db, broker.Bytes(), periodId); !ok {
		return nil, nil
	} else {
		return apidex.BrokerFeeSumByPeriodToRpc(brokerFeeSum), nil
	}
}

func (f DexFundApi) GetUserFees(address types.Address) (*apidex.RpcUserFees, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if userFees, ok := dex.GetUserFees(db, address.Bytes()); !ok {
		return nil, nil
	} else {
		return apidex.UserFeesToRpc(userFees), nil
	}
}

func (f DexFundApi) GetVxSumFunds() (*apidex.RpcVxFunds, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if vxSumFunds, ok := dex.GetVxSumFunds(db); !ok {
		return nil, nil
	} else {
		return apidex.VxFundsToRpc(vxSumFunds), nil
	}
}

func (f DexFundApi) GetVxFunds(address types.Address) (*apidex.RpcVxFunds, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if vxFunds, ok := dex.GetVxFunds(db, address.Bytes()); !ok {
		return nil, nil
	} else {
		return apidex.VxFundsToRpc(vxFunds), nil
	}
}

func (f DexFundApi) GetVxMinePool() (string, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return "", err
	}
	balance := dex.GetVxMinePool(db)
	return balance.String(), nil
}

func (f DexFundApi) GetMakerProxyAmount(periodId uint64) (string, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return "", err
	}
	balance := dex.GetMakerProxyAmountByPeriodId(db, periodId)
	return balance.String(), nil
}

func (f DexFundApi) IsPledgeVip(address types.Address) (bool, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return false, err
	}
	_, ok := dex.GetPledgeForVip(db, address)
	return ok, nil
}

func (f DexFundApi) GetPledgeForVip(address types.Address) (*dex.PledgeVip, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if info, ok := dex.GetPledgeForVip(db, address); ok {
		return info, nil
	} else {
		return nil, nil
	}
}

func (f DexFundApi) GetPledgeForVX(address types.Address) (string, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return "", err
	}
	return dex.GetPledgeForVx(db, address).String(), nil
}

func (f DexFundApi) GetPledgesForVx(address types.Address) (*apidex.RpcPledgesForVx, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if pledgesForVx, ok := dex.GetPledgesForVx(db, address); ok {
		return apidex.PledgesForVxToRpc(pledgesForVx), nil
	} else {
		return nil, nil
	}
}

func (f DexFundApi) GetPledgesForVxSum() (*apidex.RpcPledgesForVx, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	if pledgesForVxSum, ok := dex.GetPledgesForVxSum(db); ok {
		return apidex.PledgesForVxToRpc(pledgesForVxSum), nil
	} else {
		return nil, nil
	}
}

func (f DexFundApi) GetFundConfig() (map[string]string, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	configs := make(map[string]string)
	owner, _ := dex.GetOwner(db)
	configs["owner"] = owner.String()
	if timer := dex.GetTimer(db); timer != nil {
		configs["timer"] = timer.String()
	}
	if trigger := dex.GetTrigger(db); trigger != nil {
		configs["trigger"] = trigger.String()
	}
	if mineProxy := dex.GetMakerMineProxy(db); mineProxy != nil {
		configs["mineProxy"] = mineProxy.String()
	}
	if maintainer := dex.GetMaintainer(db); maintainer != nil {
		configs["maintainer"] = maintainer.String()
	}
	return configs, nil
}

func  (f DexFundApi) IsViteXStopped() (bool, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return false, err
	}
	return dex.IsViteXStopped(db), nil
}

func (f DexFundApi) GetThresholdForTradeAndMine() (map[int]*apidex.RpcThresholdForTradeAndMine, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	thresholds := make(map[int]*apidex.RpcThresholdForTradeAndMine, 4)
	for tokenType := dex.ViteTokenType; tokenType <= dex.UsdTokenType; tokenType++ {
		tradeThreshold := dex.GetTradeThreshold(db, int32(tokenType))
		mineThreshold := dex.GetMineThreshold(db, int32(tokenType))
		thresholds[tokenType] = &apidex.RpcThresholdForTradeAndMine{TradeThreshold: tradeThreshold.String(), MineThreshold: mineThreshold.String()}
	}
	return thresholds, nil
}

func (f DexFundApi) GetInviterCode(address types.Address) (uint32, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return 0, err
	}
	return dex.GetCodeByInviter(db, address), nil
}

func (f DexFundApi) GetInviteeCode(address types.Address) (uint32, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return 0, err
	}
	if inviter, err := dex.GetInviterByInvitee(db, address); err != nil {
		return 0, err
	} else {
		return dex.GetCodeByInviter(db, *inviter), nil
	}
}

func (f DexFundApi) GetPeriodJobLastPeriodId(bizType uint8) (uint64, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return 0, err
	}
	if lastPeriodId := dex.GetLastJobPeriodIdByBizType(db, bizType); err != nil {
		return 0, err
	} else {
		return lastPeriodId, nil
	}
}

func (f DexFundApi) VerifyFundBalance() (*dex.FundVerifyRes, error) {
	db, err := getDb(f.chain, types.AddressDexFund)
	if err != nil {
		return nil, err
	}
	return dex.VerifyDexFundBalance(db, getConsensusReader(f.vite)), nil
}

func getConsensusReader(vite *vite.Vite) *util.VMConsensusReader {
	return util.NewVmConsensusReader(vite.Consensus().SBPReader())
}
