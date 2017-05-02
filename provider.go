package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADE_SERVER_URL", nil),
				Description: "Brocade REST API Server URL",
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADE_USERNAME", nil),
				Description: "Brocade Username",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADE_PASSWORD", nil),
				Description: "Brocade Password",
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("BROCADE_VERIFY_SSL", true),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"brocade_pool":             resourcePool(),
			"brocade_rule":             resourceRule(),
			"brocade_traffic_ip_group": resourceTrafficIPGroup(),
			"brocade_virtual_server":   resourceVirtualServer(),
			"brocade_ssl_server_key":   resourceSSLServerKey(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ServerURL: d.Get("server_url").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
		VerifySSL: d.Get("verify_ssl").(bool),
	}

	log.Println("[INFO] Initializing Brocade client")
	return config.Client(), nil
}
