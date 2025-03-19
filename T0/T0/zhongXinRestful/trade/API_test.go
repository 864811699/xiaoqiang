package trade

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

)

func TestLogin(t *testing.T) {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(`-----BEGIN CERTIFICATE-----
MIIDVDCCAjygAwIBAgIEcImWfDANBgkqhkiG9w0BAQsFADBPMQswCQYDVQQGEwJi
ajELMAkGA1UECBMCYmoxCzAJBgNVBAcTAmJqMQswCQYDVQQKEwJ0bDELMAkGA1UE
CxMCdGwxDDAKBgNVBAMTA3lsajAeFw0yMzA3MTkwMjAxMTBaFw0zMzA3MTYwMjAx
MTBaME8xCzAJBgNVBAYTAmJqMQswCQYDVQQIEwJiajELMAkGA1UEBxMCYmoxCzAJ
BgNVBAoTAnRsMQswCQYDVQQLEwJ0bDEMMAoGA1UEAxMDeWxqMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAjrf6rkZ7UKsxmRtweBQ+0umX3qQZafmsme+C
Ub9xRYrZVK8tHY0Lhn3NUW9VgYh7qmOOCwrLLTdpi5+a6j77xiY3FZmCJCir61Bk
re9/5UejHBSvg5fhl08p6o2sfmaqWVhUlSutx0Ov/llwO3Gf5+rA+lVS/KeLIt+/
n6fbfhVDgJX1lM5JObYQx7MB8F6aPzmAlN9Mlu31lVU3qDqNYjm1H/+udjDjxCaz
ZUaUuWZ3RLh6wxvOMwrXeOBoCP3AmpzCfRkDcscGVjGRRhT8oHksUiXycrXGgUIn
s3aGaxrcQ8rqzcSABgtI/z8o6QtM71R4WtU2cfWRvepZ8UcJCwIDAQABozgwNjAV
BgNVHREEDjAMhwQKCgrbhwRqJyaqMB0GA1UdDgQWBBSrtBzOooUQQnXsPR+5XyHP
R729LTANBgkqhkiG9w0BAQsFAAOCAQEAQ9CAMMKR2S2hFoqvkrV3ikB6DrtL7giR
tVl3qPwRgf2hVb5zH4o7mtRY0JrJJFrU8iVQHHVl/oeD7+e8WU1QWhzG9SiNoXF8
N1wjD5u3vRC/Ph/dJOOxL09adNQDE2+R027AfBGRST6GcyVntmHy/mOc/946Ak0P
FklxRMdV4iiawPjsZQx/JlrTSH94oflriS/n1UTmWfpUhVY5ZGuQfFaEGszv0Oqc
24WkYymTiS35Ll+klpPFxOLRomSUDTVhYCIWWV6tF6oe4LYgOPcO7sUqgUAWrjDc
o0njWA8Ci5/+TglmUQyu6+OGANpD4vLXjK3eOSz5jS0qcjsWNC+Dpw==
-----END CERTIFICATE-----`))
	if !ok {
		panic("failed to parse root certificate")
	}

	tr:=&http.Transport{TLSClientConfig:&tls.Config{RootCAs:roots}}
	loginReq := LoginReqParams{
		UserName: "user1",
		Password: "pwd1",
	}
	client:=&http.Client{Transport:tr}
	bts, _ := json.Marshal(loginReq)
	rsp, err := client.Post("https://106.39.38.170:33808/sys/login", "application/json", bytes.NewReader(bts))
	if err != nil {
		t.Fatal( "login request fail", err)
	}
	defer rsp.Body.Close()
	bts, _ = ioutil.ReadAll(rsp.Body)
	loginRsp := RespBase{}
	if err := json.Unmarshal(bts, &loginRsp); err != nil {
		t.Fatal( "login response unmarshal fail,msg:", err, bts)
	}

	t.Fatal("login fail ,msg:", loginRsp.Message)
}

func TestLoadKey(t *testing.T) {

	// 加载证书和私钥
	cert, err := tls.LoadX509KeyPair("../etc/client.ks", "keystore-password")
	if err != nil {
		fmt.Println("加载证书和私钥失败:", err)
		return
	}

	// 加载信任存储
	caCert, err := ioutil.ReadFile("../etc/client.ts")
	if err != nil {
		fmt.Println("加载信任存储失败:", err)
		return
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// 创建一个自定义的TLS配置
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}

	// 创建一个自定义的Transport
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// 创建一个自定义的Client
	client := &http.Client{
		Transport: transport,
	}

	// 创建一个GET请求
	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取响应失败:", err)
		return
	}

	// 打印响应内容
	fmt.Println(string(body))
}

func TestLoadKey2(t *testing.T) {

	cliCrt, err := tls.LoadX509KeyPair("../etc/client.ts", "../etc/client.ks")
	if err != nil {
		fmt.Println("Loadx509keypair err:", err)
		return
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cliCrt},
		},
	}
	client := &http.Client{Transport: tr}
	t.Log(client)
}