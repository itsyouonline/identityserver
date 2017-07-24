package organization

import (
	"net/http"
	"strings"

	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/itsyouonline/identityserver/db"
)

const (
	mongoCollectionName       = "organizations"
	logoCollectionName        = "organizationLogos"
	last2FACollectionName     = "last2falogin"
	descriptionCollectionName = "organizationdescriptions"
)

//InitModels initialize models in mongo, if required.
func InitModels() {
	// TODO: Use model tags to ensure indices/constraints.
	index := mgo.Index{
		Key:    []string{"globalid"},
		Unique: true,
	}

	db.EnsureIndex(mongoCollectionName, index)

	// Index the logo collection
	index = mgo.Index{
		Key:    []string{"globalid"},
		Unique: true,
	}

	db.EnsureIndex(logoCollectionName, index)

	index = mgo.Index{
		Key:    []string{"globalid", "username"},
		Unique: true,
	}

	db.EnsureIndex(last2FACollectionName, index)

	// remove last2fa entries after 31 days
	automatic2FAExpiration := mgo.Index{
		Key:         []string{"last2fa"},
		ExpireAfter: time.Second * 3600 * 24 * 31,
		Background:  true,
	}
	db.EnsureIndex(last2FACollectionName, automatic2FAExpiration)

	// Index the info texts
	index = mgo.Index{
		Key:    []string{"globalid"},
		Unique: true,
	}

	db.EnsureIndex(descriptionCollectionName, index)
}

//Manager is used to store organizations
type Manager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

//LogoManager is used to save the logo for an organization
type LogoManager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

//Last2FAManager is used to save the date for the last 2FA login for an organization through the authorization code grant flow
type Last2FAManager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

//DescriptionManager is used to store info texts for an organization
type DescriptionManager struct {
	session    *mgo.Session
	collection *mgo.Collection
}

func getCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, mongoCollectionName)
}

//get the logo collection
func getLogoCollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, logoCollectionName)
}

//get the last 2FA collection
func getLast2FACollection(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, last2FACollectionName)
}

//get the info text collection
func getDescriptionManager(session *mgo.Session) *mgo.Collection {
	return db.GetCollection(session, descriptionCollectionName)
}

//NewManager creates and initializes a new Manager
func NewManager(r *http.Request) *Manager {
	session := db.GetDBSession(r)
	return &Manager{
		session:    session,
		collection: getCollection(session),
	}
}

//NewLogoManager creates and initializes a new LogoManager
func NewLogoManager(r *http.Request) *LogoManager {
	session := db.GetDBSession(r)
	return &LogoManager{
		session:    session,
		collection: getLogoCollection(session),
	}
}

// NewLast2FAManager creates and initializes a new Last2FAManager
func NewLast2FAManager(r *http.Request) *Last2FAManager {
	session := db.GetDBSession(r)
	return &Last2FAManager{
		session:    session,
		collection: getLast2FACollection(session),
	}
}

// NewDescriptionManager creates and initializes a new DescriptionManager
func NewDescriptionManager(r *http.Request) *DescriptionManager {
	session := db.GetDBSession(r)
	return &DescriptionManager{
		session:    session,
		collection: getDescriptionManager(session),
	}
}

// GetOrganizations gets a list of organizations.
func (m *Manager) GetOrganizations(organizationIDs []string) ([]Organization, error) {
	var organizations []Organization

	err := m.collection.Find(bson.M{"globalid": bson.M{"$in": organizationIDs}}).All(&organizations)

	return organizations, err
}

// GetSubOrganizations returns all organizations which have {globalID} as parent (including the organization with {globalID} as globalid)
//TODO: validate globalID since it is appended in the query
//TODO: put an index on the globalid field
func (m *Manager) GetSubOrganizations(globalID string) ([]Organization, error) {
	var organizations = make([]Organization, 0, 0)
	var qry = bson.M{"globalid": bson.M{"$regex": bson.RegEx{"^" + globalID + `\.`, ""}}}
	if err := m.collection.Find(qry).All(&organizations); err != nil {
		return nil, err
	}

	return organizations, nil
}

