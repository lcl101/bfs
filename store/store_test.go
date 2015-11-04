package main

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	var (
		s      *Store
		z      *Zookeeper
		v      *Volume
		err    error
		buf    []byte
		data   = []byte("test")
		n      = &Needle{}
		file   = "./test/store.idx"
		bfile  = "./test/block_1"
		ifile  = "./test/block_1.idx"
		b2file = "./test/block_2"
		i2file = "./test/block_2.idx"
		b3file = "./test/block_3"
		i3file = "./test/block_3.idx"
	)
	os.Remove(file)
	os.Remove(bfile)
	os.Remove(ifile)
	os.Remove(b2file)
	os.Remove(i2file)
	os.Remove(b3file)
	os.Remove(i3file)
	defer os.Remove(file)
	defer os.Remove(bfile)
	defer os.Remove(ifile)
	defer os.Remove(b2file)
	defer os.Remove(i2file)
	defer os.Remove(b3file)
	defer os.Remove(i3file)
	t.Log("NewStore()")
	if z, err = NewZookeeper([]string{"localhost:2181"}, time.Second*1, "/rack/test/"); err != nil {
		t.Errorf("NewZookeeper() error(%v)", err)
		goto failed
	}
	z.DelVolume(1)
	if s, err = NewStore(z, file); err != nil {
		t.Errorf("NewStore() error(%v)", err)
		goto failed

	}
	defer s.Close()
	t.Log("AddFreeVolume")
	if _, err = s.AddFreeVolume(2, "./test", "./test"); err != nil {
		t.Errorf("s.AddFreeVolume() error(%v)", err)
		goto failed
	}
	t.Log("AddVolume(1)")
	if v, err = s.AddVolume(1); err != nil {
		t.Errorf("AddVolume() error(%v)", err)
		goto failed
	}
	time.Sleep(1 * time.Second)
	t.Log("Volumes[1]")
	if v = s.Volumes[1]; v == nil {
		err = fmt.Errorf("Volume(1) not exist")
		t.Error(err)
		goto failed
	}
	if err = n.Parse(1, 1, data); err != nil {
		t.Errorf("n.Parse() error(%v)", err)
		goto failed
	}
	if err = v.Add(n); err != nil {
		t.Errorf("v.Add(1) error(%v)", err)
		goto failed
	}
	buf = v.Buffer()
	defer v.FreeBuffer(buf)
	if _, err = v.Get(1, 1, buf); err != nil {
		t.Errorf("v.Get(1) error(%v)", err)
		goto failed
	}
	t.Log("BulkVolume()")
	if err = s.BulkVolume(1, b2file, i2file); err != nil {
		t.Errorf("Bulk(1) error(%v)", err)
		goto failed
	}
	time.Sleep(3 * time.Second)
	t.Log("Volumes[1]")
	if v = s.Volumes[1]; v == nil {
		err = fmt.Errorf("Volume(1) not exist")
		t.Error(err)
		goto failed
	}
	if err = v.Add(n); err != nil {
		t.Errorf("v.Add(1) error(%v)", err)
		goto failed
	}
	if _, err = v.Get(1, 1, buf); err != nil {
		t.Errorf("v.Get(1) error(%v)", err)
		goto failed
	}
	t.Log("CompactVolume()")
	if err = s.CompactVolume(1); err != nil {
		t.Errorf("Compress(1) error(%v)", err)
		goto failed
	}
	time.Sleep(2 * time.Second)
	if v = s.Volumes[1]; v == nil {
		err = fmt.Errorf("Volume(1) not exist")
		t.Error(err)
		goto failed
	}
	if _, err = v.Get(1, 1, buf); err != nil {
		t.Errorf("v.Get(1) error(%v)", err)
		goto failed
	}
	t.Log("DelVolume(1)")
	s.DelVolume(1)
	time.Sleep(1 * time.Second)
	if v = s.Volumes[1]; v != nil {
		err = fmt.Errorf("Volume(1) exist")
		t.Error(err)
		goto failed
	}
failed:
	if err != nil {
		t.FailNow()
	}
}
