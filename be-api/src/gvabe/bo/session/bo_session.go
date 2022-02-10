// Package session contains business object (BO) and data access object (DAO) implementations for Session.
package session

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"main/src/gvabe/bo"

	"github.com/btnguyen2k/henge"

	"main/src/utils"
)

// NewSession is helper function to create new Session bo.
func NewSession(appVersion uint64, id, sessionType, idSource, appId, userId, sessionData string, expiry time.Time) *Session {
	sess := &Session{
		UniversalBo: henge.NewUniversalBo(id, appVersion, henge.UboOpt{TimeLayout: bo.UboTimeLayout, TimestampRounding: bo.UboTimestampRounding}),
	}
	if sess.GetId() == "" {
		sess.SetId(utils.UniqueId())
	}
	sess.
		SetSessionData(sessionData).
		SetIdSource(idSource).
		SetAppId(appId).
		SetUserId(userId).
		SetSessionType(sessionType).
		SetExpiry(expiry)
	return sess.sync()
}

// NewSessionFromUbo is helper function to create new Session bo from a universal bo.
func NewSessionFromUbo(ubo *henge.UniversalBo) *Session {
	if ubo == nil {
		return nil
	}
	ubo = ubo.Clone()
	sess := &Session{UniversalBo: ubo}

	if v, err := ubo.GetExtraAttrAsTimeWithLayout(FieldSessionExpiry, bo.UboTimeLayout); err != nil {
		return nil
	} else {
		sess.SetExpiry(v)
	}

	fieldListStr := []string{FieldSessionIdSource, FieldSessionAppId, FieldSessionUserId, FieldSessionSessionType}
	setterListStr := []func(string) *Session{sess.SetIdSource, sess.SetAppId, sess.SetUserId, sess.SetSessionType}
	for i, field := range fieldListStr {
		if v, err := ubo.GetExtraAttrAs(field, reddo.TypeString); err != nil {
			return nil
		} else if v != nil {
			setterListStr[i](v.(string))
		}
	}

	attrListStr := []string{AttrSessionData}
	setterListStr = []func(string) *Session{sess.SetSessionData}
	for i, attr := range attrListStr {
		if v, err := ubo.GetDataAttrAs(attr, reddo.TypeString); err != nil {
			return nil
		} else if v != nil {
			setterListStr[i](v.(string))
		}
	}

	return sess.sync()
}

const (
	FieldSessionIdSource    = "isrc"
	FieldSessionAppId       = "aid"
	FieldSessionUserId      = "uid"
	FieldSessionSessionType = "type"
	FieldSessionExpiry      = "eat"

	AttrSessionUbo  = "_ubo"
	AttrSessionData = "data"
)

// Session is the business object.
// Session inherits unique id from bo.UniversalBo.
type Session struct {
	*henge.UniversalBo `json:"_ubo"`
	sessionData        string    `json:"data"`
	idSource           string    `json:"isrc"` // identity source
	appId              string    `json:"aid"`  // id of application that is owner of the session
	userId             string    `json:"uid"`  // id of user that is owner of the session
	sessionType        string    `json:"type"` // session type
	expiry             time.Time `json:"eat"`  // timestamp when the session expires
}

// MarshalJSON implements json.encode.Marshaler.MarshalJSON.
//	TODO: lock for read?
func (sess *Session) MarshalJSON() ([]byte, error) {
	sess.sync()
	m := map[string]interface{}{
		AttrSessionUbo: sess.UniversalBo.Clone(),
		bo.SerKeyFields: map[string]interface{}{
			FieldSessionIdSource:    sess.GetIdSource(),
			FieldSessionAppId:       sess.GetAppId(),
			FieldSessionUserId:      sess.GetUserId(),
			FieldSessionSessionType: sess.GetSessionType(),
			FieldSessionExpiry:      sess.GetExpiry(),
		},
		bo.SerKeyAttrs: map[string]interface{}{
			AttrSessionData: sess.GetSessionData(),
		},
	}
	return json.Marshal(m)
}