//isDirectOwner checks if a specific user is in the owners list of an organization
func (m *Manager) isDirectOwner(globalID, username string) (isowner bool, err error) {
	matches, err := m.collection.Find(bson.M{"globalid": globalID, "owners": username}).Count()
	isowner = (matches > 0)
	return
}

func (m *Manager) isOwnerOrMember(globalID, username string, excludelist map[string]bool) (result bool, err error) {
	result, err = m.isDirectOwner(globalID, username)
	if result || err != nil {
		return
	}
	result, err = m.isDirectMember(globalID, username)
	if result || err != nil {
		return
	}
	// If not a direct owner or member, iterate through the list of owner/member organizations
	var org Organization
	err = m.collection.Find(bson.M{"globalid": globalID}).Select(bson.M{"orgowners": 1, "orgmembers": 1, "includesuborgsof": 1}).One(&org)
	if err != nil {
		return
	}
	excludelist[globalID] = true
	// Add the suborganizations from the includelist
	var includesuborgs []string
	for _, owner := range append(org.OrgOwners, org.OrgMembers...) {
		for _, include := range org.IncludeSubOrgsOf {
			if owner == include {
				includesuborgs = append(includesuborgs, include)
			}
		}
	}
	var subOrgs []Organization
	for _, subOrg := range includesuborgs {
		subOrganizations, err := m.GetSubOrganizations(subOrg)
		if err != nil && err != mgo.ErrNotFound {
			return false, err
		}
		subOrgs = append(subOrgs, subOrganizations...)
	}
	var orgs []string
	for _, suborganization := range subOrgs {
		orgs = append(orgs, suborganization.Globalid)
	}
	orgs = append(orgs, org.OrgMembers...)

	for _, owningOrganization := range append(org.OrgOwners, orgs...) {
		if _, excluded := excludelist[owningOrganization]; excluded {
			continue
		}
		result, err = m.isOwnerOrMember(owningOrganization, username, excludelist)
		if result || err != nil {
			return
		}
	}

	return
}

//IsOwner checks if a specific user is in the owners list of an organization
// or belongs to an organization that is in the owner list
// It also checks this for the parentorganizations
func (m *Manager) IsOwner(globalID, username string) (isowner bool, err error) {
	isowner, err = m.isDirectOwner(globalID, username)
	if isowner || err != nil {
		return
	}
	// If not a direct owner, check the ownership in the parent organization
	lastSubOrgSeparator := strings.LastIndex(globalID, ".")
	if lastSubOrgSeparator > 0 {
		isowner, err = m.IsOwner(globalID[:lastSubOrgSeparator], username)
		if isowner || err != nil {
			return
		}
	}
	// If not a direct or inherited owner, iterate through the list of owning organizations
	var org Organization
	err = m.collection.Find(bson.M{"globalid": globalID}).Select(bson.M{"orgowners": 1, "includesuborgsof": 1}).One(&org)
	if err != nil {
		if mgo.ErrNotFound == err {
			err = nil
		}
		return
	}
	excludelist := make(map[string]bool)
	excludelist[globalID] = true

	// Also get the suborganizations if they have been added
	var includesuborgs []string
	for _, owner := range org.OrgOwners {
		for _, include := range org.IncludeSubOrgsOf {
			if owner == include {
				includesuborgs = append(includesuborgs, include)
			}
		}
	}
	var subOrgs []Organization
	for _, subOrg := range includesuborgs {
		subOrganizations, err := m.GetSubOrganizations(subOrg)
		if err != nil && err != mgo.ErrNotFound {
			return false, err
		}
		subOrgs = append(subOrgs, subOrganizations...)
	}
	var orgs []string
	for _, suborganization := range subOrgs {
		orgs = append(orgs, suborganization.Globalid)
	}

	for _, owningOrganization := range append(orgs, org.OrgOwners...) {
		isowner, err = m.isOwnerOrMember(owningOrganization, username, excludelist)
		if isowner || err != nil {
			return
		}
	}

	return
}

