// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
)

// NaturalUser describes all the properties of a MangoPay natural user object.
type NaturalUser struct {
	User
	FirstName, LastName string
	Address             map[string]string
	Birthday            int64
	Nationality         string
	CountryOfResidence  string
	Occupation          string
	IncomeRange         int
	ProofOfIdentity     string
	ProofOfAddress      string
	Capacity            string
	service             *MangoPay // Current service
	wallets             WalletList
}

func (u *NaturalUser) String() string {
	return struct2string(u)
}

// NewNaturalUser creates a new natural user.
func (m *MangoPay) NewNaturalUser(first, last string, email string, birthday int64, nationality, country string) *NaturalUser {
	u := &NaturalUser{
		FirstName:          first,
		LastName:           last,
		Birthday:           birthday,
		Nationality:        nationality,
		CountryOfResidence: country,
	}
	u.User = User{Email: email}
	u.service = m
	return u
}

// Wallets returns user's wallets.
func (u *NaturalUser) Wallets() (WalletList, *RateLimitInfo, error) {
	ws, rateLimitInfo, err := u.service.wallets(u)
	return ws, rateLimitInfo, err
}

// Transfer gets all user's transaction.
func (u *NaturalUser) Transfers() (TransferList, *RateLimitInfo, error) {
	return u.service.transfers(u)
}

// Transfer gets all user's transaction.
func (u *NaturalUser) Transactions() (TransactionList, *RateLimitInfo, error) {
	return u.service.transactions(u)
}

// Save creates or updates a natural user. The Create API is used
// if the user's Id is an empty string. The Edit API is used when
// the Id is a non-empty string.
func (u *NaturalUser) Save() (*RateLimitInfo, error) {
	var action mangoAction
	if u.Id == "" {
		action = actionCreateNaturalUser
	} else {
		action = actionEditNaturalUser
	}

	data := JsonObject{}
	j, err := json.Marshal(u)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(j, &data); err != nil {
		return nil, err
	}

	// Force float64 to int conversion after unmarshalling.
	for _, field := range []string{"Birthday", "CreationDate"} {
		data[field] = int(data[field].(float64))
	}

	// Fields not allowed when creating a user
	if action == actionCreateNaturalUser {
		delete(data, "Id")
	}
	delete(data, "CreationDate")

	if action == actionEditNaturalUser {
		// Delete empty values so that existing ones don't get
		// overwritten with empty values.
		for k, v := range data {
			switch casted := v.(type) {
			case string:
				if casted == "" {
					delete(data, k)
				}
			case int:
				if casted == 0 {
					delete(data, k)
				}
			}
		}
	}
	if action == actionCreateNaturalUser {
		if data["IncomeRange"].(float64) == 0 {
			delete(data, "IncomeRange")
		}
	}

	user, rateLimitInfo, err := u.service.anyRequest(new(NaturalUser), action, data)
	if err != nil {
		return nil, err
	}
	serv := u.service
	*u = *(user.(*NaturalUser))
	u.service = serv
	return rateLimitInfo, nil
}

// NaturalUser finds a natural user using the user_id attribute.
func (m *MangoPay) NaturalUser(id string) (*NaturalUser, *RateLimitInfo, error) {
	u, rateLimitInfo, err := m.anyRequest(new(NaturalUser), actionFetchNaturalUser, JsonObject{"Id": id})
	if err != nil {
		return nil, nil, err
	}
	nu := u.(*NaturalUser)
	nu.service = m
	return nu, rateLimitInfo, nil
}
