// daemon.go - management of client services
// Copyright (C) 2017  David Anthony Stainton
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package util provides client utilities
package util

import (
	"github.com/katzenpost/core/crypto/rand"
	"github.com/katzenpost/core/pki"
	"github.com/katzenpost/core/wire/server"
	"github.com/op/go-logging"
)

const (
	DefaultSMTPNetwork = "tcp"
	DefaultSMTPAddress = "127.0.0.1:2525"
	DefaultPOP3Network = "tcp"
	DefaultPOP3Address = "127.0.0.1:1110"
)

var log = logging.MustGetLogger("mixclient")

// ClientDaemon handles the startup and shutdown of all client services
type ClientDaemon struct {
	config     *Config
	passphrase string
	keysDir    string
	userPKI    UserPKI
	mixPKI     pki.Client
}

// NewClientDaemon creates a new ClientDaemon given a Config
func NewClientDaemon(config *Config, passphrase string, keysDirPath string, userPKI UserPKI, mixPKI pki.Client) (*ClientDaemon, error) {
	d := ClientDaemon{
		config:     config,
		passphrase: passphrase,
		keysDir:    keysDirPath,
		userPKI:    userPKI,
		mixPKI:     mixPKI,
	}
	return &d, nil
}

// Start starts the client services
// which proxy message to and from the mixnet
// via POP3 and SMTP
func (c *ClientDaemon) Start() error {
	var smtpServer, pop3Server *server.Server
	log.Debug("Client startup.")
	log.Debug("starting smtp proxy service")
	smtpProxy := NewSubmitProxy(c.config, rand.Reader, c.userPKI)
	if len(c.config.SMTPProxy.Network) == 0 {
		smtpServer = server.New(DefaultSMTPNetwork, DefaultSMTPAddress, smtpProxy.handleSMTPSubmission, nil)
	} else {
		smtpServer = server.New(c.config.SMTPProxy.Network, c.config.SMTPProxy.Address, smtpProxy.handleSMTPSubmission, nil)
	}
	err := smtpServer.Start()
	if err != nil {
		return err
	}

	log.Debug("starting pop3 proxy service")
	pop3Proxy := NewPop3Proxy()
	if len(c.config.POP3Proxy.Network) == 0 {
		pop3Server = server.New(DefaultPOP3Network, DefaultPOP3Address, pop3Proxy.handleConnection, nil)
	} else {
		pop3Server = server.New(c.config.POP3Proxy.Network, c.config.POP3Proxy.Address, pop3Proxy.handleConnection, nil)
	}
	err = pop3Server.Start()
	if err != nil {
		return err
	}

	providerPool, err := FromAccounts(c.config.Account, c.config, c.keysDir, c.passphrase, c.mixPKI)
	if err != nil {
		return err
	}
	log.Debugf("provider pool %v", providerPool)
	return nil
}

// Stop stops the client services
func (c *ClientDaemon) Stop() {
	log.Debug("Client shutdown.")
}
