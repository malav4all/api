package handler

import (
	"log"
	"net/http"

	"gst-api/internal/models"
	"gst-api/internal/repository"
	"gst-api/pkg/response"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SyncHandler handles inbound data pushes from the external company.
// Two endpoints are exposed:
//
//	POST  /api/v1/sync  → CreateSync  (insert new records; fails on duplicate)
//	PATCH /api/v1/sync  → UpdateSync  (update existing records by unique business key)
type SyncHandler struct {
	repo *repository.SyncRepository
}

func NewSyncHandler(repo *repository.SyncRepository) *SyncHandler {
	return &SyncHandler{repo: repo}
}

// ─────────────────────────────────────────────────────────────────────────────
// POST /api/v1/sync  — external company pushes NEW records
// ─────────────────────────────────────────────────────────────────────────────

func (h *SyncHandler) CreateSync(c *gin.Context) {
	payload, ok := bindSyncPayload(c)
	if !ok {
		return
	}

	ctx := c.Request.Context()
	var result models.SyncResult

	// ── GSTHierarchy ──
	result.GSTHierarchyResult.Received = len(payload.GSTHierarchy)
	for i := range payload.GSTHierarchy {
		if err := h.repo.CreateGSTHierarchy(ctx, &payload.GSTHierarchy[i]); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				c.JSON(http.StatusConflict, gin.H{
					"success": false,
					"message": "gstHierarchy: duplicate gstHierarchyId or gstrangeCode",
					"record":  payload.GSTHierarchy[i].GSTRangeCode,
				})
				return
			}
			log.Printf("[SYNC CREATE] gstHierarchy error: %v", err)
			response.InternalServerError(c, "failed to create gstHierarchy: "+err.Error())
			return
		}
		result.GSTHierarchyResult.Upserted++
	}

	// ── Premise ──
	result.PremiseResult.Received = len(payload.Premise)
	for i := range payload.Premise {
		if err := h.repo.CreatePremise(ctx, &payload.Premise[i]); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				c.JSON(http.StatusConflict, gin.H{
					"success": false,
					"message": "premise: duplicate premiseId or premiseCode",
					"record":  payload.Premise[i].PremiseCode,
				})
				return
			}
			log.Printf("[SYNC CREATE] premise error: %v", err)
			response.InternalServerError(c, "failed to create premise: "+err.Error())
			return
		}
		result.PremiseResult.Upserted++
	}

	// ── Machine ──
	result.MachineResult.Received = len(payload.Machine)
	for i := range payload.Machine {
		if err := h.repo.CreateMachine(ctx, &payload.Machine[i]); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				c.JSON(http.StatusConflict, gin.H{
					"success": false,
					"message": "machine: duplicate machineId or machineRegistrationNo",
					"record":  payload.Machine[i].MachineRegistrationNo,
				})
				return
			}
			log.Printf("[SYNC CREATE] machine error: %v", err)
			response.InternalServerError(c, "failed to create machine: "+err.Error())
			return
		}
		result.MachineResult.Upserted++
	}

	// ── Officer ──
	result.OfficerResult.Received = len(payload.Officer)
	for i := range payload.Officer {
		if err := h.repo.CreateOfficer(ctx, &payload.Officer[i]); err != nil {
			if mongo.IsDuplicateKeyError(err) {
				c.JSON(http.StatusConflict, gin.H{
					"success": false,
					"message": "officer: duplicate officerId or officerCode",
					"record":  payload.Officer[i].OfficerCode,
				})
				return
			}
			log.Printf("[SYNC CREATE] officer error: %v", err)
			response.InternalServerError(c, "failed to create officer: "+err.Error())
			return
		}
		result.OfficerResult.Upserted++
	}

	log.Printf("[SYNC CREATE] gstHierarchy:%+v premise:%+v machine:%+v officer:%+v",
		result.GSTHierarchyResult, result.PremiseResult, result.MachineResult, result.OfficerResult)

	response.Created(c, "data created successfully", result)
}

// ─────────────────────────────────────────────────────────────────────────────
// PATCH /api/v1/sync  — external company pushes UPDATED records
// ─────────────────────────────────────────────────────────────────────────────

