package chain

import (
	"cess-gateway/configs"
	. "cess-gateway/internal/logger"
	"cess-gateway/tools"
	"math/big"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

// File meta information on chain
func FileMetaInfoOnChain(phrase, userwallet, filename, fileid, filehash string, public bool, backups uint8, filesize int64, downloadfee *big.Int) error {
	var (
		err         error
		accountInfo types.AccountInfo
	)
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic]: %v", err)
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
	api := getSubstrateAPI()
	defer func() {
		releaseSubstrateAPI()
		err := recover()
		if err != nil {
			Err.Sugar().Errorf("[panic]: %v", err)
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

	fileid_bytes, err := types.EncodeToBytes(fileid)
	if err != nil {
		return errors.Wrap(err, "EncodeToBytes")
	}

	c, err := types.NewCall(meta, ChainTx_FileBank_DeleteFile, types.NewBytes(fileid_bytes))
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
						if events.FileBank_DeleteFile[i].Acc == types.NewAccountID(keyring.PublicKey) && string(events.FileBank_DeleteFile[i].Fileid) == fileid {
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
