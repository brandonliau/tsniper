package repository

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"Tsniper/internal/shared"
	"Tsniper/pkg/config"
	"Tsniper/pkg/database"
	"Tsniper/pkg/logger"

	"github.com/bwmarrin/discordgo"
)

type snipeRepo struct {
	mu             sync.RWMutex
	dCfg           *config.DiscordConfig
	sCfg           *config.ServiceConfig
	session        *discordgo.Session
	db             database.Database
	logger         logger.Logger
	registered     map[string]string
	trackedCount   map[string]map[string]int
	trackedIndices map[string][]string
}

func NewSnipeRepo(dCfg *config.DiscordConfig, sCfg *config.ServiceConfig, s *discordgo.Session, db database.Database, logger logger.Logger) *snipeRepo {
	registered := make(map[string]string)
	trackedCount := make(map[string]map[string]int)
	trackedIndices := make(map[string][]string)
	snipeRepo := &snipeRepo{
		dCfg:           dCfg,
		sCfg:           sCfg,
		session:        s,
		db:             db,
		logger:         logger,
		registered:     registered,
		trackedCount:   trackedCount,
		trackedIndices: trackedIndices,
	}
	return snipeRepo
}

// registered commands
func (repo *snipeRepo) Register(name, id string) {
	repo.registered[name] = id
}

func (repo *snipeRepo) Registered() map[string]string {
	return repo.registered
}

// in-memory repository
func (repo *snipeRepo) TrackedIndices(campus, season string) []string {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	return repo.trackedIndices[campus+season]
}

func (repo *snipeRepo) Sync() {
	// sync users
	start := time.Now()
	rows, _ := repo.db.Query("SELECT DISTINCT user_id FROM snipes")
	defer rows.Close()
	dbUserList := make([]string, 0)
	var userID string
	for rows.Next() {
		rows.Scan(&userID)
		dbUserList = append(dbUserList, userID)
	}
	var lastUserID string
	currentGuildMembers := make([]string, 0)
	guild, _ := repo.session.State.Guild(repo.dCfg.Guild)
	for {
		gm, _ := repo.session.GuildMembers(guild.ID, lastUserID, 1000)
		if len(gm) < 1 {
			break
		}
		for _, member := range gm {
			currentGuildMembers = append(currentGuildMembers, member.User.ID)
		}
		lastUserID = gm[len(gm)-1].User.ID
	}

	// clear snipes and remove user from db if they are in dbUserList but not in currentGuildMembers
	diff := shared.Difference(dbUserList, currentGuildMembers)
	for _, userID := range diff {
		repo.ClearSnipe(userID)
	}
	repo.logger.Info("Synced %d users in %v", len(currentGuildMembers), time.Since(start))

	// clean off-season snipes
	start = time.Now()
	placeholders := strings.Repeat("?,", len(repo.sCfg.Seasons))
	placeholders = placeholders[:len(placeholders)-1]
	query := fmt.Sprintf("DELETE FROM snipes WHERE season NOT IN (%s)", placeholders)
	params := make([]interface{}, len(repo.sCfg.Seasons))
	for i, season := range repo.sCfg.Seasons {
		params[i] = season
	}
	rowsAffected, _ := repo.db.ExecWithResult(query, params...)
	repo.logger.Info("Cleared %d off season in %v", rowsAffected, time.Since(start))

	// sync tracking
	start = time.Now()
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for _, campus := range repo.sCfg.Campuses {
		for _, season := range repo.sCfg.Seasons {
			repo.trackedCount[campus+season] = make(map[string]int)
			repo.trackedIndices[campus+season] = make([]string, 0)
		}
	}
	rows, _ = repo.db.Query("SELECT course_index, campus, season, COUNT(*) AS count FROM snipes GROUP BY course_index, campus, season")
	defer rows.Close()
	numSnipes := 0
	var index, campus, season string
	var count int
	for rows.Next() {
		rows.Scan(&index, &campus, &season, &count)
		repo.trackedCount[campus+season][index] = count
		repo.trackedIndices[campus+season] = append(repo.trackedIndices[campus+season], index)
		numSnipes += 1
	}
	repo.logger.Info("Synced %d snipes in %v", numSnipes, time.Since(start))
	repo.logger.Debug("%v", repo.trackedCount)
	repo.logger.Debug("%v", repo.trackedIndices)
}

func (repo *snipeRepo) Add(index, campus, season string) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.trackedCount[campus+season][index] += 1
	if repo.trackedCount[campus+season][index] == 1 {
		repo.trackedIndices[campus+season] = append(repo.trackedIndices[campus+season], index)
	}
	repo.logger.Debug("%v", repo.trackedCount)
	repo.logger.Debug("%v", repo.trackedIndices)
}

func (repo *snipeRepo) Remove(index, campus, season string) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.trackedCount[campus+season][index] -= 1
	if repo.trackedCount[campus+season][index] == 0 {
		delete(repo.trackedCount[campus+season], index)
		trackedIndices := repo.trackedIndices[campus+season]
		for i, v := range trackedIndices {
			if v == index {
				trackedIndices[i] = trackedIndices[len(trackedIndices)-1]
				repo.trackedIndices[campus+season] = trackedIndices[:len(trackedIndices)-1]
				break
			}
		}
	}
	repo.logger.Debug("%v", repo.trackedCount)
	repo.logger.Debug("%v", repo.trackedIndices)
}

