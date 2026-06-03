package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
}

type PlayerStatus struct {
	Player    Player          `json:"player"`
	Companion CompanionState  `json:"companion"`
	Resources PlayerResources `json:"resources"`
}

type Player struct {
	ID            int64  `json:"id"`
	DisplayName   string `json:"displayName"`
	Level         int    `json:"level"`
	TotalXP       int64  `json:"totalXp"`
	CurrencyCents int64  `json:"currencyCents"`
}

type CompanionState struct {
	CompanionName     string    `json:"name"`
	MoodScore         float64   `json:"moodScore"`
	Satiety           int       `json:"satiety"`
	Connection        int       `json:"connection"`
	Energy            int       `json:"energy"`
	Comfort           int       `json:"comfort"`
	Playfulness       int       `json:"playfulness"`
	Inspiration       int       `json:"inspiration"`
	Cleanliness       int       `json:"cleanliness"`
	Status            string    `json:"status"`
	LastInteractionAt time.Time `json:"lastInteractionAt"`
}

type PlayerResources struct {
	Fans        int64 `json:"fans"`
	Memes       int64 `json:"memes"`
	LostItems   int64 `json:"lostItems"`
	Confidence  int64 `json:"confidence"`
	Receipts    int64 `json:"receipts"`
	Patterns    int64 `json:"patterns"`
	GlitchDrops int64 `json:"glitchDrops"`
}

func Connect(ctx context.Context, databaseURL string) (*Store, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create PostgreSQL pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return &Store{Pool: pool}, nil
}

func (s *Store) Close() {
	s.Pool.Close()
}

func (s *Store) RunMigrations(ctx context.Context, migrationsDir string) error {
	if _, err := s.Pool.Exec(ctx, `
		create table if not exists schema_migrations (
			version text primary key,
			applied_at timestamptz not null default now()
		);
	`); err != nil {
		return fmt.Errorf("create schema_migrations table: %w", err)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".sql") {
			migrationFiles = append(migrationFiles, name)
		}
	}

	sort.Strings(migrationFiles)

	for _, fileName := range migrationFiles {
		alreadyApplied, err := s.migrationApplied(ctx, fileName)
		if err != nil {
			return err
		}

		if alreadyApplied {
			continue
		}

		path := filepath.Join(migrationsDir, fileName)
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", fileName, err)
		}

		tx, err := s.Pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("begin migration %s: %w", fileName, err)
		}

		if _, err := tx.Exec(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("execute migration %s: %w", fileName, err)
		}

		if _, err := tx.Exec(ctx, `
			insert into schema_migrations (version)
			values ($1)
		`, fileName); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("record migration %s: %w", fileName, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("commit migration %s: %w", fileName, err)
		}
	}

	return nil
}

func (s *Store) migrationApplied(ctx context.Context, version string) (bool, error) {
	var exists bool
	err := s.Pool.QueryRow(ctx, `
		select exists (
			select 1
			from schema_migrations
			where version = $1
		)
	`, version).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}

	return exists, nil
}

func (s *Store) SeedDevPlayer(ctx context.Context) error {
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin dev seed: %w", err)
	}

	defer tx.Rollback(ctx)

	var playerID int64
	if err := tx.QueryRow(ctx, `
		insert into players (display_name)
		values ('Soryn')
		on conflict (display_name) do update
		set updated_at = now()
		returning id
	`).Scan(&playerID); err != nil {
		return fmt.Errorf("upsert dev player: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		insert into companion_states (
			player_id,
			companion_name,
			mood_score,
			satiety,
			connection,
			energy,
			comfort,
			playfulness,
			inspiration,
			cleanliness,
			status,
			last_interaction_at
		)
		values (
			$1,
			'Nami-chan',
			80.00,
			85,
			95,
			75,
			90,
			80,
			80,
			80,
			'awake',
			now()
		)
		on conflict (player_id) do update
		set
			companion_name = excluded.companion_name,
			mood_score = excluded.mood_score,
			satiety = excluded.satiety,
			connection = excluded.connection,
			energy = excluded.energy,
			comfort = excluded.comfort,
			playfulness = excluded.playfulness,
			inspiration = excluded.inspiration,
			cleanliness = excluded.cleanliness,
			status = excluded.status,
			last_interaction_at = excluded.last_interaction_at,
			updated_at = now()
	`, playerID); err != nil {
		return fmt.Errorf("upsert companion state: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		insert into player_resources (player_id)
		values ($1)
		on conflict (player_id) do update
		set updated_at = now()
	`, playerID); err != nil {
		return fmt.Errorf("upsert player resources: %w", err)
	}

	if _, err := tx.Exec(ctx, `
		insert into activity_log (player_id, event_type, message)
		values ($1, 'dev_seed', 'Dev player Soryn and Nami-chan were seeded.')
	`, playerID); err != nil {
		return fmt.Errorf("insert activity log: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit dev seed: %w", err)
	}

	return nil
}

func (s *Store) GetDevPlayerStatus(ctx context.Context) (*PlayerStatus, error) {
	var status PlayerStatus

	err := s.Pool.QueryRow(ctx, `
		select
			p.id,
			p.display_name,
			p.level,
			p.total_xp,
			p.currency_cents,
			c.companion_name,
			c.mood_score,
			c.satiety,
			c.connection,
			c.energy,
			c.comfort,
			c.playfulness,
			c.inspiration,
			c.cleanliness,
			c.status,
			c.last_interaction_at,
			r.fans,
			r.memes,
			r.lost_items,
			r.confidence,
			r.receipts,
			r.patterns,
			r.glitch_drops
		from players p
		join companion_states c on c.player_id = p.id
		join player_resources r on r.player_id = p.id
		where p.display_name = 'Soryn'
	`).Scan(
		&status.Player.ID,
		&status.Player.DisplayName,
		&status.Player.Level,
		&status.Player.TotalXP,
		&status.Player.CurrencyCents,
		&status.Companion.CompanionName,
		&status.Companion.MoodScore,
		&status.Companion.Satiety,
		&status.Companion.Connection,
		&status.Companion.Energy,
		&status.Companion.Comfort,
		&status.Companion.Playfulness,
		&status.Companion.Inspiration,
		&status.Companion.Cleanliness,
		&status.Companion.Status,
		&status.Companion.LastInteractionAt,
		&status.Resources.Fans,
		&status.Resources.Memes,
		&status.Resources.LostItems,
		&status.Resources.Confidence,
		&status.Resources.Receipts,
		&status.Resources.Patterns,
		&status.Resources.GlitchDrops,
	)

	if err != nil {
		return nil, fmt.Errorf("get dev player status: %w", err)
	}

	return &status, nil
}