//OrganizationIsOwner checks if organization2 is an owner of organization1
func (m *Manager) OrganizationIsOwner(globalID, organization string) (isowner bool, err error) {
	matches, err := m.collection.Find(bson.M{"globalid": globalID, "orgowners": organization}).Count()
	isowner = (matches > 0)
	return
}

//isDirectMember checks if a specific user is in the members list of an organization
func (m *Manager) isDirectMember(globalID, username string) (ismember bool, err error) {
	matches, err := m.collection.Find(bson.M{"globalid": globalID, "members": username}).Count()
	ismember = (matches > 0)
	return
}

//IsMember checks if a specific user is in the members list of an organization
// or belongs to an organization that is in the member list
// it also checks this for the parentorganization
func (m *Manager) IsMember(globalID, username string) (result bool, err error) {
	result, err = m.isDirectMember(globalID, username)
	if result || err != nil {
		return
	}
	// If not a direct member, check the membership in the parent organization
	lastSubOrgSeperator := strings.LastIndex(globalID, ".")
	if lastSubOrgSeperator > 0 {
		result, err = m.IsMember(globalID[:lastSubOrgSeperator], username)
		if result || err != nil {
			return
		}
	}
	// If not a direct or inherited member, iterate through the list of member organizations
	var org Organization
	err = m.collection.Find(bson.M{"globalid": globalID}).Select(bson.M{"orgmembers": 1, "includesuborgsof": 1}).One(&org)
	if err != nil {
		if mgo.ErrNotFound == err {
			err = nil
		}
		return
	}
	excludelist := make(map[string]bool)
	excludelist[globalID] = true

	// Also get the suborganizations if they have been added
	var includesuborgs []string
	for _, member := range org.OrgMembers {
		for _, include := range org.IncludeSubOrgsOf {
			if member == include {
				includesuborgs = append(includesuborgs, include)
			}
		}
	}
	var subOrgs []Organization
	for _, subOrg := range includesuborgs {
		subOrganizations, err := m.GetSubOrganizations(subOrg)
		if err != nil && err != mgo.ErrNotFound {
			return false, err
		}
		subOrgs = append(subOrgs, subOrganizations...)
	}
	var orgs []string
	for _, suborganization := range subOrgs {
		orgs = append(orgs, suborganization.Globalid)
	}
	for _, owningOrganization := range append(org.OrgMembers, orgs...) {
		result, err = m.isOwnerOrMember(owningOrganization, username, excludelist)
		if result || err != nil {
			return
		}
	}

	return
}

//OrganizationIsMember checks if organization2 is a member of organization1
func (m *Manager) OrganizationIsMember(globalID, organization string) (ismember bool, err error) {
	matches, err := m.collection.Find(bson.M{"globalid": globalID, "orgmembers": organization}).Count()
	ismember = (matches > 0)
	return
}

//OrganizationIsPartOf checks if organization2 is a member or an owner of organization1
func (m *Manager) OrganizationIsPartOf(globalID, organization string) (bool, error) {
	condition := []interface{}{
		bson.M{"orgmembers": organization},
		bson.M{"orgowners": organization},
	}
	qry := []interface{}{
		bson.M{"globalid": globalID},
		bson.M{"$or": condition},
	}
	occurrences, err := m.collection.Find(bson.M{"$and": qry}).Count()
	return occurrences == 1, err
}

// AllByUser get organizations for certain user.
func (m *Manager) AllByUser(username string) ([]Organization, error) {
	var organizations []Organization
	//TODO: handle this a bit smarter, select only the ones where the user is owner first, and take select only the org name
	//do the same for the orgs where the username is member but not owners
	//No need to pull in 1000's of records for this

	condition := []interface{}{
		bson.M{"members": username},
		bson.M{"owners": username},
	}

	err := m.collection.Find(bson.M{"$or": condition}).All(&organizations)

	return organizations, err
}

