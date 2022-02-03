// Package session contains business object (BO) and data access object (DAO) implementations for Session.
package session

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/btnguyen2k/consu/reddo"

	"github.com/btnguyen2k/henge"

	"main/src/utils"
)

// NewSession is helper function to create new Session bo.
func NewSession(appVersion uint64, id, sessionType, idSource, appId, userId, sessionData string, expiry time.Time) *Session {
	sess := &Session{
		UniversalBo: henge.NewUniversalBo(id, appVersion),
		sessionData: strings.TrimSpace(sessionData),
		idSource:    strings.TrimSpace(strings.ToLower(idSource)),
		appId:       strings.TrimSpace(strings.ToLower(appId)),
		userId:      strings.TrimSpace(strings.ToLower(userId)),
		sessionType: strings.TrimSpace(sessionType),
		expiry:      expiry,
	}
	if sess.GetId() == "" {
		sess.SetId(utils.UniqueId())
	}
	return sess.sync()
}

// NewSessionFromUbo is helper function to create new Session bo from a universal bo.
func NewSessionFromUbo(ubo *henge.UniversalBo) *Session {
	if ubo == nil {
		return nil
	}
	sess := Session{UniversalBo: &henge.UniversalBo{}}
	if err := json.Unmarshal([]byte(ubo.GetDataJson()), &sess); err != nil {
		log.Print(fmt.Sprintf("[WARN] NewSessionFromUbo - error unmarshalling JSON data: %e", err))
		// log.Print(err)
		return nil
	}
	sess.UniversalBo = ubo.Clone()
	if sessionType, err := sess.GetExtraAttrAs(FieldSession_SessionType, reddo.TypeString); err == nil {
		sess.sessionType = sessionType.(string)
	}
	if idSource, err := sess.GetExtraAttrAs(FieldSession_IdSource, reddo.TypeString); err == nil {
		sess.idSource = idSource.(string)
	}
	if appId, err := sess.GetExtraAttrAs(FieldSession_AppId, reddo.TypeString); err == nil {
		sess.appId = appId.(string)
	}
	if userId, err := sess.GetExtraAttrAs(FieldSession_UserId, reddo.TypeString); err == nil {
		sess.userId = userId.(string)
	}
	if expiry, err := sess.GetExtraAttrAsTimeWithLayout(FieldSession_Expiry, henge.TimeLayout); err == nil {
		sess.expiry = expiry
	}
	if data, err := sess.GetDataAttrAs(AttrSession_Data, reddo.TypeString); err == nil && data != nil {
		sess.sessionData = data.(string)
	}
	return &sess
}

const (
	FieldSession_IdSource    = "isrc"
	FieldSession_AppId       = "aid"
	FieldSession_UserId      = "uid"
	FieldSession_SessionType = "type"
	FieldSession_Expiry      = "eat"

	AttrSession_Ubo  = "_ubo"
	AttrSession_Data = "data"
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
		AttrSession_Ubo:          sess.UniversalBo.Clone(),
		FieldSession_IdSource:    sess.idSource,
		FieldSession_AppId:       sess.appId,
		FieldSession_UserId:      sess.userId,
		FieldSession_SessionType: sess.sessionType,
		FieldSession_Expiry:      sess.expiry,
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
	var err error
	if m[AttrSession_Ubo] != nil {
		js, _ := json.Marshal(m[AttrSession_Ubo])
		if err := json.Unmarshal(js, &sess.UniversalBo); err != nil {
			return err
		}
	}
	if sess.idSource, err = reddo.ToString(m[FieldSession_IdSource]); err != nil {
		return err
	}
	if sess.appId, err = reddo.ToString(m[FieldSession_AppId]); err != nil {
		return err
	}
	if sess.userId, err = reddo.ToString(m[FieldSession_UserId]); err != nil {
		return err
	}
	if sess.sessionType, err = reddo.ToString(m[FieldSession_SessionType]); err != nil {
		return err
	}
	if sessionData, err := sess.GetDataAttrAs(AttrSession_Data, reddo.TypeString); err != nil {
		return err
	} else if sessionData != nil {
		sess.sessionData = sessionData.(string)
	}
	if sess.expiry, err = reddo.ToTimeWithLayout(m[FieldSession_Expiry], henge.TimeLayout); err != nil {
		return err
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
	sess.expiry = value
	return sess
}

// IsExpired returns true if the session expired, false otherwise.
func (sess *Session) IsExpired() bool {
	return sess.expiry.Before(time.Now())
}

func (sess *Session) sync() *Session {
	sess.SetExtraAttr(FieldSession_IdSource, sess.idSource)
	sess.SetExtraAttr(FieldSession_AppId, sess.appId)
	sess.SetExtraAttr(FieldSession_UserId, sess.userId)
	sess.SetExtraAttr(FieldSession_Expiry, sess.expiry)
	sess.SetExtraAttr(FieldSession_SessionType, sess.sessionType)
	sess.SetDataAttr(AttrSession_Data, sess.sessionData)
	sess.UniversalBo.Sync()
	return sess
}
