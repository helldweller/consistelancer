package main
// перенести в db

import (
    "encoding/json"
)

var (
    upstreams Upstreams
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