// AllByUserChain returns all organizations where the user is involved, explicitly or implicit
func (m *Manager) AllByUserChain(username string) ([]string, error) {
	var ownedOrgs []string

	// Get organization names where this user is an explicit owner or member
	userOrgs, err := m.AllByUser(username)
	if err != nil {
		if err != mgo.ErrNotFound {
			return nil, err
		}
		err = nil
	}

	// Now load all suborganizations to account for owner- and member-heritance
	var suborganizations []Organization
	for _, org := range userOrgs {
		orgs, err := m.GetSubOrganizations(org.Globalid)
		if err != nil {
			if err != mgo.ErrNotFound {
				return nil, err
			}
			err = nil
		}
		suborganizations = append(suborganizations, orgs...)
	}

	newOrgsFound := true
	var orgs []string

	var parentOrgs []string

	// We only want the globalids. Also remove possible duplicates
	for _, org := range append(userOrgs, suborganizations...) {
		exists := false
		for _, value := range orgs {
			if org.Globalid == value {
				exists = true
				break
			}
		}
		if !exists {
			orgs = append(orgs, org.Globalid)
		}
	}

	// Copy the elements in the slice
	for _, org := range orgs {
		ownedOrgs = append(ownedOrgs, org)
	}

	var orgsFoundThisIteration []string

	// while we find new organizations
	for newOrgsFound {

		// Get the organizations where our currently know organizations are owners or members.
		for _, org := range orgs {
			ownedbyorg, err := m.AllByOrg(org)
			if err != nil {
				if err != mgo.ErrNotFound {
					return nil, err
				}
				err = nil
			}
			for _, obOrg := range ownedbyorg {
				alreadyFound := false
				for _, knownOrg := range ownedOrgs {
					if obOrg.Globalid == knownOrg {
						alreadyFound = true
						break
					}
				}
				if !alreadyFound {
					ownedOrgs = append(ownedOrgs, obOrg.Globalid)
					orgs = append(orgs, obOrg.Globalid)
					orgsFoundThisIteration = append(orgsFoundThisIteration, obOrg.Globalid)
				}
			}
		}

		// Now get the subOrganizations we haven't discovered yet
		for _, org := range orgsFoundThisIteration {
			subOrgs, err := m.GetSubOrganizations(org)
			if err != nil {
				if err != mgo.ErrNotFound {
					return nil, err
				}
				err = nil
			}
			for _, subOrg := range subOrgs {
				alreadyFound := false
				for _, knownOrg := range ownedOrgs {
					if subOrg.Globalid == knownOrg {
						alreadyFound = true
						break
					}
				}
				if !alreadyFound {
					ownedOrgs = append(ownedOrgs, subOrg.Globalid)
					orgs = append(orgs, subOrg.Globalid)
					orgsFoundThisIteration = append(orgsFoundThisIteration, subOrg.Globalid)
				}
			}
		}

		// Get parent orgs
		for _, possibleChild := range orgs {
			subOrgSeperator := strings.LastIndex(possibleChild, ".")
			for subOrgSeperator > 0 {
				alreadyFoud := false
				for _, parentOrg := range parentOrgs {
					if parentOrg == possibleChild[:subOrgSeperator] {
						alreadyFoud = true
						break
					}
				}
				if !alreadyFoud {
					for _, knownOrg := range ownedOrgs {
						if knownOrg == possibleChild[:subOrgSeperator] {
							alreadyFoud = true
							break
						}
					}
				}
				if !alreadyFoud {
					parentOrgs = append(parentOrgs, possibleChild[:subOrgSeperator])
				}
				subOrgSeperator = strings.LastIndex(possibleChild[:subOrgSeperator], ".")
			}
		}

		// Get orgs where parentorgs are owner or member and parentorgs are on the includechildren list.
		for _, parentOrg := range parentOrgs {
			includedOrgs, err := m.AllByOrg(parentOrg)
			if err != nil {
				if err != mgo.ErrNotFound {
					return nil, err
				}
				err = nil
			}
			for _, includedOrg := range includedOrgs {
				alreadyFound := false
				for _, org := range ownedOrgs {
					if includedOrg.Globalid == org {
						alreadyFound = true
						break
					}
				}
				if alreadyFound {
					break
				}
				for _, subOrgsIncludedOf := range includedOrg.IncludeSubOrgsOf {
					if parentOrg == subOrgsIncludedOf {
						ownedOrgs = append(ownedOrgs, includedOrg.Globalid)
						orgs = append(orgs, includedOrg.Globalid)
						orgsFoundThisIteration = append(orgsFoundThisIteration, includedOrg.Globalid)
					}
				}
			}
		}

		newOrgsFound = len(orgsFoundThisIteration) > 0

		// Copy orgsFoundThisIteration to orgs && clear orgsFoundThisIteration
		orgs = []string{}
		for _, org := range orgsFoundThisIteration {
			orgs = append(orgs, org)
		}
		orgsFoundThisIteration = []string{}
	}

	return ownedOrgs, nil
}

