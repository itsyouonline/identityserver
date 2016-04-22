package invitations

import (
	"net/http"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoOrganizationRequestCollectionName = "join-organization-invitations"
)

//InvitationManager is used to store organizations
type InvitationManager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getOrganizationRequestCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoOrganizationRequestCollectionName)
}

//NewInvitationManager creates and initializes a new InvitationManager
func NewInvitationManager(r *http.Request) *InvitationManager {
	session := db.GetDBSession(r)
	return &InvitationManager{
		session:    session,
		collection: getOrganizationRequestCollection(session),
	}
}

// GetByUser gets all invitations for a user.
func (o *InvitationManager) GetByUser(username string) ([]JoinOrganizationInvitation, error) {
	orgRequests := []JoinOrganizationInvitation{}

	err := o.collection.Find(bson.M{"user": username}).All(&orgRequests)

	return orgRequests, err
}

// GetPendingByOrganization gets all pending invitations for a user.
func (o *InvitationManager) GetPendingByOrganization(globalid string) ([]JoinOrganizationInvitation, error) {
	orgRequests := []JoinOrganizationInvitation{}

	err := o.collection.Find(bson.M{"organization": globalid, "status": RequestPending}).All(&orgRequests)

	return orgRequests, err
}

//Get get an invitation by it's content, not really this useful, TODO: just make an exists method
func (o *InvitationManager) Get(username string, organization string, role string, status InvitationStatus) (*JoinOrganizationInvitation, error) {
	var orgRequest JoinOrganizationInvitation

	query := bson.M{
		"user":         username,
		"role":         role,
		"organization": organization,
		"status":       status,
	}

	err := o.collection.Find(query).One(&orgRequest)

	return &orgRequest, err
}

// Save save/update an invitation
func (o *InvitationManager) Save(invite *JoinOrganizationInvitation) error {
	_, err := o.collection.Upsert(
		bson.M{
			"user":         invite.User,
			"organization": invite.Organization,
			"role":         invite.Role,
		}, invite)

	return err
}
