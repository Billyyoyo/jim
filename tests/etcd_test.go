package tests

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	. "jim/common/utils"
	"testing"
	"time"
)

func createClient() (cli *clientv3.Client) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"172.18.1.2:2379", "172.18.1.3:2379", "172.18.1.4:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	return
}

func TestGet(t *testing.T) {
	cli := createClient()
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	resp, err := kv.Get(context.Background(), "key", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		println(err.Error())
		return
	}
	if resp.Count > 0 {
		for _, kv := range resp.Kvs {
			println(string(kv.Key), ":", string(kv.Value))
		}
	}
}

//********很重要：withPrexxx  获取上一次操作的值   比如put和delete， 变化前的值
func TestLease(t *testing.T) {
	cli := createClient()
	defer cli.Close()
	lease := clientv3.NewLease(cli)
	grant, err := lease.Grant(context.Background(), 5)
	if err != nil {
		println(err.Error())
		return
	}
	kv := clientv3.NewKV(cli)
	_, err = kv.Put(context.Background(), "key05", "1232123234", clientv3.WithLease(grant.ID))
	if err != nil {
		println(err.Error())
		return
	}
	/**自动续期*/ //其实还是定时调用keepaliveonce
	leaseChan, err := lease.KeepAlive(context.Background(), grant.ID)
	if err != nil {
		printl(err.Error())
		return
	}
	//var aliveResp *clientv3.LeaseKeepAliveResponse
	for {
		select {
		case aliveResp := <-leaseChan:
			curr := GetCurrentMS()
			printl(aliveResp.TTL, curr)
		}
	}

	/**手动续期*/
	//	tick := time.Tick(3 * time.Second)
	//label:
	//	for {
	//		select {
	//		case <-tick:
	//			resp, errr := lease.KeepAliveOnce(context.Background(), grant.ID)
	//			if errr != nil {
	//				printl(errr.Error())
	//				break label
	//			}
	//		}
	//	}
}

func TestTransaction(t *testing.T) {
	cli := createClient()
	defer cli.Close()
	kv := clientv3.NewKV(cli)
	tx := kv.Txn(context.Background())
	tx.If(clientv3.Compare(clientv3.Value("key01"), "<", "3")).
		Then(clientv3.OpPut("key01", "5")).
		Else(clientv3.OpPut("key01", "0")).Commit()
}

func TestWatch(t *testing.T) {
	cli := createClient()
	defer cli.Close()
	watcher := cli.Watch(context.Background(), "key0", clientv3.WithPrefix())
	for {
		select {
		case ev := <-watcher:
			printj(ev)
		}
	}
}

func TestTest(t *testing.T) {
	//defer handleFatal()
	//happenFatal()
	//select {}
}

func handleFatal() {
	if r := recover(); r != nil {
		printl("fatal err:", r)
	}
}

func happenFatal() {
	//panic("some one died")
	m := 2
	n := m - 2
	printl(23 / n)
}
