package blobkv

import (
	"context"

	"go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type KVServer struct {
	store *Store
	etcdserverpb.UnimplementedKVServer
}

type WatchServer struct {
	store *Store
	etcdserverpb.UnimplementedWatchServer
}

func NewKVServer(store *Store) *KVServer {
	return &KVServer{store: store}
}

func NewWatchServer(store *Store) *WatchServer {
	return &WatchServer{store: store}
}

func (s *KVServer) Range(ctx context.Context, req *etcdserverpb.RangeRequest) (*etcdserverpb.RangeResponse, error) {
	kv, err := s.store.Get(ctx, string(req.Key))
	if err != nil {
		return nil, err
	}

	return &etcdserverpb.RangeResponse{
		Header: &etcdserverpb.ResponseHeader{},
		Kvs: []*mvccpb.KeyValue{
			{
				Key:            []byte(kv.Key),
				Value:          kv.Value,
				CreateRevision: kv.CreateRev,
				ModRevision:   kv.ModRev,
				Version:       kv.Version,
				Lease:         kv.Lease,
			},
		},
	}, nil
}

func (s *KVServer) Put(ctx context.Context, req *etcdserverpb.PutRequest) (*etcdserverpb.PutResponse, error) {
	kv := &KeyValue{
		Key:       string(req.Key),
		Value:     req.Value,
		Version:   1,
		CreateRev: s.store.createRev,
		ModRev:    s.store.modRev,
	}

	err := s.store.Put(ctx, kv)
	if err != nil {
		return nil, err
	}

	return &etcdserverpb.PutResponse{
		Header: &etcdserverpb.ResponseHeader{},
	}, nil
}

func (s *KVServer) DeleteRange(ctx context.Context, req *etcdserverpb.DeleteRangeRequest) (*etcdserverpb.DeleteRangeResponse, error) {
	err := s.store.Delete(ctx, string(req.Key))
	if err != nil {
		return nil, err
	}

	return &etcdserverpb.DeleteRangeResponse{
		Header: &etcdserverpb.ResponseHeader{},
	}, nil
}

func (s *WatchServer) Watch(stream etcdserverpb.Watch_WatchServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		watchCreate := req.GetCreateRequest()
		ch, err := s.store.Watch(stream.Context(), string(watchCreate.Key))
		if err != nil {
			return err
		}

		go func() {
			for kv := range ch {
				resp := &etcdserverpb.WatchResponse{
					Header: &etcdserverpb.ResponseHeader{},
					Events: []*mvccpb.Event{
						{
							Type: mvccpb.PUT,
							Kv: &mvccpb.KeyValue{
								Key:            []byte(kv.Key),
								Value:          kv.Value,
								CreateRevision: kv.CreateRev,
								ModRevision:   kv.ModRev,
								Version:       kv.Version,
								Lease:         kv.Lease,
							},
						},
					},
				}
				if err := stream.Send(resp); err != nil {
					return
				}
			}
		}()
	}
}