// AllByOrg get organizations where certain organization is a member/owner.
func (m *Manager) AllByOrg(globalID string) ([]Organization, error) {
	var organizations []Organization
	//TODO: handle this a bit smarter, select only the ones where the user is owner first, and take select only the org name
	//do the same for the orgs where the username is member but not owners
	//No need to pull in 1000's of records for this

	condition := []interface{}{
		bson.M{"orgmembers": globalID},
		bson.M{"orgowners": globalID},
	}

	err := m.collection.Find(bson.M{"$or": condition}).All(&organizations)

	return organizations, err
}

// Get organization by ID.
func (m *Manager) Get(id string) (*Organization, error) {
	var organization Organization

	objectID := bson.ObjectIdHex(id)

	if err := m.collection.FindId(objectID).One(&organization); err != nil {
		return nil, err
	}

	return &organization, nil
}

// GetByName gets an organization by Name.
func (m *Manager) GetByName(globalID string) (organization *Organization, err error) {
	err = m.collection.Find(bson.M{"globalid": globalID}).One(&organization)
	// Check if organization isnt a nullpointer in case no org was found with this globalID
	if organization != nil && organization.RequiredScopes == nil {
		organization.RequiredScopes = []RequiredScope{}
	}
	return
}

func (m *LogoManager) GetByName(globalID string) (organization *Organization, err error) {
	err = m.collection.Find(bson.M{"globalid": globalID}).One(&organization)
	return
}

// Exists checks if an organization exists.
func (m *Manager) Exists(globalID string) bool {
	count, _ := m.collection.Find(bson.M{"globalid": globalID}).Count()

	return count == 1
}

// Exists checks if an organization and logo entry exists.
func (m *LogoManager) Exists(globalID string) bool {
	count, _ := m.collection.Find(bson.M{"globalid": globalID}).Count()

	return count == 1
}

// Exists checks if an organization - user combination entry exists.
func (m *Last2FAManager) Exists(globalID string, username string) bool {
	condition := []interface{}{
		bson.M{"globalid": globalID},
		bson.M{"username": username},
	}
	count, _ := m.collection.Find(bson.M{"$and": condition}).Count()

	return count == 1
}

// Create a new organization.
func (m *Manager) Create(organization *Organization) error {
	// TODO: Validation!

	err := m.collection.Insert(organization)
	if mgo.IsDup(err) {
		return db.ErrDuplicate
	}
	return err
}

// Create a new organization entry in the organization logo collection
func (m *LogoManager) Create(organization *Organization) error {
	var orgLogo OrganizationLogo

	orgLogo.Globalid = organization.Globalid

	err := m.collection.Insert(orgLogo)
	if mgo.IsDup(err) {
		return db.ErrDuplicate
	}
	return err
}

// SaveMember save or update member
func (m *Manager) SaveMember(organization *Organization, username string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"members": username}})
}

