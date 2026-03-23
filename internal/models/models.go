package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ---------------------------------------------------------------------------
// GSTHierarchy
// ---------------------------------------------------------------------------

// GSTHierarchy represents a single GST zone/commissionerate/division/range record.
// Unique key: gstrangeCode
// GSTHierarchyId: client-facing business ID (provided by client, unique)
type GSTHierarchy struct {
	ID                     primitive.ObjectID `bson:"_id,omitempty"            json:"id,omitempty"`
	GSTHierarchyId         string             `bson:"gstHierarchyId"           json:"gstHierarchyId"           binding:"required"`
	GSTZoneName            string             `bson:"gstZoneName"              json:"gstZoneName"              binding:"required"`
	GSTZoneCode            string             `bson:"gstZoneCode"              json:"gstZoneCode"              binding:"required"`
	GSTCommissionerateName string             `bson:"gstCommissionerateName"   json:"gstCommissionerateName"   binding:"required"`
	GSTCommissionerateCode string             `bson:"gstCommissionerateCode"   json:"gstCommissionerateCode"   binding:"required"`
	GSTDivisionCode        string             `bson:"gstDivisionCode"          json:"gstDivisionCode"          binding:"required"`
	GSTDivisionName        string             `bson:"gstDivisionName"          json:"gstDivisionName"          binding:"required"`
	GSTRangeName           string             `bson:"gstrangeName"             json:"gstrangeName"             binding:"required"`
	GSTRangeCode           string             `bson:"gstrangeCode"             json:"gstrangeCode"             binding:"required"`
	RangeStateName         string             `bson:"rangestateName"           json:"rangestateName"`
	RangePincode           string             `bson:"rangepincode"             json:"rangepincode"`
	IsDeleted              bool               `bson:"isDeleted"                json:"isDeleted"`
	DeletedAt              *time.Time         `bson:"deletedAt,omitempty"      json:"deletedAt,omitempty"`
	CreatedAt              time.Time          `bson:"createdAt"                json:"createdAt"`
	UpdatedAt              time.Time          `bson:"updatedAt"                json:"updatedAt"`
}

// ---------------------------------------------------------------------------
// Premise
// ---------------------------------------------------------------------------

// Premise represents a manufacturing premise / factory.
// Unique key: premiseCode
// PremiseId: client-facing business ID (provided by client, unique)
// GSTHierarchyId: reference to the parent GST hierarchy record
type Premise struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"       json:"id,omitempty"`
	PremiseId        string             `bson:"premiseId"           json:"premiseId"           binding:"required"`
	GSTHierarchyId   string             `bson:"gstHierarchyId"      json:"gstHierarchyId"      binding:"required"`
	ManufacturerName string             `bson:"manufacturerName"    json:"manufacturerName"    binding:"required"`
	PremiseName      string             `bson:"premiseName"         json:"premiseName"         binding:"required"`
	PremiseCode      string             `bson:"premiseCode"         json:"premiseCode"         binding:"required"`
	GSTIN            string             `bson:"gstin"               json:"gstin"               binding:"required"`
	PremiseAddress   string             `bson:"premiseAddress"      json:"premiseAddress"`
	PremiseCity      string             `bson:"premiseCity"         json:"premiseCity"`
	PremisePincode   string             `bson:"premisePincode"      json:"premisePincode"`
	PremiseLatLong   string             `bson:"premiseLatLong"      json:"premiseLatLong"`
	IsDeleted        bool               `bson:"isDeleted"           json:"isDeleted"`
	DeletedAt        *time.Time         `bson:"deletedAt,omitempty" json:"deletedAt,omitempty"`
	CreatedAt        time.Time          `bson:"createdAt"           json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt"           json:"updatedAt"`
}

// ---------------------------------------------------------------------------
// Machine
// ---------------------------------------------------------------------------