// UnmarshalJSON implements json.decode.Unmarshaler.UnmarshalJSON.
//	TODO: lock for write?
func (sess *Session) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if m[AttrSessionUbo] != nil {
		js, _ := json.Marshal(m[AttrSessionUbo])
		if err := json.Unmarshal(js, &sess.UniversalBo); err != nil {
			return err
		}
	}
	if _cols, ok := m[bo.SerKeyFields].(map[string]interface{}); ok {
		fieldListStr := []string{FieldSessionIdSource, FieldSessionAppId, FieldSessionUserId, FieldSessionSessionType}
		setterListStr := []func(string) *Session{sess.SetIdSource, sess.SetAppId, sess.SetUserId, sess.SetSessionType}
		for i, field := range fieldListStr {
			if v, err := reddo.ToString(_cols[field]); err != nil {
				return err
			} else {
				setterListStr[i](v)
			}
		}

		if v, err := reddo.ToTimeWithLayout(_cols[FieldSessionExpiry], bo.UboTimeLayout); err != nil {
			return err
		} else {
			sess.SetExpiry(v)
		}
	}
	if _attrs, ok := m[bo.SerKeyAttrs].(map[string]interface{}); ok {
		attrListStr := []string{AttrSessionData}
		setterListStr := []func(string) *Session{sess.SetSessionData}
		for i, attr := range attrListStr {
			if v, err := reddo.ToString(_attrs[attr]); err != nil {
				return err
			} else {
				setterListStr[i](v)
			}
		}
	}

	sess.sync()
	return nil
}

// GetIdSource returns session's 'id-source' value.
func (sess *Session) GetIdSource() string {
	return sess.idSource
}

// SetIdSource sets session's 'id-source' value.
func (sess *Session) SetIdSource(value string) *Session {
	sess.idSource = strings.TrimSpace(strings.ToLower(value))
	return sess
}

// GetAppId returns session's 'app-id' value.
func (sess *Session) GetAppId() string {
	return sess.appId
}

// SetAppId sets session's 'app-id' value.
func (sess *Session) SetAppId(value string) *Session {
	sess.appId = strings.TrimSpace(strings.ToLower(value))
	return sess
}

// GetUserId returns session's 'user-id' value.
func (sess *Session) GetUserId() string {
	return sess.userId
}

// SetUserId sets session's 'user-id' value.
func (sess *Session) SetUserId(value string) *Session {
	sess.userId = strings.TrimSpace(strings.ToLower(value))
	return sess
}

// GetSessionType returns session's 'session-type' value.
func (sess *Session) GetSessionType() string {
	return sess.sessionType
}

// SetSessionType sets session's 'session-type' value.
func (sess *Session) SetSessionType(value string) *Session {
	sess.sessionType = strings.TrimSpace(value)
	return sess
}

// GetSessionData returns session's 'session-data' value.
func (sess *Session) GetSessionData() string {
	return sess.sessionData
}

// SetSessionData sets session's 'session-data' value.
func (sess *Session) SetSessionData(value string) *Session {
	sess.sessionData = strings.TrimSpace(value)
	return sess
}

// GetExpiry returns session's 'expiry' value.
func (sess *Session) GetExpiry() time.Time {
	return sess.expiry
}

// SetExpiry sets session's 'expiry' value.
func (sess *Session) SetExpiry(value time.Time) *Session {
	sess.expiry = sess.RoundTimestamp(value)
	return sess
}

// IsExpired returns true if the session expired, false otherwise.
func (sess *Session) IsExpired() bool {
	return sess.expiry.Before(time.Now())
}

func (sess *Session) sync() *Session {
	sess.SetExtraAttr(FieldSessionIdSource, sess.idSource)
	sess.SetExtraAttr(FieldSessionAppId, sess.appId)
	sess.SetExtraAttr(FieldSessionUserId, sess.userId)
	sess.SetExtraAttr(FieldSessionExpiry, sess.expiry)
	sess.SetExtraAttr(FieldSessionSessionType, sess.sessionType)
	sess.SetDataAttr(AttrSessionData, sess.sessionData)
	sess.UniversalBo.Sync()
	return sess
}