// RemoveMember remove member
func (m *Manager) RemoveMember(organization *Organization, username string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"members": username}})
}

// SaveOwner save or update owners
func (m *Manager) SaveOwner(organization *Organization, owner string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"owners": owner}})
}

// RemoveOwner remove owner
func (m *Manager) RemoveOwner(organization *Organization, owner string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"owners": owner}})
}

// SaveOrgMember save or update organization member
func (m *Manager) SaveOrgMember(organization *Organization, organizationID string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"orgmembers": organizationID}})
}

// RemoveOrgMember remove organization member
func (m *Manager) RemoveOrgMember(organization *Organization, organizationID string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"orgmembers": organizationID}})
}

// SaveOrgOwner save or update owners
func (m *Manager) SaveOrgOwner(organization *Organization, organizationID string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"orgowners": organizationID}})
}

// RemoveOrgOwner remove owner
func (m *Manager) RemoveOrgOwner(organization *Organization, organizationID string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"orgowners": organizationID}})
}

// AddIncludeSubOrgOf adds an organization to the list of orgs who's suborgs are included in the owner/member hierarchy
func (m *Manager) AddIncludeSubOrgOf(globalId, orgMemberId string) error {
	return m.collection.Update(
		bson.M{"globalid": globalId},
		bson.M{"$addToSet": bson.M{"includesuborgsof": orgMemberId}})
}

// RemoveIncludeSubOrgOf removes an organization from the list of orgs who's suborgs are included in the owner/member hierarchy
func (m *Manager) RemoveIncludeSubOrgOf(globalId, orgMemberId string) error {
	return m.collection.Update(
		bson.M{"globalid": globalId},
		bson.M{"$pull": bson.M{"includesuborgsof": orgMemberId}})
}

func (m *Manager) AddDNS(organization *Organization, dnsName string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"dns": dnsName}})
}

func (m *Manager) UpdateDNS(organization *Organization, oldDNSName string, newDNSName string) error {
	err := m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"dns": oldDNSName}})
	if err != nil {
		return err
	}
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$addToSet": bson.M{"dns": newDNSName}})
}

// RemoveDNS remove DNS
func (m *Manager) RemoveDNS(organization *Organization, dns string) error {
	return m.collection.Update(
		bson.M{"globalid": organization.Globalid},
		bson.M{"$pull": bson.M{"dns": dns}})
}

// Remove removes the organization
func (m *Manager) Remove(globalid string) error {
	return m.collection.Remove(bson.M{"globalid": globalid})
}

// Remove the organization logo
func (m *LogoManager) Remove(globalid string) error {
	return m.collection.Remove(bson.M{"globalid": globalid})
}

// Remove the Last2FA entries for this organization
func (m *Last2FAManager) RemoveByOrganization(globalid string) error {
	_, err := m.collection.RemoveAll(bson.M{"globalid": globalid})
	return err
}

//Remove the Last2FA entries for this user
func (m *Last2FAManager) RemoveByUser(username string) error {
	_, err := m.collection.RemoveAll(bson.M{"username": username})
	return err
}

// Remove removes the organization descriptions
func (m *DescriptionManager) Remove(globalid string) error {
	return m.collection.Remove(bson.M{"globalid": globalid})
}

// UpdateMembership Updates a user his role in an organization
func (m *Manager) UpdateMembership(globalid string, username string, oldrole string, newrole string) error {
	qry := bson.M{"globalid": globalid}
	pull := bson.M{
		"$pull": bson.M{oldrole: username},
	}
	push := bson.M{
		"$addToSet": bson.M{newrole: username},
	}
	err := m.collection.Update(qry, pull)
	if err != nil {
		return err
	}
	return m.collection.Update(qry, push)
}