func (h *SyncHandler) UpdateSync(c *gin.Context) {
	payload, ok := bindSyncPayload(c)
	if !ok {
		return
	}

	ctx := c.Request.Context()
	var result models.SyncResult
	var notFound []string

	// ── GSTHierarchy ──
	result.GSTHierarchyResult.Received = len(payload.GSTHierarchy)
	for _, rec := range payload.GSTHierarchy {
		req := models.UpdateGSTHierarchyRequest{
			GSTZoneName:            rec.GSTZoneName,
			GSTZoneCode:            rec.GSTZoneCode,
			GSTCommissionerateName: rec.GSTCommissionerateName,
			GSTCommissionerateCode: rec.GSTCommissionerateCode,
			GSTDivisionCode:        rec.GSTDivisionCode,
			GSTDivisionName:        rec.GSTDivisionName,
			GSTRangeName:           rec.GSTRangeName,
			GSTRangeCode:           rec.GSTRangeCode,
			RangeStateName:         rec.RangeStateName,
			RangePincode:           rec.RangePincode,
		}
		found, err := h.repo.UpdateGSTHierarchyByBusinessId(ctx, rec.GSTHierarchyId, req)
		if err != nil {
			log.Printf("[SYNC UPDATE] gstHierarchy error: %v", err)
			response.InternalServerError(c, "failed to update gstHierarchy: "+err.Error())
			return
		}
		if found {
			result.GSTHierarchyResult.Modified++
		} else {
			notFound = append(notFound, "gstHierarchy:"+rec.GSTHierarchyId)
		}
	}

	// ── Premise ──
	result.PremiseResult.Received = len(payload.Premise)
	for _, rec := range payload.Premise {
		req := models.UpdatePremiseRequest{
			ManufacturerName: rec.ManufacturerName,
			PremiseName:      rec.PremiseName,
			PremiseCode:      rec.PremiseCode,
			GSTIN:            rec.GSTIN,
			PremiseAddress:   rec.PremiseAddress,
			PremiseCity:      rec.PremiseCity,
			PremisePincode:   rec.PremisePincode,
			PremiseLatLong:   rec.PremiseLatLong,
		}
		found, err := h.repo.UpdatePremiseByBusinessId(ctx, rec.PremiseId, req)
		if err != nil {
			log.Printf("[SYNC UPDATE] premise error: %v", err)
			response.InternalServerError(c, "failed to update premise: "+err.Error())
			return
		}
		if found {
			result.PremiseResult.Modified++
		} else {
			notFound = append(notFound, "premise:"+rec.PremiseId)
		}
	}

	// ── Machine ──
	result.MachineResult.Received = len(payload.Machine)
	for _, rec := range payload.Machine {
		req := models.UpdateMachineRequest{
			MachineName:           rec.MachineName,
			MachineRegistrationNo: rec.MachineRegistrationNo,
			MachineType:           rec.MachineType,
			MachineMake:           rec.MachineMake,
			MachineModel:          rec.MachineModel,
			MachineSerialNo:       rec.MachineSerialNo,
			WorkingStatus:         rec.WorkingStatus,
		}
		found, err := h.repo.UpdateMachineByBusinessId(ctx, rec.MachineId, req)
		if err != nil {
			log.Printf("[SYNC UPDATE] machine error: %v", err)
			response.InternalServerError(c, "failed to update machine: "+err.Error())
			return
		}
		if found {
			result.MachineResult.Modified++
		} else {
			notFound = append(notFound, "machine:"+rec.MachineId)
		}
	}

	// ── Officer ──
	result.OfficerResult.Received = len(payload.Officer)
	for _, rec := range payload.Officer {
		req := models.UpdateOfficerRequest{
			OfficerName:        rec.OfficerName,
			OfficerCode:        rec.OfficerCode,
			OfficerDesignation: rec.OfficerDesignation,
			OfficerMobileNo:    rec.OfficerMobileNo,
			OfficerEmail:       rec.OfficerEmail,
			SSOId:              rec.SSOId,
		}
		found, err := h.repo.UpdateOfficerByBusinessId(ctx, rec.OfficerId, req)
		if err != nil {
			log.Printf("[SYNC UPDATE] officer error: %v", err)
			response.InternalServerError(c, "failed to update officer: "+err.Error())
			return
		}
		if found {
			result.OfficerResult.Modified++
		} else {
			notFound = append(notFound, "officer:"+rec.OfficerId)
		}
	}

	log.Printf("[SYNC UPDATE] gstHierarchy:%+v premise:%+v machine:%+v officer:%+v notFound:%v",
		result.GSTHierarchyResult, result.PremiseResult, result.MachineResult, result.OfficerResult, notFound)

	resp := gin.H{"result": result}
	if len(notFound) > 0 {
		resp["notFound"] = notFound
	}
	response.Success(c, "data updated successfully", resp)
}

// ─────────────────────────────────────────────────────────────────────────────
// shared helper
// ─────────────────────────────────────────────────────────────────────────────

func bindSyncPayload(c *gin.Context) (models.SyncPayload, bool) {
	var payload models.SyncPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.BadRequest(c, "invalid JSON payload: "+err.Error())
		return payload, false
	}
	if len(payload.GSTHierarchy) == 0 &&
		len(payload.Premise) == 0 &&
		len(payload.Machine) == 0 &&
		len(payload.Officer) == 0 {
		response.BadRequest(c, "payload must contain at least one module array with records")
		return payload, false
	}
	return payload, true
}
