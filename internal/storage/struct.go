package storage

import (
	"time"

	"github.com/brocaar/lorawan"
	uuid "github.com/gofrs/uuid"
	"github.com/lib/pq"
	"github.com/lib/pq/hstore"
)

type Device struct {
	DevEUI                    lorawan.EUI64     `db:"dev_eui"`
	CreatedAt                 time.Time         `db:"created_at"`
	UpdatedAt                 time.Time         `db:"updated_at"`
	LastSeenAt                *time.Time        `db:"last_seen_at"`
	ApplicationID             int64             `db:"application_id"`
	DeviceProfileID           uuid.UUID         `db:"device_profile_id"`
	Name                      string            `db:"name"`
	Description               string            `db:"description"`
	SkipFCntCheck             bool              `db:"-"`
	ReferenceAltitude         float64           `db:"-"`
	DeviceStatusBattery       *float32          `db:"device_status_battery"`
	DeviceStatusMargin        *int              `db:"device_status_margin"`
	DeviceStatusExternalPower bool              `db:"device_status_external_power_source"`
	DR                        *int              `db:"dr"`
	Latitude                  *float64          `db:"latitude"`
	Longitude                 *float64          `db:"longitude"`
	Altitude                  *float64          `db:"altitude"`
	DevAddr                   lorawan.DevAddr   `db:"dev_addr"`
	AppSKey                   lorawan.AES128Key `db:"app_s_key"`
	Variables                 hstore.Hstore     `db:"variables"`
	Tags                      hstore.Hstore     `db:"tags"`
	IsDisabled                bool              `db:"-"`
	DataTime                  int64             `db:"data_time"`
	Lat                       float64           `db:"lat"`
	Lng                       float64           `db:"lng"`
	DeviceProfileName         string            `db:"device_profile_name"`
	OrganizationId            int64             `db:"organization_id"`
	DeviceType                int64             `db:"device_type"`
}

type User struct {
	ID            int64
	IsAdmin       bool
	IsActive      bool
	SessionTTL    int32
	CreatedAt     time.Time
	UpdatedAt     time.Time
	PasswordHash  string
	Email         string
	EmailVerified bool
	EmailOld      string
	Note          string
	ExternalID    *string
	WebKey        string
	IosKey        string
	AndroidKey    string
	PhoneNumber   string
	ZoneIDList    pq.Int64Array
	Name          string
	Username      string
	Training      bool
}

type Notification struct {
	Id          int64
	SenderId    int64
	ReceiverId  int64
	Message     string
	CategoryId  int64
	IsRead      bool
	SendTime    time.Time
	ReadTime    time.Time
	DeletedTime time.Time
	SenderIp    string
	ReaderIp    string
	IsDeleted   bool
	DeviceName  string
	DevEui      string
}
