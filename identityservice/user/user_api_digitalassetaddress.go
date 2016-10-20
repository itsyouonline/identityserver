package user

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/db/user"
)

// RegisterNewDigitalAssetAddress is the handler for POST /users/{username}/digitalwallet
// Register a new address
func (api UsersAPI) RegisterNewDigitalAssetAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)

	currency := user.DigitalAssetAddress{}

	if err := json.NewDecoder(r.Body).Decode(&currency); err != nil {
		log.Debug("Error while decoding the body: ", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(currency.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	//Check if this label is already used
	_, err = u.GetDigitalAssetAddressByLabel(currency.Label)
	if err == nil {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	if err := userMgr.SaveVirtualCurrency(username, currency); err != nil {
		log.Error("ERROR while saving virtual currency:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u, err = userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	newWallet := u.DigitalWallet

	// respond with created phone number.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newWallet)
}

// GetDigitalWallet is handler for GET /users/{username}/digitalwallet
func (api UsersAPI) GetDigitalWallet(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	userMgr := user.NewManager(r)

	useobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if len(useobj.DigitalWallet) == 0 {
		useobj.DigitalWallet = make([]user.DigitalAssetAddress, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(useobj.DigitalWallet)
}

// GetDigitalAssetAddress It is handler for GET /users/{username}/digitalwallet/{label}
func (api UsersAPI) GetDigitalAssetAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	userobj, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	currency, err := userobj.GetDigitalAssetAddressByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currency)
}

// UpdateDigitalAssetAddress is the handler for PUT /users/{username}/digitalwallet/{label}
// Update the label and/or value of an existing address.
func (api UsersAPI) UpdateDigitalAssetAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	oldlabel := mux.Vars(r)["label"]

	newcurrency := user.DigitalAssetAddress{}
	if err := json.NewDecoder(r.Body).Decode(&newcurrency); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !isValidLabel(newcurrency.Label) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userMgr := user.NewManager(r)

	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = u.GetDigitalAssetAddressByLabel(oldlabel)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if oldlabel != newcurrency.Label {
		_, err = u.GetAddressByLabel(newcurrency.Label)
		if err == nil {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	if err = userMgr.SaveVirtualCurrency(username, newcurrency); err != nil {
		log.Error("ERROR while saving address - ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if oldlabel != newcurrency.Label {
		if err := userMgr.RemoveVirtualCurrency(username, oldlabel); err != nil {
			log.Error(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newcurrency)
}

// DeleteDigitalAssetAddress is the handler for DELETE /users/{username}/digitalwallet/{label}
func (api UsersAPI) DeleteDigitalAssetAddress(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	label := mux.Vars(r)["label"]
	userMgr := user.NewManager(r)

	u, err := userMgr.GetByName(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err = u.GetDigitalAssetAddressByLabel(label)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err = userMgr.RemoveVirtualCurrency(username, label); err != nil {
		log.Error("ERROR while saving address:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
