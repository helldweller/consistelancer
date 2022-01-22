package db

import (
    "encoding/json"
    "math/rand"
)

var (
    UpstreamList Upstreams
)

type Upstream struct {
    Host      string
    Port      string
    Scheme    string
}
type Upstreams struct {
    Items []Upstream
}

func (u *Upstreams) AddItem(item Upstream) []Upstream {
    u.Items = append(u.Items, item)
    return u.Items
}

func (u *Upstreams) MarshalBinary() ([]byte, error) {
    return json.Marshal(u)
}

func (u *Upstreams) UnmarshalBinary(data []byte) error {
    if err := json.Unmarshal(data, &u); err != nil {
        return err
    }
    return nil
}

func (u *Upstreams) GetRandomItem() Upstream {
	i := rand.Intn(len(u.Items))
	return u.Items[i]
}