// db user management
func (repo *snipeRepo) Snipes(userID string) [][]string {
	snipes := make([][]string, 0)
	rows, _ := repo.db.Query("SELECT course_index, campus, season FROM snipes WHERE user_id = ?", userID)
	defer rows.Close()
	var index, campus, season string
	for rows.Next() {
		rows.Scan(&index, &campus, &season)
		snipes = append(snipes, []string{index, campus, season})
	}
	sort.SliceStable(snipes, func(i, j int) bool {
		return snipes[i][0] < snipes[j][0]
	})
	return snipes
}

func (repo *snipeRepo) IsSniping(userID, index, campus, season string) bool {
	row, _ := repo.db.QueryRow("SELECT 1 FROM snipes WHERE course_index = ? AND season = ? AND campus = ? AND user_id = ?", index, season, campus, userID)
	var exist int
	err := row.Scan(&exist)
	return err == nil // return true if course exists in db
}

// db snipe management
func (repo *snipeRepo) AddSnipe(userID, index, campus, season string) {
	repo.db.Exec("INSERT INTO snipes (user_id, course_index, campus, season) VALUES (?, ?, ?, ?)", userID, index, campus, season)
}

func (repo *snipeRepo) RemoveSnipe(userID, index, campus, season string) {
	repo.db.Exec("DELETE FROM snipes WHERE course_index = ? AND season = ? AND campus = ? AND user_id = ?", index, season, campus, userID)
}

func (repo *snipeRepo) ClearSnipe(userID string) {
	repo.db.Exec("DELETE FROM snipes WHERE user_id = ?", userID)
}

// db course management
func (repo *snipeRepo) CourseEntry(index, campus, season string) shared.CourseEntry {
	query := `SELECT course_index, title, course_string, section, instructors, notes, meeting
		FROM courses WHERE course_index = ? AND campus = ? AND season = ?`

	row, _ := repo.db.QueryRow(query, index, campus, season)
	var course_index, title, courseString, section, instructors, notes, meeting string
	row.Scan(&course_index, &title, &courseString, &section, &instructors, &notes, &meeting)
	return shared.CourseEntry{
		Title:        title,
		CourseString: courseString,
		Index:        course_index,
		Section:      section,
		Instructors:  instructors,
		Notes:        notes,
		Meeting:      meeting,
	}
}

func (repo *snipeRepo) Exists(index, campus, season string) bool {
	query := "SELECT 1 FROM courses WHERE course_index = ? AND campus = ? AND season = ?"
	row, _ := repo.db.QueryRow(query, index, campus, season)
	var exist int
	err := row.Scan(&exist)
	return err == nil // return true if course exists in db
}

func (repo *snipeRepo) SnipeUsers(index, campus, season string) []string {
	rows, _ := repo.db.Query("SELECT user_id FROM snipes WHERE course_index = ? AND season = ?", index, season)
	defer rows.Close()
	var users []string
	var userID string
	for rows.Next() {
		rows.Scan(&userID)
		users = append(users, userID)
	}
	return users
}

func (repo *snipeRepo) SnipeCount(index, campus, season string) int {
	row, _ := repo.db.QueryRow("SELECT count(*) FROM snipes WHERE course_index = ? AND season = ? AND campus = ?", index, season, campus)
	var count int
	row.Scan(&count)
	return count
}

func (repo *snipeRepo) LastOpen(index, campus, season string) int {
	query := "SELECT last_open FROM courses WHERE course_index = ? AND campus = ? AND season = ?"
	row, _ := repo.db.QueryRow(query, index, campus, season)
	var lastOpen int
	row.Scan(&lastOpen)
	return lastOpen
}

func (repo *snipeRepo) LastOpens(campus, season string) map[string]int64 {
	lastOpens := make(map[string]int64)
	query := "SELECT course_index, last_open FROM courses WHERE campus = ? AND season = ?"
	rows, _ := repo.db.Query(query, campus, season)
	defer rows.Close()
	var index string
	var lastOpen int64
	for rows.Next() {
		rows.Scan(&index, &lastOpen)
		lastOpens[index] = lastOpen
	}
	return lastOpens
}

func (repo *snipeRepo) UpdateLastOpen(index, campus, season string, lastOpen int64) {
	query := "UPDATE courses SET last_open = ? WHERE course_index = ? AND campus = ? AND season = ?"
	repo.db.Exec(query, lastOpen, index, campus, season)
}

// discord user management
func (repo *snipeRepo) Campus(userID string) string {
	var member *discordgo.Member
	var err error
	member, err = repo.session.State.Member(repo.dCfg.Guild, userID)
	if err != nil {
		member, err = repo.session.GuildMember(repo.dCfg.Guild, userID)
		if err != nil {
			repo.logger.Error("Failed to get guild member %s: %v", userID, err)
			return repo.sCfg.DefaultCampus
		}
	}
	for _, roleID := range member.Roles {
		role, _ := repo.session.State.Role(repo.dCfg.Guild, roleID)
		switch role.Name {
		case "New Brunswick":
			return "NB"
		case "Newark":
			return "NK"
		case "Camden":
			return "CM"
		}
	}
	return repo.sCfg.DefaultCampus
}
