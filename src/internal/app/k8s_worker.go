package app

import (
    "context"
    "crypto/tls"
    "errors"
    "fmt"
    "io/ioutil"
    "time"
    "net/http"
    gjson "github.com/tidwall/gjson"

    log "package/main/internal/logger"
    "package/main/internal/config"
    "package/main/internal/locker"
    "package/main/internal/db"
)

func k8sApiWatcher(groupCtx context.Context, interval int, msgChannel chan db.Upstreams, conf *config.Config) error {
    defer close(msgChannel)
    log.Info("Starting k8sApiWatcher")
    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    lockclient := locker.Initialize(groupCtx, conf.RedisHost + ":" + conf.RedisPort)
    for {
        select {
        case <-ticker.C:
            if !lockclient.IsMaster() {
                log.Info("I'am slave!")
                continue
            }
            result, err := getEndpoints(conf)
            if err != nil {
                log.Error(err.Error())
            }
            msgChannel<- result
        case <-groupCtx.Done():
            log.Error("Closing k8sApiWatcher goroutine")
            lockclient.Close()
            return groupCtx.Err()
        }
    }
    return nil
}

func getEndpoints(conf *config.Config) (db.Upstreams, error) {
    result := db.Upstreams{}
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    reqUri := "https://" + conf.K8sApiHost + ":" + conf.K8sApiPort + "/api/v1/namespaces/" + conf.K8sNamespace + "/endpoints/" + conf.K8sService
    req, err := http.NewRequest("GET", string(reqUri), nil)
    if err != nil {
        return result, err
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", "Bearer " + conf.K8sToken)
    resp, err := client.Do(req)
    if err != nil {
        return result, err
    }
    if resp.StatusCode != http.StatusOK {
        return result, fmt.Errorf("K8s api response status: %s", resp.Status)
    }
    defer resp.Body.Close()
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return result, err
    }
    result, err = parseUpstreamList(body)
    if err != nil {
        return result, err
    }
    return result, nil
}

func parseUpstreamList(json []byte) (db.Upstreams, error) {
    result := db.Upstreams{}
    var err error = nil
    scheme := ""
    ips := gjson.GetBytes(json, "subsets.#.addresses.#.ip|@flatten")
    ports := gjson.GetBytes(json, "subsets.#.ports.#.port|@flatten")
    if (len(ips.Array()) < 1) || (len(ports.Array()) < 1) {
        err = errors.New("Ip or port list is empty")
    }
    for _, ip := range ips.Array() {
        for _, port := range ports.Array() {
            if (port.String() == "443") {
                scheme = "https"
            } else {
                scheme = "http"
            }
            upstream := db.Upstream{Host: ip.String(), Port: port.String(), Scheme: scheme}
            result.AddItem(upstream)
        }
    }
    return result, err
}
