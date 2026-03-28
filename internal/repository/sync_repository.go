package repository

import (
	"context"
	"time"

	"gst-api/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection names – change here if needed.
const (
	ColGSTHierarchy = "gst_hierarchies"
	ColPremise      = "premises"
	ColMachine      = "machines"
	ColOfficer      = "officers"
)

// SyncRepository exposes all DB operations for the four modules.
type SyncRepository struct {
	db *mongo.Database
}

// New creates a SyncRepository and ensures indexes exist on startup.
func New(db *mongo.Database) (*SyncRepository, error) {
	r := &SyncRepository{db: db}
	// if err := r.ensureIndexes(context.Background()); err != nil {
	// 	return nil, err
	// }
	return r, nil
}

// ---------------------------------------------------------------------------
// GSTHierarchy
// ---------------------------------------------------------------------------

// UpsertGSTHierarchies is kept for backward-compat (not used by new sync endpoints).
func (r *SyncRepository) UpsertGSTHierarchies(ctx context.Context, records []models.GSTHierarchy) (models.ModuleResult, error) {
	result := models.ModuleResult{Received: len(records)}
	if len(records) == 0 {
		return result, nil
	}
	col := r.db.Collection(ColGSTHierarchy)
	now := time.Now().UTC()

	for _, rec := range records {
		rec.UpdatedAt = now
		filter := bson.M{"gstrangeCode": rec.GSTRangeCode}
		update := bson.M{
			"$set":         rec,
			"$setOnInsert": bson.M{"createdAt": now},
		}
		res, err := col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return result, err
		}
		result.Upserted += int(res.UpsertedCount)
		result.Modified += int(res.ModifiedCount)
	}
	return result, nil
}

// CreateGSTHierarchy inserts a new record. Returns duplicate key error if
// gstHierarchyId or gstrangeCode already exists.
func (r *SyncRepository) CreateGSTHierarchy(ctx context.Context, rec *models.GSTHierarchy) error {
	now := time.Now().UTC()
	rec.CreatedAt = now
	rec.UpdatedAt = now
	rec.IsDeleted = false
	_, err := r.db.Collection(ColGSTHierarchy).InsertOne(ctx, rec)
	return err
}

// UpdateGSTHierarchyByBusinessId updates fields on the record matching gstHierarchyId.
// Returns (false, nil) when no record matched.
func (r *SyncRepository) UpdateGSTHierarchyByBusinessId(ctx context.Context, gstHierarchyId string, req models.UpdateGSTHierarchyRequest) (bool, error) {
	fields := buildUpdateFields(map[string]string{
		"gstZoneName":            req.GSTZoneName,
		"gstZoneCode":            req.GSTZoneCode,
		"gstCommissionerateName": req.GSTCommissionerateName,
		"gstCommissionerateCode": req.GSTCommissionerateCode,
		"gstDivisionCode":        req.GSTDivisionCode,
		"gstDivisionName":        req.GSTDivisionName,
		"gstrangeName":           req.GSTRangeName,
		"gstrangeCode":           req.GSTRangeCode,
		"rangestateName":         req.RangeStateName,
		"rangepincode":           req.RangePincode,
	})
	if len(fields) == 0 {
		return false, ErrNoFieldsToUpdate
	}
	fields["updatedAt"] = time.Now().UTC()
	res, err := r.db.Collection(ColGSTHierarchy).UpdateOne(ctx,
		bson.M{"gstHierarchyId": gstHierarchyId, "isDeleted": bson.M{"$ne": true}},
		bson.M{"$set": fields},
	)
	if err != nil {
		return false, err
	}
	return res.MatchedCount > 0, nil
}

// ---------------------------------------------------------------------------
// Premise
// ---------------------------------------------------------------------------

func (r *SyncRepository) UpsertPremises(ctx context.Context, records []models.Premise) (models.ModuleResult, error) {
	result := models.ModuleResult{Received: len(records)}
	if len(records) == 0 {
		return result, nil
	}
	col := r.db.Collection(ColPremise)
	now := time.Now().UTC()

	for _, rec := range records {
		rec.UpdatedAt = now
		filter := bson.M{"premiseCode": rec.PremiseCode}
		update := bson.M{
			"$set":         rec,
			"$setOnInsert": bson.M{"createdAt": now},
		}
		res, err := col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return result, err
		}
		result.Upserted += int(res.UpsertedCount)
		result.Modified += int(res.ModifiedCount)
	}
	return result, nil
}

func (r *SyncRepository) CreatePremise(ctx context.Context, rec *models.Premise) error {
	now := time.Now().UTC()
	rec.CreatedAt = now
	rec.UpdatedAt = now
	rec.IsDeleted = false
	_, err := r.db.Collection(ColPremise).InsertOne(ctx, rec)
	return err
}

