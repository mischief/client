// common_test.go - Tests for common code of the noise based wire protocol.
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

package vault

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVaultOpenSeal(t *testing.T) {
	assert := assert.New(t)

	tmpfile, err := ioutil.TempFile("", "example")
	assert.NoError(err, "TempFile failed")
	v1 := Vault{
		Passphrase: "up up down down left right right left",
		Path:       tmpfile.Name(),
	}
	plaintext1 := "war is peace freedom is slavery ignorance is strength"
	err = v1.Seal([]byte(plaintext1))
	assert.NoError(err, "Vault Seal failed")
	plaintext2, err := v1.Open()
	assert.NoError(err, "Vault Open failed")
	assert.Equal(plaintext1, string(plaintext2))
	os.Remove(tmpfile.Name())
}