// UpdateOrgMembership Updates an organization role in another organization
func (m *Manager) UpdateOrgMembership(globalid string, organization string, oldrole string, newrole string) error {
	qry := bson.M{"globalid": globalid}
	pull := bson.M{
		"$pull": bson.M{oldrole: organization},
	}
	push := bson.M{
		"$addToSet": bson.M{newrole: organization},
	}
	err := m.collection.Update(qry, pull)
	if err != nil {
		return err
	}
	return m.collection.Update(qry, push)
}

// CountByUser counts the amount of organizations by user
func (m *Manager) CountByUser(username string) (int, error) {
	qry := bson.M{"owners": username}
	return m.collection.Find(qry).Count()
}

// CountByOrganization counts the amount of organizations where the organization is an owner
func (m *Manager) CountByOrganization(organization string) (int, error) {
	qry := bson.M{"orgowners": organization}
	return m.collection.Find(qry).Count()
}

// RemoveUser Removes a user from an organization
func (m *Manager) RemoveUser(globalID string, username string) error {
	qry := bson.M{"globalid": globalID}
	update := bson.M{"$pull": bson.M{"owners": username, "members": username}}
	return m.collection.Update(qry, update)
}

// RemoveOrganization Removes an organization as member or owner from another organization
func (m *Manager) RemoveOrganization(globalID string, organization string) error {
	qry := bson.M{"globalid": globalID}
	update := bson.M{"$pull": bson.M{"orgowners": organization, "orgmembers": organization}}
	return m.collection.Update(qry, update)
}

// GetValidity gets the 2FA validity duration in seconds
func (m *Manager) GetValidity(globalID string) (int, error) {
	var org Organization
	err := m.collection.Find(bson.M{"globalid": globalID}).One(&org)
	if err != nil {
		return 0, err
	}
	seconds := org.SecondsValidity
	if seconds == -1 { //special value to avoid confusion with mongo null
		return 0, err
	} else if seconds == 0 { //mongo null for int field, use default duration
		return 3600 * 24 * 7, err //7 days default duration
	} else {
		return seconds, err
	}
}

func (m *Manager) SetValidity(globalID string, secondsDuration int) error {
	if secondsDuration == 0 {
		secondsDuration = -1 //assign -1 if duration should be zero to avoid confusion with mongo null
	}
	return m.collection.Update(
		bson.M{"globalid": globalID},
		bson.M{"$set": bson.M{"secondsvalidity": secondsDuration}})
}

// SaveLogo save or update logo
func (m *LogoManager) SaveLogo(globalID string, logo string) (*mgo.ChangeInfo, error) {
	return m.collection.Upsert(
		bson.M{"globalid": globalID},
		bson.M{"$set": bson.M{"logo": logo}})
}

// GetLogo Gets the logo from an organization
func (m *LogoManager) GetLogo(globalID string) (string, error) {
	var org *OrganizationLogo
	err := m.collection.Find(bson.M{"globalid": globalID}).One(&org)
	if err != nil {
		return "", err
	}
	return org.Logo, err
}

// RemoveLogo Removes the logo from an organization
func (m *LogoManager) RemoveLogo(globalID string) error {
	qry := bson.M{"globalid": globalID}
	update := bson.M{"$unset": bson.M{"logo": 1}}
	return m.collection.Update(qry, update)
}

// SetLast2FA Set the last successful 2FA time
func (m *Last2FAManager) SetLast2FA(globalID string, username string) error {
	now := time.Now()
	condition := []interface{}{
		bson.M{"globalid": globalID},
		bson.M{"username": username},
	}
	_, err := m.collection.Upsert(
		bson.M{"$and": condition},
		bson.M{"$set": bson.M{"last2fa": now}})
	return err
}

// GetLast2FA Gets the date of the last successful 2FA login, if no failed login attempts have occurred since then
func (m *Last2FAManager) GetLast2FA(globalID string, username string) (db.DateTime, error) {
	var l2fa *UserLast2FALogin
	condition := []interface{}{
		bson.M{"globalid": globalID},
		bson.M{"username": username},
	}
	err := m.collection.Find(bson.M{"$and": condition}).One(&l2fa)
	return l2fa.Last2FA, err
}

