package core

import "github.com/koinotice/vite/common/types"

func ConvertVoteToAddress(votes []*Vote) []types.Address {
	var result []types.Address
	for _, v := range votes {
		result = append(result, v.Addr)
	}
	return result
}
