package chain

import (
	. "cess-gateway/internal/logger"
	"cess-gateway/tools"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// Get scheduler information on the cess chain
func GetSchedulerInfo() ([]SchedulerInfo, error) {
	var (
		err  error
		data []SchedulerInfo
	)
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileMap, FileMap_SchedulerInfo)
	}

	key, err := types.CreateStorageKey(meta, State_FileMap, FileMap_SchedulerInfo)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileMap, FileMap_SchedulerInfo)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return nil, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileMap, FileMap_SchedulerInfo)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileMap, FileMap_SchedulerInfo)
	}
	return data, nil
}

// Get file meta information on the cess chain
func GetFileMetaInfo(fileid int64) (FileMetaInfo, error) {
	var (
		err  error
		data FileMetaInfo
	)

	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileBank, FileMap_FileMetaInfo)
	}
	fileid_s := fmt.Sprintf("%d", fileid)
	id, err := types.EncodeToBytes(fileid_s)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:EncodeToBytes]", State_FileBank, FileMap_FileMetaInfo)
	}
	key, err := types.CreateStorageKey(meta, State_FileBank, FileMap_FileMetaInfo, id)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileBank, FileMap_FileMetaInfo)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileBank, FileMap_FileMetaInfo)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileBank, FileMap_FileMetaInfo)
	}
	return data, nil
}

// Get user information on the cess chain
func GetSpaceDetailsInfo(prk string) ([]UserSpaceListInfo, error) {
	var (
		err  error
		data []UserSpaceListInfo
	)

	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileBank, FileBank_UserSpaceList)
	}

	keyring, err := signature.KeyringPairFromSecret(prk, 0)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:KeyringPairFromSecret]", State_FileBank, FileBank_UserSpaceList)
	}
	b, err := types.EncodeToBytes(types.NewAccountID(keyring.PublicKey))
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:KeyringPairFromSecret]", State_FileBank, FileBank_UserSpaceList)
	}

	key, err := types.CreateStorageKey(meta, State_FileBank, FileBank_UserSpaceList, b)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileBank, FileBank_UserSpaceList)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileBank, FileBank_UserSpaceList)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileBank, FileBank_UserSpaceList)
	}
	return data, nil
}

// Get user space information on the cess chain
func GetUserSpaceInfo(prk string) (UserStorageSpace, error) {
	var (
		err  error
		data UserStorageSpace
	)

	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileBank, FileBank_UserSpaceInfo)
	}

	keyring, err := signature.KeyringPairFromSecret(prk, 0)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:KeyringPairFromSecret]", State_FileBank, FileBank_UserSpaceList)
	}

	b, err := types.EncodeToBytes(types.NewAccountID(keyring.PublicKey))
	if err != nil {
		return data, err
	}
	key, err := types.CreateStorageKey(meta, State_FileBank, FileBank_UserSpaceInfo, b)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileBank, FileBank_UserSpaceInfo)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileBank, FileBank_UserSpaceInfo)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileBank, FileBank_UserSpaceInfo)
	}
	return data, nil
}

// Get file meta information on the cess chain
func GetFilelistInfo(wallet string) ([]types.Bytes, error) {
	var (
		err  error
		data []types.Bytes
	)

	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_FileBank, FileBank_UserFilelistInfo)
	}

	bytes, err := tools.DecodeToPub(wallet, tools.ChainCessTestPrefix)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:DecodeToPub]", State_FileBank, FileBank_UserFilelistInfo)
	}
	b, err := types.EncodeToBytes(types.NewAccountID(bytes))
	if err != nil {
		return data, err
	}
	key, err := types.CreateStorageKey(meta, State_FileBank, FileBank_UserFilelistInfo, b)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_FileBank, FileBank_UserFilelistInfo)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_FileBank, FileBank_UserFilelistInfo)
	}
	if !ok {
		return data, errors.Errorf("[%v.%v:GetStorageLatest value is nil]", State_FileBank, FileBank_UserFilelistInfo)
	}
	return data, nil
}

//Query sold space information on the cess chain
func QuerySoldSpace() (uint64, error) {
	var (
		err  error
		data types.U128
	)
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic] %v", err)
		}
	}()
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:GetMetadataLatest]", State_Sminer, Sminer_PurchasedSpace)
	}

	key, err := types.CreateStorageKey(meta, State_Sminer, Sminer_PurchasedSpace)
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:CreateStorageKey]", State_Sminer, Sminer_PurchasedSpace)
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return 0, errors.Wrapf(err, "[%v.%v:GetStorageLatest]", State_Sminer, Sminer_PurchasedSpace)
	}
	if !ok {
		return 0, nil
	}
	return data.Uint64(), nil
}
