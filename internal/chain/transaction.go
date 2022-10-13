package chain

import (
	"cess-gateway/configs"
	. "cess-gateway/internal/logger"
	"cess-gateway/tools"
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

func UploadDeclaration(transactionPrK, filehash, filename string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	api, err := GetRpcClient_Safe(configs.C.RpcAddr)
	defer Free()
	if err != nil {
		return txhash, errors.Wrap(err, "[GetRpcClient_Safe]")
	}

	meta, err := GetMetadata(api)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetMetadataLatest]")
	}

	c, err := types.NewCall(meta, ChainTx_FileBank_UploadDeclaration, types.NewBytes([]byte(filehash)), types.NewBytes([]byte(filename)))
	if err != nil {
		return txhash, errors.Wrap(err, "[NewCall]")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "[NewExtrinsic]")
	}

	genesisHash, err := GetGenesisHash(api)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetGenesisHash]")
	}

	rv, err := GetRuntimeVersion(api)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetRuntimeVersion]")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", configs.PublicKey, nil)
	if err != nil {
		return txhash, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil {
		return txhash, errors.Wrap(err, "[GetStorageLatest]")
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
		return txhash, errors.Wrap(err, "[Sign]")
	}

	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		return txhash, errors.Wrap(err, "[SubmitAndWatchExtrinsic]")
	}
	defer sub.Unsubscribe()
	timeout := time.After(configs.TimeToWaitEvents)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				events := MyEventRecords{}
				txhash, _ = types.EncodeToHex(status.AsInBlock)
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

				if len(events.FileBank_UploadDeclaration) > 0 {
					for i := 0; i < len(events.FileBank_UploadDeclaration); i++ {
						if string(events.FileBank_UploadDeclaration[i].File_hash) == filehash {
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

// Delete files in chain
func DeleteFileOnChain(phrase, fileid string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()

	var txhash string
	var accountInfo types.AccountInfo

	api, err := GetRpcClient_Safe(configs.C.RpcAddr)
	defer Free()
	if err != nil {
		return txhash, errors.Wrap(err, "GetRpcClient_Safe")
	}

	meta, err := GetMetadata(api)
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadataLatest")
	}

	c, err := types.NewCall(meta, ChainTx_FileBank_DeleteFile, types.NewBytes([]byte(fileid)))
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
				txhash, _ = types.EncodeToHex(status.AsInBlock)
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

				if len(events.FileBank_DeleteFile) > 0 {
					for i := 0; i < len(events.FileBank_DeleteFile); i++ {
						if string(events.FileBank_DeleteFile[i].Acc[:]) == string(configs.PublicKey) {
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
				txhash, _ = types.EncodeToHex(status.AsInBlock)
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
				txhash, _ = types.EncodeToHex(status.AsInBlock)
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
				txhash, _ = types.EncodeToHex(status.AsInBlock)
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

type FileHash struct {
	Hash [68]types.U8
}

func TestUpload(rpcaddr, privkey string) (string, error) {

	var txhash string
	var accountInfo types.AccountInfo

	api, err := NewRpcClient(rpcaddr)
	if err != nil {
		return txhash, errors.Wrap(err, "NewRpcClient")
	}

	keyring, err := signature.KeyringPairFromSecret(privkey, 0)
	if err != nil {
		return txhash, errors.Wrap(err, "KeyringPairFromSecret")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadata")
	}
	var hash [68]types.U8
	for i := 0; i < 68; i++ {
		hash[i] = types.U8(i)
	}
	fmt.Println(len(hash), " : ", hash)
	//b, _ := types.EncodeToBytes(hash)

	value := uint8(100)
	c, err := types.NewCall(meta, "FileBank.test_upload", hash, types.NewU8(value))
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return txhash, errors.Wrap(err, "GetGenesisHash")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetRuntimeVersion")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey, nil)
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

	// Sign the transaction
	err = ext.Sign(keyring, o)
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
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				return txhash, nil
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}

// Query file meta info
func TestQuery(rpcaddr, privkey string) (types.U8, error) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()
	var data types.U8

	api, err := NewRpcClient(rpcaddr)
	if err != nil {
		return data, errors.Wrap(err, "NewRpcClient")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrap(err, "GetMetadata")
	}
	var hash [68]types.U8
	for i := 0; i < 68; i++ {
		hash[i] = types.U8(i)
	}
	b, _ := types.Encode(hash)
	key, err := types.CreateStorageKey(meta, State_FileBank, "TestFile1", b)
	if err != nil {
		return data, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return data, errors.New(ERR_Empty)
	}
	return data, nil
}

const (
	file = iota
	filler
)

type enumtest struct {
	File types.U8
}

func TestUpload2(rpcaddr, privkey string) (string, error) {

	var txhash string
	var accountInfo types.AccountInfo

	api, err := NewRpcClient(rpcaddr)
	if err != nil {
		return txhash, errors.Wrap(err, "NewRpcClient")
	}

	keyring, err := signature.KeyringPairFromSecret(privkey, 0)
	if err != nil {
		return txhash, errors.Wrap(err, "KeyringPairFromSecret")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadata")
	}
	var hash [68]types.U8
	for i := 0; i < 68; i++ {
		hash[i] = types.U8(i)
	}
	fmt.Println(len(hash), " : ", hash)
	//b, _ := types.EncodeToBytes(hash)

	c, err := types.NewCall(meta, "FileBank.test_upload_level2", hash, types.U8(1))
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return txhash, errors.Wrap(err, "GetGenesisHash")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetRuntimeVersion")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey, nil)
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

	// Sign the transaction
	err = ext.Sign(keyring, o)
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
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				return txhash, nil
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}

// Query file meta info
func TestQuery2(rpcaddr, privkey string) (types.U8, error) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()
	var data types.U8

	api, err := NewRpcClient(rpcaddr)
	if err != nil {
		return data, errors.Wrap(err, "NewRpcClient")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrap(err, "GetMetadata")
	}
	var hash [68]types.U8
	for i := 0; i < 68; i++ {
		hash[i] = types.U8(i)
	}
	b, _ := types.Encode(hash)
	key, err := types.CreateStorageKey(meta, State_FileBank, "TestFile2", b)
	if err != nil {
		return data, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return data, errors.New(ERR_Empty)
	}
	return data, nil
}

type ssss struct {
	Hash [68]types.U8
	File types.U8
}

func TestUpload3(rpcaddr, privkey string) (string, error) {

	var txhash string
	var accountInfo types.AccountInfo

	api, err := NewRpcClient(rpcaddr)
	if err != nil {
		return txhash, errors.Wrap(err, "NewRpcClient")
	}

	keyring, err := signature.KeyringPairFromSecret(privkey, 0)
	if err != nil {
		return txhash, errors.Wrap(err, "KeyringPairFromSecret")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadata")
	}
	var ssssss ssss
	var hash [68]types.U8
	for i := 0; i < 68; i++ {
		hash[i] = types.U8(i)
	}
	fmt.Println(len(hash), " : ", hash)
	//b, _ := types.EncodeToBytes(hash)
	ssssss.Hash = hash
	ssssss.File = types.U8(1)
	c, err := types.NewCall(meta, "FileBank.test_upload_level3", hash, ssssss)
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return txhash, errors.Wrap(err, "GetGenesisHash")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetRuntimeVersion")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey, nil)
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

	// Sign the transaction
	err = ext.Sign(keyring, o)
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
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				return txhash, nil
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}

type ccccc1 struct {
	Index types.U8
	Value [4]types.U8
}
type ccccc2 struct {
	Index types.U8
	Value [8]types.U16
}
type cccc struct {
	IPv4 ccccc1
	IPv6 ccccc2
}

func TestRegisterScheduler(rpcaddr, privkey, stash string) (string, error) {

	var txhash string
	var accountInfo types.AccountInfo

	api, err := NewRpcClient(rpcaddr)
	if err != nil {
		return txhash, errors.Wrap(err, "NewRpcClient")
	}

	keyring, err := signature.KeyringPairFromSecret(privkey, 0)
	if err != nil {
		return txhash, errors.Wrap(err, "KeyringPairFromSecret")
	}

	keyring2, err := signature.KeyringPairFromSecret(stash, 0)
	if err != nil {
		return txhash, errors.Wrap(err, "KeyringPairFromSecret")
	}
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetMetadata")
	}

	var ccc cccc
	var hash [4]types.U8
	for i := 0; i < 4; i++ {
		hash[i] = types.U8(i + 1)
	}

	// ccc.IPv4.Index = types.U8(0)
	ccc.IPv4.Index = 0
	ccc.IPv4.Value = hash
	// b, err := types.EncodeToBytes(ccc.IPv4)
	c, err := types.NewCall(meta, "FileMap.registration_scheduler", types.NewAccountID(keyring2.PublicKey), ccc.IPv4)
	if err != nil {
		return txhash, errors.Wrap(err, "NewCall")
	}

	ext := types.NewExtrinsic(c)
	if err != nil {
		return txhash, errors.Wrap(err, "NewExtrinsic")
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return txhash, errors.Wrap(err, "GetGenesisHash")
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return txhash, errors.Wrap(err, "GetRuntimeVersion")
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyring.PublicKey, nil)
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

	// Sign the transaction
	err = ext.Sign(keyring, o)
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
				txhash, _ = types.EncodeToHex(status.AsInBlock)
				return txhash, nil
			}
		case err = <-sub.Err():
			return txhash, errors.Wrap(err, "<-sub")
		case <-timeout:
			return txhash, errors.New(ERR_Timeout)
		}
	}
}

// Query file meta info
func TestQueryFile(rpcaddr, privkey, fid string) (FileMetaInfo, error) {
	defer func() {
		if err := recover(); err != nil {
			Err.Sugar().Errorf("%v", tools.RecoverError(err))
		}
	}()
	var data FileMetaInfo

	api, err := NewRpcClient(rpcaddr)
	if err != nil {
		return data, errors.Wrap(err, "NewRpcClient")
	}

	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return data, errors.Wrap(err, "GetMetadata")
	}
	b, err := types.Encode(types.NewBytes([]byte(fid)))
	if err != nil {
		return data, errors.Wrap(err, "[EncodeToBytes]")
	}
	key, err := types.CreateStorageKey(meta, State_FileBank, "File", b)
	if err != nil {
		return data, errors.Wrap(err, "[CreateStorageKey]")
	}

	ok, err := api.RPC.State.GetStorageLatest(key, &data)
	if err != nil {
		return data, errors.Wrap(err, "[GetStorageLatest]")
	}
	if !ok {
		return data, errors.New(ERR_Empty)
	}
	return data, nil
}
