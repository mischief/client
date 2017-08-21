// db_test.go - db tests
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

package egress

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/katzenpost/core/sphinx/constants"
	"github.com/stretchr/testify/require"
)

func TestDBBasics(t *testing.T) {
	require := require.New(t)

	dbFile, err := ioutil.TempFile("", "db_test1")
	require.NoError(err, "unexpected TempFile error")
	defer func() {
		err := os.Remove(dbFile.Name())
		require.NoError(err, "unexpected os.Remove error")
	}()
	store, err := New(dbFile.Name())
	require.NoError(err, "unexpected New() error")

	rid := []byte{1, 2, 3, 4}
	recipientID := [constants.RecipientIDLength]byte{}
	copy(recipientID[:], rid)
	s := StorageBlock{
		SenderProvider:    "acme.com",
		RecipientProvider: "nsa.gov",
		RecipientID:       &recipientID,
		Payload:           []byte(`"The time has come," the Walrus said`),
	}
	id := []byte{1, 2, 3, 4, 5, 6}
	surbID := [constants.SURBIDLength]byte{}
	copy(surbID[:], id)

	err = store.Put(&surbID, &s)
	require.NoError(err, "unexpected storeMessage() error")

	surbs, err := store.GetKeys()
	require.NoError(err, "unexpected GetKeys() error")

	err = store.Remove(&surbs[0])
	require.NoError(err, "unexpected Remove() error")

	surbs, err = store.GetKeys()
	require.NoError(err, "unexpected GetKeys() error")

	require.Equal(len(surbs), 0, "expected zero length")

	err = store.Close()
	require.NoError(err, "unexpected Close() error")
}