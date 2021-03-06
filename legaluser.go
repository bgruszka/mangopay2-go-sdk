// Copyright 2014 Mathias Monnerville. All rights reserved.
// Use of this source code is governed by a GPL
// license that can be found in the LICENSE file.

package mango

import (
	"encoding/json"
)

// LegalUser describes all the properties of a MangoPay legal user object.
type LegalUser struct {
	User
	Name                                  string
	LegalPersonType                       string
	HeadquartersAddress                   string
	LegalRepresentativeFirstName          string
	LegalRepresentativeLastName           string
	LegalRepresentativeAddress            string
	LegalRepresentativeEmail              string
	LegalRepresentativeBirthday           int64
	LegalRepresentativeNationality        string
	LegalRepresentativeCountryOfResidence string
	Statute                               string
	ProofOfRegistration                   string
	ShareholderDeclaration                string
	service                               *MangoPay // Current service
	wallets                               WalletList
}

func (u *LegalUser) String() string {
	return struct2string(u)
}

// NewLegalUser creates a new legal user.
func (m *MangoPay) NewLegalUser(name string, email string, personType string, legalFirstName, legalLastName string, birthday int64, nationality string, country string) *LegalUser {
	u := &LegalUser{
		Name:                                  name,
		LegalPersonType:                       personType,
		LegalRepresentativeFirstName:          legalFirstName,
		LegalRepresentativeLastName:           legalLastName,
		LegalRepresentativeBirthday:           birthday,
		LegalRepresentativeNationality:        nationality,
		LegalRepresentativeCountryOfResidence: country}
	u.User = User{Email: email}
	u.service = m
	return u
}

// Wallets returns user's wallets.
func (u *LegalUser) Wallets() (WalletList, *RateLimitInfo, error) {
	ws, rateLimitInfo, err := u.service.wallets(u)
	return ws, rateLimitInfo, err
}

// Transfer gets all user's transaction.
func (u *LegalUser) Transfers() (TransferList, *RateLimitInfo, error) {
	return u.service.transfers(u)
}

// Save creates or updates a legal user. The Create API is used
// if the user's Id is an empty string. The Edit API is used when
// the Id is a non-empty string.
func (u *LegalUser) Save() (*RateLimitInfo, error) {
	var action mangoAction
	if u.Id == "" {
		action = actionCreateLegalUser
	} else {
		action = actionEditLegalUser
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
	for _, field := range []string{"LegalRepresentativeBirthday", "CreationDate"} {
		data[field] = int(data[field].(float64))
	}

	// Fields not allowed when creating a user
	if action == actionCreateLegalUser {
		delete(data, "Id")
	}
	delete(data, "CreationDate")

	if action == actionEditLegalUser {
		// Delete empty values so that existing ones don't get
		// overwritten with empty values.
		for k, v := range data {
			switch v.(type) {
			case string:
				if v.(string) == "" {
					delete(data, k)
				}
			case int:
				if v.(int) == 0 {
					delete(data, k)
				}
			}
		}
	}

	ins, rateLimitInfo, err := u.service.anyRequest(new(LegalUser), action, data)
	if err != nil {
		return nil, err
	}
	serv := u.service
	*u = *(ins.(*LegalUser))
	u.service = serv
	return rateLimitInfo, nil
}

// LegalUser finds a legal user using the user_id attribute.
func (m *MangoPay) LegalUser(id string) (*LegalUser, *RateLimitInfo, error) {
	u := new(LegalUser)
	ins, rateLimitInfo, err := m.anyRequest(u, actionFetchLegalUser, JsonObject{"Id": id})
	if err != nil {
		return nil, nil, err
	}
	lu := ins.(*LegalUser)
	lu.service = m
	return lu, rateLimitInfo, nil
}
