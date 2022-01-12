package main

import (
    "context"
    "crypto/tls"
    "errors"
    "fmt"
    "io/ioutil"
    "time"
    "net/http"

    gjson "github.com/tidwall/gjson"
    log "github.com/sirupsen/logrus"
    "package/main/locker"
    "package/main/db"
)

func k8sApiWatcher(interval int, msgChannel chan db.Upstreams, groupCtx context.Context) error {
    log.Info("Starting k8sApiWatcher")
    ticker := time.NewTicker(time.Duration(interval) * time.Second)
    lockclient := locker.Initialize(groupCtx, config.RedisHost + ":" + config.RedisPort)
    for {
        select {
        case <-ticker.C:
            if !lockclient.IsMaster() {
                log.Info("I'am slave!")
                continue
            }
            result, err := getEndpoints(config.K8sService)
            if err != nil {
                log.Error(err)
                //return err
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

func getEndpoints(k8sServiceName string) (db.Upstreams, error) {
    result := db.Upstreams{}
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    client := &http.Client{Transport: tr}
    reqUri := "https://" + config.K8sApiHost + ":" + config.K8sApiPort + "/api/v1/namespaces/" + config.K8sNamespace + "/endpoints/" + k8sServiceName
    req, err := http.NewRequest("GET", string(reqUri), nil)
    if err != nil {
        return result, err
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("Authorization", "Bearer " + config.K8sToken)
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