// Machine represents a manufacturing machine registered at a premise.
// Unique key: machineRegistrationNo
// MachineId: client-facing business ID (provided by client, unique)
// PremiseId: reference to the parent premise record
// GSTHierarchyId: reference to the parent GST hierarchy record
type Machine struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty"              json:"id,omitempty"`
	MachineId             string             `bson:"machineId"                  json:"machineId"                  binding:"required"`
	PremiseId             string             `bson:"premiseId"                  json:"premiseId"                  binding:"required"`
	GSTHierarchyId        string             `bson:"gstHierarchyId"             json:"gstHierarchyId"             binding:"required"`
	MachineName           string             `bson:"machineName"                json:"machineName"                binding:"required"`
	MachineRegistrationNo string             `bson:"machineRegistrationNo"      json:"machineRegistrationNo"      binding:"required"`
	MachineType           string             `bson:"machineType"                json:"machineType"`
	MachineMake           string             `bson:"machineMake"                json:"machineMake"`
	MachineModel          string             `bson:"machineModel"               json:"machineModel"`
	MachineSerialNo       string             `bson:"machineSerialNo"            json:"machineSerialNo"`
	WorkingStatus         string             `bson:"workingStatus"              json:"workingStatus"`
	IsDeleted             bool               `bson:"isDeleted"                  json:"isDeleted"`
	DeletedAt             *time.Time         `bson:"deletedAt,omitempty"        json:"deletedAt,omitempty"`
	CreatedAt             time.Time          `bson:"createdAt"                  json:"createdAt"`
	UpdatedAt             time.Time          `bson:"updatedAt"                  json:"updatedAt"`
}

// ---------------------------------------------------------------------------
// Officer
// ---------------------------------------------------------------------------

// Officer represents a government officer.
// Unique key: officerCode
// OfficerId: client-facing business ID (provided by client, unique)
// GSTHierarchyId: reference to the parent GST hierarchy record
type Officer struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"        json:"id,omitempty"`
	OfficerId          string             `bson:"officerId"            json:"officerId"            binding:"required"`
	GSTHierarchyId     string             `bson:"gstHierarchyId"       json:"gstHierarchyId"       binding:"required"`
	OfficerName        string             `bson:"officerName"          json:"officerName"          binding:"required"`
	OfficerCode        string             `bson:"officerCode"          json:"officerCode"          binding:"required"`
	OfficerDesignation string             `bson:"officerDesignation"   json:"officerDesignation"`
	OfficerMobileNo    string             `bson:"officerMobileNo"      json:"officerMobileNo"`
	OfficerEmail       string             `bson:"officerEmail"         json:"officerEmail"`
	SSOId              string             `bson:"ssoId"                json:"ssoId"`
	IsDeleted          bool               `bson:"isDeleted"            json:"isDeleted"`
	DeletedAt          *time.Time         `bson:"deletedAt,omitempty"  json:"deletedAt,omitempty"`
	CreatedAt          time.Time          `bson:"createdAt"            json:"createdAt"`
	UpdatedAt          time.Time          `bson:"updatedAt"            json:"updatedAt"`
}

// ---------------------------------------------------------------------------
// Sync Payload  (combined inbound payload from the external company)
// ---------------------------------------------------------------------------

// SyncPayload is the combined inbound request from the external company.
type SyncPayload struct {
	GSTHierarchy []GSTHierarchy `json:"gstHierarchy"`
	Premise      []Premise      `json:"premise"`
	Machine      []Machine      `json:"machine"`
	Officer      []Officer      `json:"officer"`
}

// SyncResult summarises what happened for each module.
type SyncResult struct {
	GSTHierarchyResult ModuleResult `json:"gstHierarchy"`
	PremiseResult      ModuleResult `json:"premise"`
	MachineResult      ModuleResult `json:"machine"`
	OfficerResult      ModuleResult `json:"officer"`
}

// ModuleResult gives per-module upsert counts.
type ModuleResult struct {
	Received int `json:"received"`
	Upserted int `json:"upserted"`
	Modified int `json:"modified"`
}

