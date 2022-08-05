package chain

import (
	"cess-gateway/configs"
	. "cess-gateway/internal/logger"
	"cess-gateway/tools"
	"fmt"
	"math/big"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

//
func UploadDeclaration(transactionPrK, filehash, filename string) (string, int, error) {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api, err := NewRpcClient(configs.C.RpcAddr)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "NewRpcClient")
	}
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()
	keyring, err := signature.KeyringPairFromSecret(transactionPrK, 0)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[KeyringPairFromSecret]")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[GetMetadataLatest]")
	}

	c, err := types.NewCall(meta, ChainTx_FileBank_UploadDeclaration, types.NewBytes([]byte(filehash)), types.NewBytes([]byte(filename)))
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[NewExtrinsic]")
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[GetBlockHash]")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[GetRuntimeVersionLatest]")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[CreateStorageKey System  Account]")
	}

	keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[CreateStorageKey System Events]")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return "", configs.Code_500, errors.New("[GetStorageLatest value is empty]")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return "", configs.Code_500, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(configs.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				txhash := fmt.Sprintf("%#x", status.AsInBlock)
				h, err := api.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return txhash, configs.Code_600, err
				}

				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					Out.Sugar().Infof("[%v]Decode event err:%v", txhash, err)
				}

				for i := 0; i < len(events.FileBank_UploadDeclaration); i++ {
					if string(events.FileBank_UploadDeclaration[i].FileHash) == filehash {
						return txhash, configs.Code_200, nil
					}
				}
				return txhash, configs.Code_600, errors.Errorf("events.FileBank_FillerUpload not found")
			}
		case err = <-sub.Err():
			return "", configs.Code_500, err
		case <-timeout:
			return "", configs.Code_500, errors.New("Timeout")
		}
	}
}

// File meta information on chain
func FileMetaInfoOnChain(phrase, userwallet, filename, fileid, filehash string, public bool, backups uint8, filesize int64, downloadfee *big.Int) error {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api, err := NewRpcClient(configs.C.RpcAddr)
	if err != nil {
		return errors.Wrap(err, "NewRpcClient")
	}
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()
	keyring, err := signature.KeyringPairFromSecret(phrase, 0)
	if err != nil {
		return errors.Wrap(err, "KeyringPairFromSecret")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return errors.Wrap(err, "GetMetadataLatest")
	}

	c, err := types.NewCall(
		meta,
		ChainTx_FileBank_Upload,
		types.Bytes([]byte(userwallet)),
		types.Bytes([]byte(filename)),
		types.Bytes([]byte(fileid)),
		types.Bytes([]byte(filehash)),
		types.Bool(public),
		types.U8(backups),
		types.U64(filesize),
		types.NewU128(*downloadfee),
	)
	if err != nil {
		return errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return errors.Wrap(err, "GetBlockHash")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return errors.Wrap(err, "GetRuntimeVersionLatest")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey)
	if err != nil {
		return errors.Wrap(err, "CreateStorageKey")
	}

	keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return errors.Wrap(err, "CreateStorageKey Events")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return errors.Wrap(err, "GetStorageLatest")
	}
	if !ok {
		return errors.New("GetStorageLatest return value is empty")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return errors.Wrap(err, "Sign")
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return errors.Wrap(err, "SubmitAndWatchExtrinsic")
	}
	defer sub.Unsubscribe()
	var head *types.Header
	t := tools.RandomInRange(10000000, 99999999)
	timeout := time.After(configs.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				head, err = api.RPC.Chain.GetHeader(status.AsInBlock)
				if err == nil {
					Out.Sugar().Infof("[%v] [%v]", t, head.Number)
				}
				h, err := api.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return errors.Wrapf(err, "[%v]", t)
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					Out.Sugar().Infof("[%v]Decode event err:%v", t, err)
				}
				if events.FileBank_FileUpload != nil {
					for i := 0; i < len(events.FileBank_FileUpload); i++ {
						if events.FileBank_FileUpload[i].Acc == types.NewAccountID(keyring.PublicKey) {
							return nil
						}
					}
					return errors.Errorf("[%v]events.FileBank_FileUpload data err", t)
				}
				return errors.Errorf("[%v]events.FileBank_FileUpload not found", t)
			}
		case err = <-sub.Err():
			return errors.Wrapf(err, "[%v]", t)
		case <-timeout:
			return errors.Errorf("[%v]upload file meta info timeout,please check your Internet!", t)
		}
	}
}

