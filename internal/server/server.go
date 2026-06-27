package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SorynMochi/Namagotchi/internal/database"
)

type Server struct {
	Store     *database.Store
	StartedAt time.Time
}

type StatusResponse struct {
	Server      string `json:"server"`
	Database    string `json:"database"`
	Timestamp   string `json:"timestamp"`
	Uptime      string `json:"uptime"`
	Version     string `json:"version"`
	OnlineUsers int    `json:"onlineUsers"`
}

type MessageResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

type PlayerOnlineTickResponse struct {
	OK            bool  `json:"ok"`
	OnlineSeconds int64 `json:"onlineSeconds"`
}

type DevWardrobeSpawnResponse struct {
	OK      bool                        `json:"ok"`
	Message string                      `json:"message"`
	ItemID  int64                       `json:"itemId"`
	Detail  database.WardrobeItemDetail `json:"detail"`
}

type WardrobeItemActionRequest struct {
	ItemID  int64  `json:"itemId"`
	SlotKey string `json:"slotKey"`
}

type WardrobeItemActionResponse struct {
	OK      bool                        `json:"ok"`
	Message string                      `json:"message"`
	Detail  database.WardrobeItemDetail `json:"detail"`
}

type GatheringRequest struct {
	Task string `json:"task"`
}

type CareActionRequest struct {
	Action string `json:"action"`
}

func New(store *database.Store, startedAt time.Time) *Server {
	return &Server{
		Store:     store,
		StartedAt: startedAt,
	}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/status", s.HandleStatus)
	mux.HandleFunc("/health", s.HandleStatus)

	mux.HandleFunc("/api/auth/register", s.requireRateLimit("auth_register", 8, 10*time.Minute, s.HandleAuthRegister))
	mux.HandleFunc("/api/auth/login", s.requireRateLimit("auth_login", 10, 10*time.Minute, s.HandleAuthLogin))
	mux.HandleFunc("/api/auth/logout", s.requireCSRF(s.HandleAuthLogout))
	mux.HandleFunc("/api/auth/me", s.HandleAuthMe)
	mux.HandleFunc("/api/auth/display-name", s.requireAuth(s.requireCSRF(s.HandleAuthDisplayName)))
	mux.HandleFunc("/api/auth/csrf", s.HandleCSRFToken)
	mux.HandleFunc("/api/auth/google/start", s.HandleAuthGoogleStart)
	mux.HandleFunc("/api/auth/google/callback", s.HandleAuthGoogleCallback)

	mux.HandleFunc("/api/dev/unlock", s.requireDev(s.requireCSRF(s.requireRateLimit("dev_unlock", 6, 10*time.Minute, s.HandleDevUnlock))))
	mux.HandleFunc("/api/dev/lock", s.requireDev(s.requireCSRF(s.HandleDevLock)))
	mux.HandleFunc("/api/dev/seed-player", s.requireDev(s.requireDevUnlock(s.withDevAudit("seed-player", s.requireCSRF(s.HandleSeedDevPlayer)))))
	mux.HandleFunc("/api/dev/force-tick", s.requireDev(s.requireDevUnlock(s.withDevAudit("force-tick", s.requireCSRF(s.HandleForceTick)))))
	mux.HandleFunc("/api/dev/reset-playdeck-streak", s.requireDev(s.requireDevUnlock(s.withDevAudit("reset-playdeck-streak", s.requireCSRF(s.HandleResetPlaydeckStreak)))))
	mux.HandleFunc("/api/dev/rewind-care-decay", s.requireDev(s.requireDevUnlock(s.withDevAudit("rewind-care-decay", s.requireCSRF(s.HandleRewindCareDecay)))))
	mux.HandleFunc("/api/dev/finish-care", s.requireDev(s.requireDevUnlock(s.withDevAudit("finish-care", s.requireCSRF(s.HandleFinishCareAction)))))
	mux.HandleFunc("/api/dev/spawn-wardrobe-item", s.requireDev(s.requireDevUnlock(s.withDevAudit("spawn-wardrobe-item", s.requireCSRF(s.HandleSpawnDevWardrobeItem)))))
	mux.HandleFunc("/api/dev/audit-logs", s.requireDev(s.requireDevUnlock(s.HandleDevAuditLogs)))
	mux.HandleFunc("/api/dev/security-events", s.requireDev(s.requireDevUnlock(s.HandleDevSecurityEvents)))
	mux.HandleFunc("/dev", s.requireDev(s.HandleDevConsolePage))
	mux.HandleFunc("/dev/", s.requireDev(s.HandleDevConsolePage))

	mux.HandleFunc("/api/player/status", s.requireAuth(s.HandlePlayerStatus))
	mux.HandleFunc("/api/player/sync", s.requireAuth(s.requireCSRF(s.HandlePlayerSync)))
	mux.HandleFunc("/api/player/online-tick", s.requireAuth(s.requireCSRF(s.HandlePlayerOnlineTick)))
	mux.HandleFunc("/api/player/wardrobe/item", s.requireAuth(s.HandleWardrobeItemDetail))
	mux.HandleFunc("/api/player/wardrobe/equip", s.requireAuth(s.requireCSRF(s.HandleEquipWardrobeItem)))
	mux.HandleFunc("/api/player/wardrobe/unequip", s.requireAuth(s.requireCSRF(s.HandleUnequipWardrobeItem)))
	mux.HandleFunc("/api/player/settle-ticks", s.requireAuth(s.requireCSRF(s.HandleSettleTicks)))
	mux.HandleFunc("/api/player/gathering", s.requireAuth(s.requireCSRF(s.HandleGatheringTask)))
	mux.HandleFunc("/api/player/care", s.requireAuth(s.requireCSRF(s.HandleCareAction)))
	mux.HandleFunc("/api/nami/messages", s.requireAuth(s.HandleNamiMessages))

	mux.Handle("/", http.FileServer(http.Dir("web")))

	return s.withSecurityHeaders(s.withRequestBodyLimit(mux))
}

