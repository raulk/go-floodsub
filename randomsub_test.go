package floodsub

import (
	"bytes"
	"context"
	"strconv"
	"testing"
	"time"
)

func TestRandomSubBasic(t *testing.T) {
	var err error
	ctx := context.Background()
	hosts := getNetHosts(t, ctx, 50)

	rs := make([]*PubSub, len(hosts))
	for i, h := range hosts {
		if rs[i], err = NewRandomSub(ctx, h); err != nil {
			t.Fatal(err)
		}
	}

	subs := make([]*Subscription, len(hosts))
	for i, rs := range rs {
		if subs[i], err = rs.Subscribe("glasgow2018"); err != nil {
			t.Fatal(err)
		}
	}

	denseConnect(t, hosts)

	time.Sleep(2 * time.Second)

	for i, rs := range rs {
		data := []byte(strconv.Itoa(i))
		if err := rs.Publish("glasgow2018", data); err != nil {
			t.Errorf("error while publishing message: %v", err)
		}

	INNER:
		for j, sub := range subs {
			ctx2, cancelFn := context.WithTimeout(ctx, 2*time.Second)
			m, err := sub.Next(ctx2)
			cancelFn()
			if err != nil {
				t.Errorf("failed to receive next message on subscription; idx msg sent: %d, "+
					"idx sub failed: %d, err: %v", i, j, err)
				break INNER
			}
			if !bytes.Equal(m.Data, data) {
				t.Errorf("expected values to be equal, publisher: %d, sub: %d, "+
					"expected: %s, rcvd: %s", i, j, string(data), string(m.Data))
			}
		}
	}
}