// Delete files in chain
func DeleteFileOnChain(phrase, fileid string) error {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api, err := NewRpcClient(configs.C.RpcAddr)
	if err != nil {
		return errors.Wrap(err, "NewRpcClient")
	}
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()
	keyring, err := signature.KeyringPairFromSecret(phrase, 0)
	if err != nil {
		return errors.Wrap(err, "KeyringPairFromSecret")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return errors.Wrap(err, "GetMetadataLatest")
	}

	c, err := types.NewCall(meta, ChainTx_FileBank_DeleteFile, types.NewBytes([]byte(fileid)))
	if err != nil {
		return errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return errors.Wrap(err, "GetBlockHash")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return errors.Wrap(err, "GetRuntimeVersionLatest")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey)
	if err != nil {
		return errors.Wrap(err, "CreateStorageKey")
	}

	keye, err := types.CreateStorageKey(meta, "System", "Events", nil)
	if err != nil {
		return errors.Wrap(err, "CreateStorageKey Events")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return errors.Wrap(err, "GetStorageLatest")
	}
	if !ok {
		return errors.New("GetStorageLatest return value is empty")
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	// Sign the transaction
	err = ext.Sign(keyring, o)
	if err != nil {
		return errors.Wrap(err, "Sign")
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return errors.Wrap(err, "SubmitAndWatchExtrinsic")
	}
	defer sub.Unsubscribe()

	var head *types.Header
	t := tools.RandomInRange(10000000, 99999999)
	timeout := time.After(configs.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				head, err = api.RPC.Chain.GetHeader(status.AsInBlock)
				if err == nil {
					Out.Sugar().Infof("[%v] [%v]", t, head.Number)
				}
				h, err := api.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return errors.Wrapf(err, "[%v]", t)
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					Out.Sugar().Infof("[%v]Decode event err:%v", t, err)
				}
				if events.FileBank_DeleteFile != nil {
					for i := 0; i < len(events.FileBank_DeleteFile); i++ {
						if events.FileBank_DeleteFile[i].Acc == types.NewAccountID(keyring.PublicKey) && string(events.FileBank_DeleteFile[i].Fileid) == string(fileid) {
							return nil
						}
					}
					return errors.Errorf("[%v]events.FileBank_DeleteFile data err", t)
				}
				return errors.Errorf("[%v]events.FileBank_DeleteFile not found", t)
			}
		case err = <-sub.Err():
			return errors.Wrapf(err, "[%v]", t)
		case <-timeout:
			return errors.Errorf("[%v]delete file timeout,please check your Internet!", t)
		}
	}
}

//
func GetAddressFromPrk(prk string, prefix []byte) (string, error) {
	keyring, err := signature.KeyringPairFromSecret(prk, 0)
	if err != nil {
		return "", errors.Wrap(err, "[KeyringPairFromSecret]")
	}
	addr, err := tools.Encode(keyring.PublicKey, tools.ChainCessTestPrefix)
	if err != nil {
		return "", errors.Wrap(err, "[Encode]")
	}
	return addr, nil
}

//
func GetPubkeyFromPrk(prk string) ([]byte, error) {
	keyring, err := signature.KeyringPairFromSecret(prk, 0)
	if err != nil {
		return nil, errors.Wrap(err, "[KeyringPairFromSecret]")
	}
	return keyring.PublicKey, nil
}

func BuySpacePackage(package_type types.U8, count types.U128) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	api, err := NewRpcClient(configs.C.RpcAddr)
	if err != nil {
		return txhash, errors.Wrap(err, "NewRpcClient")
	}

	meta, err := GetMetadata(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadata")
	}

	c, err := types.NewCall(meta, ChainTx_FileBank_BuyPackage, package_type, count)
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := GetGenesisHash(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetGenesisHash")
	}

	rv, err := GetRuntimeVersion(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetRuntimeVersion")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", configs.PublicKey, nil)
	if err != nil {
		return txhash, errors.Wrap(err, "CreateStorageKey")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "GetStorageLatest")
	}

	if !ok {
		return txhash, errors.New(ERR_Empty)
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	kring, err := GetKeyring()
	if err != nil {
		return txhash, errors.Wrap(err, "GetKeyring")
	}

	// Sign the transaction
	err = ext.Sign(kring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "Sign")
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return txhash, errors.Wrap(err, "SubmitAndWatchExtrinsic")
	}

	defer sub.Unsubscribe()
	timeout := time.After(configs.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				txhash, _ = types.EncodeToHexString(status.AsInBlock)
				keye, err := GetKeyEvents()
				if err != nil {
					return txhash, errors.Wrap(err, "GetKeyEvents")
				}
				h, err := api.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "GetStorageRaw")
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					Out.Sugar().Infof("[%v]Decode event err:%v", txhash, err)
				}

				if len(events.FileBank_BuyPackage) > 0 {
					for i := 0; i < len(events.FileBank_DeleteFile); i++ {
						if events.FileBank_BuyPackage[i].Acc == types.NewAccountID(configs.PublicKey) {
							return txhash, nil
						}
					}
				}

				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}