func (s *Server) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	response := StatusResponse{
		Server:      "online",
		Database:    s.databaseStatus(),
		Timestamp:   time.Now().Format(time.RFC3339),
		Uptime:      time.Since(s.StartedAt).Round(time.Second).String(),
		Version:     "phase-3b",
		OnlineUsers: 1,
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *Server) HandleSeedDevPlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.Store.SeedDevPlayer(r.Context()); err != nil {
		log.Printf("seed dev player failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to seed dev player")
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{
		OK:      true,
		Message: "Dev player Soryn and Nami-chan are ready.",
	})
}

func (s *Server) HandleSpawnDevWardrobeItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	result, err := s.Store.SpawnDevRandomWardrobeItem(r.Context())
	if err != nil {
		log.Printf("spawn dev wardrobe item failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to spawn wardrobe item")
		return
	}

	writeJSON(w, http.StatusOK, DevWardrobeSpawnResponse{
		OK:      true,
		Message: "Random wardrobe item spawned.",
		ItemID:  result.ItemID,
		Detail:  result.Detail,
	})
}

func (s *Server) HandlePlayerStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	status, err := s.playerStatusForRequest(r)
	if err != nil {
		log.Printf("get player status failed: %v", err)
		writeError(w, http.StatusNotFound, "player status not found")
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (s *Server) HandlePlayerSync(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := s.syncPlayerState(r.Context()); err != nil {
		log.Printf("sync player state failed: %v", err)
		writeError(w, http.StatusInternalServerError, "player sync failed")
		return
	}

	status, err := s.playerStatusForRequest(r)
	if err != nil {
		log.Printf("get player status after sync failed: %v", err)
		writeError(w, http.StatusNotFound, "player status not found")
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (s *Server) HandlePlayerOnlineTick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	playerID, err := s.playerIDForRequest(r)
	if err != nil {
		log.Printf("get player for online tick failed: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	onlineSeconds, err := s.Store.TrackPlayerOnlineTime(r.Context(), playerID)
	if err != nil {
		log.Printf("track player online time failed: %v", err)
		writeError(w, http.StatusInternalServerError, "online tick failed")
		return
	}

	writeJSON(w, http.StatusOK, PlayerOnlineTickResponse{
		OK:            true,
		OnlineSeconds: onlineSeconds,
	})
}
func (s *Server) syncPlayerState(ctx context.Context) error {

	if _, err := s.Store.SettleDevTicks(ctx, 0); err != nil {
		return err
	}

	if err := s.Store.SettleDevCareActions(ctx); err != nil {
		return err
	}

	if err := s.Store.SettleDevCareDecay(ctx); err != nil {
		return err
	}

	if err := s.Store.GenerateDevPassiveNamiMessages(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Server) HandleWardrobeItemDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	rawItemID := strings.TrimSpace(r.URL.Query().Get("id"))
	itemID, err := strconv.ParseInt(rawItemID, 10, 64)
	if err != nil || itemID < 1 {
		writeError(w, http.StatusBadRequest, "id must be a positive whole number")
		return
	}

	playerID, err := s.playerIDForRequest(r)
	if err != nil {
		log.Printf("get player for wardrobe item failed: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	detail, err := s.Store.GetWardrobeItemDetail(
		r.Context(),
		playerID,
		itemID,
		r.URL.Query().Get("compareSlot"),
	)
	if err != nil {
		log.Printf("get wardrobe item detail failed: %v", err)
		writeError(w, http.StatusNotFound, "wardrobe item not found")
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

func (s *Server) HandleSettleTicks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	result, err := s.Store.SettleDevTicks(r.Context(), 0)
	if err != nil {
		log.Printf("settle ticks failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to settle ticks")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleForceTick(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ticks := int64(1)

	if rawTicks := r.URL.Query().Get("ticks"); rawTicks != "" {
		parsedTicks, err := strconv.ParseInt(rawTicks, 10, 64)
		if err != nil || parsedTicks < 1 {
			writeError(w, http.StatusBadRequest, "ticks must be a positive whole number")
			return
		}

		ticks = parsedTicks
	}

	result, err := s.Store.SettleDevTicks(r.Context(), ticks)
	if err != nil {
		log.Printf("force tick failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to force tick")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleResetPlaydeckStreak(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	playerID, err := s.Store.DevPlayerID(r.Context())
	if err != nil {
		log.Printf("get dev player for reset playdeck streak failed: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	commandTag, err := s.Store.Pool.Exec(r.Context(), `
        update player_tick_state
        set playdeck_streak = 0,
            last_tick_at = now(),
            updated_at = now()
        where player_id = $1
    `, playerID)
	if err != nil {
		log.Printf("reset playdeck streak failed: %v", err)
		writeError(w, http.StatusInternalServerError, "failed to reset playdeck streak")
		return
	}

	if commandTag.RowsAffected() == 0 {
		writeError(w, http.StatusNotFound, "playdeck streak state not found")
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{
		OK:      true,
		Message: "Playdeck streak reset.",
	})
}

func (s *Server) HandleRewindCareDecay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	rawHours := strings.TrimSpace(r.URL.Query().Get("hours"))
	if rawHours == "" {
		rawHours = "2"
	}

	hours, err := strconv.ParseFloat(rawHours, 64)
	if err != nil || hours <= 0 {
		writeError(w, http.StatusBadRequest, "hours must be a positive number")
		return
	}

	duration := time.Duration(hours * float64(time.Hour))

	if err := s.Store.RewindDevCareDecay(r.Context(), duration); err != nil {
		log.Printf("rewind care decay failed: %v", err)
		writeError(w, http.StatusInternalServerError, "rewind care decay failed")
		return
	}

	if err := s.Store.SettleDevCareDecay(r.Context()); err != nil {
		log.Printf("settle dev care decay failed: %v", err)
		writeError(w, http.StatusInternalServerError, "settle dev care decay failed")
		return
	}

	status, err := s.playerStatusForRequest(r)
	if err != nil {
		log.Printf("get player status after care decay rewind failed: %v", err)
		writeError(w, http.StatusInternalServerError, "player status failed")
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (s *Server) HandleGatheringTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var request GatheringRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid gathering request")
		return
	}

	request.Task = strings.TrimSpace(strings.ToLower(request.Task))
	accountID, err := accountIDForRequest(r)
	if err != nil {
		log.Printf("get account for gathering task failed: %v", err)
		writeError(w, http.StatusUnauthorized, "login required")
		return
	}

	if err := s.Store.SetGatheringTaskForAccount(r.Context(), accountID, request.Task); err != nil {
		log.Printf("set gathering task failed: %v", err)
		writeError(w, http.StatusBadRequest, "invalid gathering task")
		return
	}

	writeJSON(w, http.StatusOK, MessageResponse{
		OK:      true,
		Message: "Gathering task updated.",
	})
}

func (s *Server) HandleCareAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var request CareActionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid care action request")
		return
	}

	request.Action = strings.TrimSpace(strings.ToLower(request.Action))
	result, err := s.Store.StartOrQueueDevCareAction(r.Context(), request.Action)
	if err != nil {
		log.Printf("care action failed: %v", err)
		writeError(w, http.StatusBadRequest, "invalid care action")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (s *Server) HandleNamiMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	accountID, err := accountIDForRequest(r)
	if err != nil {
		log.Printf("get account for nami messages failed: %v", err)
		writeError(w, http.StatusUnauthorized, "login required")
		return
	}

	messages, err := s.Store.GetRecentNamiMessagesForAccount(r.Context(), accountID, 50)
	if err != nil {
		log.Printf("get nami messages failed: %v", err)
		writeError(w, http.StatusNotFound, "nami messages not found")
		return
	}

	writeJSON(w, http.StatusOK, messages)
}

func (s *Server) databaseStatus() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := s.Store.Pool.Ping(ctx); err != nil {
		log.Printf("database status check failed: %v", err)
		return "offline"
	}

	return "online"
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to write JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, MessageResponse{
		OK:      false,
		Message: message,
	})
}

func (s *Server) HandleEquipWardrobeItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request, err := decodeWardrobeItemActionRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wardrobe equip request")
		return
	}

	if request.ItemID < 1 {
		writeError(w, http.StatusBadRequest, "itemId must be a positive whole number")
		return
	}

	playerID, err := s.playerIDForRequest(r)
	if err != nil {
		log.Printf("get player for equip failed: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	detail, err := s.Store.EquipWardrobeItem(r.Context(), playerID, request.ItemID, request.SlotKey)
	if err != nil {
		log.Printf("equip wardrobe item failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, WardrobeItemActionResponse{
		OK:      true,
		Message: "Wardrobe item equipped.",
		Detail:  detail,
	})
}

func (s *Server) HandleUnequipWardrobeItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	request, err := decodeWardrobeItemActionRequest(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid wardrobe unequip request")
		return
	}

	if request.ItemID < 1 {
		writeError(w, http.StatusBadRequest, "itemId must be a positive whole number")
		return
	}

	playerID, err := s.playerIDForRequest(r)
	if err != nil {
		log.Printf("get player for unequip failed: %v", err)
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	detail, err := s.Store.UnequipWardrobeItem(r.Context(), playerID, request.ItemID, request.SlotKey)
	if err != nil {
		log.Printf("unequip wardrobe item failed: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, WardrobeItemActionResponse{
		OK:      true,
		Message: "Wardrobe item removed.",
		Detail:  detail,
	})
}

func decodeWardrobeItemActionRequest(r *http.Request) (WardrobeItemActionRequest, error) {
	var request WardrobeItemActionRequest

	if r.Body != nil && r.ContentLength != 0 {
		defer r.Body.Close()

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			return request, err
		}
	}

	if request.ItemID == 0 {
		itemID, _ := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("itemId")), 10, 64)
		request.ItemID = itemID
	}

	if request.SlotKey == "" {
		request.SlotKey = strings.TrimSpace(r.URL.Query().Get("slotKey"))
	}

	return request, nil
}