// RemoveLast2FA Removes the entry of the last successful 2FA login for this organization - user combination
func (m *Last2FAManager) RemoveLast2FA(globalID string, username string) error {
	condition := []interface{}{
		bson.M{"globalid": globalID},
		bson.M{"username": username},
	}
	return m.collection.Remove(bson.M{"$and": condition})
}

// AddRequiredScope adds a required scope
func (m *Manager) AddRequiredScope(globalId string, requiredScope RequiredScope) error {
	qry := bson.M{"globalid": globalId}
	update := bson.M{"$push": bson.M{"requiredscopes": requiredScope}}
	return m.collection.Update(qry, update)
}

// UpdateRequiredScope updates a required scope
func (m *Manager) UpdateRequiredScope(globalId string, oldRequiredScope string, newRequiredScope RequiredScope) error {
	qry := bson.M{
		"globalid":             globalId,
		"requiredscopes.scope": oldRequiredScope,
	}
	update := bson.M{
		"$set": bson.M{
			"requiredscopes.$": newRequiredScope,
		},
	}
	return m.collection.Update(qry, update)
}

// DeleteRequiredScope deletes a required scope
func (m *Manager) DeleteRequiredScope(globalId string, requiredScope string) error {
	return m.collection.Update(bson.M{"globalid": globalId},
		bson.M{"$pull": bson.M{"requiredscopes": bson.M{"scope": requiredScope}}})
}

func (m *Manager) ListByUserOrGlobalID(username string, globalIds []string) (error, []Organization) {
	var organizations []Organization
	qry := bson.M{
		"$or": []bson.M{
			{"owners": username},
			{"members": username},
			{"globalid": bson.M{
				"$in": globalIds},
			},
		},
	}
	err := m.collection.Find(qry).All(&organizations)
	return err, organizations
}

// TransferSubOrg transfers the suborganization, and all of its suborganizations
// to a new parent organization
func (m *Manager) TransferSubOrg(globalId string, newParent string) error {
	suborgseparator := strings.LastIndex(globalId, ".")
	suborgs, err := m.GetSubOrganizations(globalId)
	if err != nil {
		return err
	}

	suborgids := make([]string, len(suborgs))
	for i, subOrg := range suborgs {
		suborgids[i] = subOrg.Globalid
	}

	for _, org := range append(suborgids, globalId) {
		newGlobalid := newParent + org[suborgseparator:]
		qry := bson.M{
			"globalid": org,
		}
		update := bson.M{
			"$set": bson.M{
				"globalid": newGlobalid,
			},
		}
		err = m.collection.Update(qry, update)
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveDescription saves a description for an organization
func (m *DescriptionManager) SaveDescription(globalId string, text LocalizedInfoText) error {
	_, err := m.collection.Upsert(bson.M{"globalid": globalId}, bson.M{"$addToSet": bson.M{"infotexts": text}})
	return err
}

// UpdateDescription updates a description for an organization
func (m *DescriptionManager) UpdateDescription(globalId string, text LocalizedInfoText) error {
	err := m.collection.Update(bson.M{"globalid": globalId}, bson.M{"$pull": bson.M{"infotexts": bson.M{"langkey": text.LangKey}}})
	if err != nil {
		return err
	}
	_, err = m.collection.Upsert(bson.M{"globalid": globalId}, bson.M{"$addToSet": bson.M{"infotexts": text}})
	return err
}

// DeleteDescription deletes a (translated) description for an organization
func (m *DescriptionManager) DeleteDescription(globalId, langKey string) error {
	return m.collection.Update(bson.M{"globalid": globalId}, bson.M{"$pull": bson.M{"infotexts": bson.M{"langkey": langKey}}})
}

// GetDescription get all descriptions for an organization
func (m *DescriptionManager) GetDescription(globalId string) (OrganizationInfoText, error) {
	var info OrganizationInfoText
	err := m.collection.Find(bson.M{"globalid": globalId}).One(&info)
	return info, err
}