func UpgradeSpacePackage(package_type types.U8, count types.U128) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	api, err := NewRpcClient(configs.C.RpcAddr)
	if err != nil {
		return txhash, errors.Wrap(err, "NewRpcClient")
	}

	meta, err := GetMetadata(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadata")
	}

	c, err := types.NewCall(meta, ChainTx_FileBank_UpgradePackage, package_type, count)
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := GetGenesisHash(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetGenesisHash")
	}

	rv, err := GetRuntimeVersion(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetRuntimeVersion")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", configs.PublicKey, nil)
	if err != nil {
		return txhash, errors.Wrap(err, "CreateStorageKey")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "GetStorageLatest")
	}

	if !ok {
		return txhash, errors.New(ERR_Empty)
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	kring, err := GetKeyring()
	if err != nil {
		return txhash, errors.Wrap(err, "GetKeyring")
	}

	// Sign the transaction
	err = ext.Sign(kring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "Sign")
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return txhash, errors.Wrap(err, "SubmitAndWatchExtrinsic")
	}

	defer sub.Unsubscribe()
	timeout := time.After(configs.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				txhash, _ = types.EncodeToHexString(status.AsInBlock)
				keye, err := GetKeyEvents()
				if err != nil {
					return txhash, errors.Wrap(err, "GetKeyEvents")
				}
				h, err := api.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "GetStorageRaw")
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					Out.Sugar().Infof("[%v]Decode event err:%v", txhash, err)
				}

				if len(events.FileBank_PackageUpgrade) > 0 {
					for i := 0; i < len(events.FileBank_PackageUpgrade); i++ {
						if events.FileBank_PackageUpgrade[i].Acc == types.NewAccountID(configs.PublicKey) {
							return txhash, nil
						}
					}
				}

				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}

func Renewal() (string, error) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	api, err := NewRpcClient(configs.C.RpcAddr)
	if err != nil {
		return txhash, errors.Wrap(err, "NewRpcClient")
	}

	meta, err := GetMetadata(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadata")
	}

	c, err := types.NewCall(meta, ChainTx_FileBank_RenewalPackage)
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := GetGenesisHash(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetGenesisHash")
	}

	rv, err := GetRuntimeVersion(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetRuntimeVersion")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", configs.PublicKey, nil)
	if err != nil {
		return txhash, errors.Wrap(err, "CreateStorageKey")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "GetStorageLatest")
	}

	if !ok {
		return txhash, errors.New(ERR_Empty)
	}

	o := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(accountInfo.Nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	kring, err := GetKeyring()
	if err != nil {
		return txhash, errors.Wrap(err, "GetKeyring")
	}

	// Sign the transaction
	err = ext.Sign(kring, o)
	if err != nil {
		return txhash, errors.Wrap(err, "Sign")
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return txhash, errors.Wrap(err, "SubmitAndWatchExtrinsic")
	}

	defer sub.Unsubscribe()
	timeout := time.After(configs.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				txhash, _ = types.EncodeToHexString(status.AsInBlock)
				keye, err := GetKeyEvents()
				if err != nil {
					return txhash, errors.Wrap(err, "GetKeyEvents")
				}
				h, err := api.RPC.State.GetStorageRaw(keye, status.AsInBlock)
				if err != nil {
					return txhash, errors.Wrap(err, "GetStorageRaw")
				}
				err = types.EventRecordsRaw(*h).DecodeEventRecords(meta, &events)
				if err != nil {
					Out.Sugar().Infof("[%v]Decode event err:%v", txhash, err)
				}

				if len(events.FileBank_PackageRenewal) > 0 {
					for i := 0; i < len(events.FileBank_PackageRenewal); i++ {
						if events.FileBank_PackageRenewal[i].Acc == types.NewAccountID(configs.PublicKey) {
							return txhash, nil
						}
					}
				}

				return txhash, errors.New(ERR_Failed)
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}
