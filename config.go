package main

import (
	"crypto/sha1"
	"crypto/tls"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	stingray "github.com/whitepages/go-stingray"
)

// Config brocade login
type Config struct {
	ServerURL string
	Username  string
	Password  string
	VerifySSL bool
}

// Client initialise brocade connection
func (c *Config) Client() *stingray.Client {
	if c.VerifySSL {
		log.Printf("[INFO] Secure Brocade Client configured for server %s", c.ServerURL)
		return stingray.NewClient(nil, c.ServerURL, c.Username, c.Password)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	log.Printf("[INFO] Brocade Client configured for server %s", c.ServerURL)
	return stingray.NewClient(httpClient, c.ServerURL, c.Username, c.Password)
}

func expandStringList(configured []interface{}) []string {
	vs := make([]string, 0, len(configured))
	for _, v := range configured {
		vs = append(vs, v.(string))
	}
	return vs
}

// hashString returns a hash of the input for use as a StateFunc
func hashString(v interface{}) string {
	switch v.(type) {
	case string:
		hash := sha1.Sum([]byte(v.(string)))
		return hex.EncodeToString(hash[:])
	default:
		return ""
	}
}

func setBool(target **bool, d *schema.ResourceData, key string) {
	*target = stingray.Bool(d.Get(key).(bool))
}

func setInt(target **int, d *schema.ResourceData, key string) {
	*target = stingray.Int(d.Get(key).(int))
}

func setString(target **string, d *schema.ResourceData, key string) {
	*target = stingray.String(d.Get(key).(string))
}

func setStringList(target **[]string, d *schema.ResourceData, key string) {
	list := expandStringList(d.Get(key).([]interface{}))
	*target = &list
}

func setStringSet(target **[]string, d *schema.ResourceData, key string) {
	list := expandStringList(d.Get(key).(*schema.Set).List())
	*target = &list
}