func (r *SyncRepository) UpdatePremiseByBusinessId(ctx context.Context, premiseId string, req models.UpdatePremiseRequest) (bool, error) {
	fields := buildUpdateFields(map[string]string{
		"manufacturerName": req.ManufacturerName,
		"premiseName":      req.PremiseName,
		"premiseCode":      req.PremiseCode,
		"gstin":            req.GSTIN,
		"premiseAddress":   req.PremiseAddress,
		"premiseCity":      req.PremiseCity,
		"premisePincode":   req.PremisePincode,
		"premiseLatLong":   req.PremiseLatLong,
	})
	if len(fields) == 0 {
		return false, ErrNoFieldsToUpdate
	}
	fields["updatedAt"] = time.Now().UTC()
	res, err := r.db.Collection(ColPremise).UpdateOne(ctx,
		bson.M{"premiseId": premiseId, "isDeleted": bson.M{"$ne": true}},
		bson.M{"$set": fields},
	)
	if err != nil {
		return false, err
	}
	return res.MatchedCount > 0, nil
}

// ---------------------------------------------------------------------------
// Machine
// ---------------------------------------------------------------------------

func (r *SyncRepository) UpsertMachines(ctx context.Context, records []models.Machine) (models.ModuleResult, error) {
	result := models.ModuleResult{Received: len(records)}
	if len(records) == 0 {
		return result, nil
	}
	col := r.db.Collection(ColMachine)
	now := time.Now().UTC()

	for _, rec := range records {
		rec.UpdatedAt = now
		filter := bson.M{"machineRegistrationNo": rec.MachineRegistrationNo}
		update := bson.M{
			"$set":         rec,
			"$setOnInsert": bson.M{"createdAt": now},
		}
		res, err := col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return result, err
		}
		result.Upserted += int(res.UpsertedCount)
		result.Modified += int(res.ModifiedCount)
	}
	return result, nil
}

func (r *SyncRepository) CreateMachine(ctx context.Context, rec *models.Machine) error {
	now := time.Now().UTC()
	rec.CreatedAt = now
	rec.UpdatedAt = now
	rec.IsDeleted = false
	_, err := r.db.Collection(ColMachine).InsertOne(ctx, rec)
	return err
}

func (r *SyncRepository) UpdateMachineByBusinessId(ctx context.Context, machineId string, req models.UpdateMachineRequest) (bool, error) {
	fields := buildUpdateFields(map[string]string{
		"machineName":           req.MachineName,
		"machineRegistrationNo": req.MachineRegistrationNo,
		"machineType":           req.MachineType,
		"machineMake":           req.MachineMake,
		"machineModel":          req.MachineModel,
		"machineSerialNo":       req.MachineSerialNo,
		"workingStatus":         req.WorkingStatus,
	})
	if len(fields) == 0 {
		return false, ErrNoFieldsToUpdate
	}
	fields["updatedAt"] = time.Now().UTC()
	res, err := r.db.Collection(ColMachine).UpdateOne(ctx,
		bson.M{"machineId": machineId, "isDeleted": bson.M{"$ne": true}},
		bson.M{"$set": fields},
	)
	if err != nil {
		return false, err
	}
	return res.MatchedCount > 0, nil
}

// ---------------------------------------------------------------------------
// Officer
// ---------------------------------------------------------------------------

func (r *SyncRepository) UpsertOfficers(ctx context.Context, records []models.Officer) (models.ModuleResult, error) {
	result := models.ModuleResult{Received: len(records)}
	if len(records) == 0 {
		return result, nil
	}
	col := r.db.Collection(ColOfficer)
	now := time.Now().UTC()

	for _, rec := range records {
		rec.UpdatedAt = now
		filter := bson.M{"officerCode": rec.OfficerCode}
		update := bson.M{
			"$set":         rec,
			"$setOnInsert": bson.M{"createdAt": now},
		}
		res, err := col.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			return result, err
		}
		result.Upserted += int(res.UpsertedCount)
		result.Modified += int(res.ModifiedCount)
	}
	return result, nil
}

func (r *SyncRepository) CreateOfficer(ctx context.Context, rec *models.Officer) error {
	now := time.Now().UTC()
	rec.CreatedAt = now
	rec.UpdatedAt = now
	rec.IsDeleted = false
	_, err := r.db.Collection(ColOfficer).InsertOne(ctx, rec)
	return err
}

func (r *SyncRepository) UpdateOfficerByBusinessId(ctx context.Context, officerId string, req models.UpdateOfficerRequest) (bool, error) {
	fields := buildUpdateFields(map[string]string{
		"officerName":        req.OfficerName,
		"officerCode":        req.OfficerCode,
		"officerDesignation": req.OfficerDesignation,
		"officerMobileNo":    req.OfficerMobileNo,
		"officerEmail":       req.OfficerEmail,
		"ssoId":              req.SSOId,
	})
	if len(fields) == 0 {
		return false, ErrNoFieldsToUpdate
	}
	fields["updatedAt"] = time.Now().UTC()
	res, err := r.db.Collection(ColOfficer).UpdateOne(ctx,
		bson.M{"officerId": officerId, "isDeleted": bson.M{"$ne": true}},
		bson.M{"$set": fields},
	)
	if err != nil {
		return false, err
	}
	return res.MatchedCount > 0, nil
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// ErrNoFieldsToUpdate is returned when an update request has no non-empty fields.
var ErrNoFieldsToUpdate = mongo.ErrNoDocuments

// buildUpdateFields filters out empty string values so we only $set sent fields.
func buildUpdateFields(input map[string]string) bson.M {
	out := bson.M{}
	for k, v := range input {
		if v != "" {
			out[k] = v
		}
	}
	return out
}