// ---------------------------------------------------------------------------
// Pagination & Query helpers
// ---------------------------------------------------------------------------

// ListQuery holds common query parameters for all GET list endpoints.
type ListQuery struct {
	Page      int64  `form:"page"`
	Limit     int64  `form:"limit"`
	Search    string `form:"search"`    // generic text search
	SortBy    string `form:"sortBy"`    // field name to sort by
	SortOrder string `form:"sortOrder"` // "asc" or "desc"
	Deleted   bool   `form:"deleted"`   // if true, include soft-deleted records
}

// Normalize applies defaults and caps to a ListQuery.
func (q *ListQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 20
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	if q.SortOrder != "asc" && q.SortOrder != "desc" {
		q.SortOrder = "desc"
	}
	if q.SortBy == "" {
		q.SortBy = "createdAt"
	}
}

// Skip returns the MongoDB skip value.
func (q *ListQuery) Skip() int64 { return (q.Page - 1) * q.Limit }

// PaginatedResponse wraps list results with pagination metadata.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int64       `json:"page"`
	Limit      int64       `json:"limit"`
	TotalPages int64       `json:"totalPages"`
}

// NewPaginatedResponse builds a PaginatedResponse.
func NewPaginatedResponse(data interface{}, total, page, limit int64) PaginatedResponse {
	pages := total / limit
	if total%limit != 0 {
		pages++
	}
	return PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: pages,
	}
}

// ---------------------------------------------------------------------------
// Update Request DTOs  (used by PUT /api/v1/<module>/:id)
// All fields are optional — only non-zero fields will be applied.
// ---------------------------------------------------------------------------

type UpdateGSTHierarchyRequest struct {
	GSTZoneName            string `json:"gstZoneName"`
	GSTZoneCode            string `json:"gstZoneCode"`
	GSTCommissionerateName string `json:"gstCommissionerateName"`
	GSTCommissionerateCode string `json:"gstCommissionerateCode"`
	GSTDivisionCode        string `json:"gstDivisionCode"`
	GSTDivisionName        string `json:"gstDivisionName"`
	GSTRangeName           string `json:"gstrangeName"`
	GSTRangeCode           string `json:"gstrangeCode"`
	RangeStateName         string `json:"rangestateName"`
	RangePincode           string `json:"rangepincode"`
}

type UpdatePremiseRequest struct {
	GSTHierarchyId   string `json:"gstHierarchyId"`
	ManufacturerName string `json:"manufacturerName"`
	PremiseName      string `json:"premiseName"`
	PremiseCode      string `json:"premiseCode"`
	GSTIN            string `json:"gstin"`
	PremiseAddress   string `json:"premiseAddress"`
	PremiseCity      string `json:"premiseCity"`
	PremisePincode   string `json:"premisePincode"`
	PremiseLatLong   string `json:"premiseLatLong"`
}

type UpdateMachineRequest struct {
	PremiseId             string `json:"premiseId"`
	GSTHierarchyId        string `json:"gstHierarchyId"`
	MachineName           string `json:"machineName"`
	MachineRegistrationNo string `json:"machineRegistrationNo"`
	MachineType           string `json:"machineType"`
	MachineMake           string `json:"machineMake"`
	MachineModel          string `json:"machineModel"`
	MachineSerialNo       string `json:"machineSerialNo"`
	WorkingStatus         string `json:"workingStatus"`
}

type UpdateOfficerRequest struct {
	GSTHierarchyId     string `json:"gstHierarchyId"`
	OfficerName        string `json:"officerName"`
	OfficerCode        string `json:"officerCode"`
	OfficerDesignation string `json:"officerDesignation"`
	OfficerMobileNo    string `json:"officerMobileNo"`
	OfficerEmail       string `json:"officerEmail"`
	SSOId              string `json:"ssoId"`
}
