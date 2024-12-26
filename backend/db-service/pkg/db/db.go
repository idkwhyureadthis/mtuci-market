package db

import (
	"auth-service/internal/model"
	"bufio"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"golang.org/x/crypto/bcrypt"
)

const CARDS_PER_PAGE = 10

type DB struct {
	conn *sql.DB
}

func New(connURL string) *DB {
	conn, err := sql.Open("pgx", connURL)
	if err != nil {
		log.Fatal(err)
	}
	err = setupMigrations(conn)
	if err != nil {
		log.Fatal(err)
	}
	return &DB{conn: conn}
}

func (d *DB) Stop() {
	d.conn.Close()
}

func setupMigrations(conn *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	err := goose.Up(conn, "internal/migrations")
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteCard(id int) error {
	_, err := d.conn.Exec("DELETE FROM cards WHERE id = $1", id)
	return err
}

func (d *DB) CardsOf(id int) ([]*model.Card, error) {
	cards := []*model.Card{}
	rows, err := d.conn.Query("SELECT c.card_status, c.id, c.price, c.name, u.name, u.dorm_number, u.room, u.telegram_id, c.description FROM cards c INNER JOIN users u ON u.id = c.created_by WHERE c.created_by = $1", id)
	if err != nil {
		return cards, err
	}
	for rows.Next() {
		card := model.Card{}
		rows.Scan(&card.Status, &card.Id, &card.Price, &card.Name, &card.CreatorName, &card.DormNumber, &card.Room, &card.Telegram, &card.Description)
		photos, err := d.conn.Query("SELECT photo_link from card_photos WHERE card_id = $1", card.Id)
		if err != nil {
			return cards, err
		}
		for photos.Next() {
			var cardPhoto string
			photos.Scan(&cardPhoto)
			card.Photos = append(card.Photos, cardPhoto)
		}
		cards = append(cards, &card)
	}
	return cards, nil
}

func (d *DB) VerifyToken(id int, token string) error {
	s := sha256.New()
	s.Write([]byte(token))
	var cryptedStored string
	err := d.conn.QueryRow("SELECT crypted_refresh FROM users WHERE id = $1", id).Scan(&cryptedStored)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(cryptedStored), s.Sum(nil)); err != nil {
		return errors.New("wrong refresh provided")
	}
	return nil
}

func (d *DB) GetUserData(id int) (*model.UserData, error) {
	data := model.UserData{}
	err := d.conn.QueryRow("SELECT name, room, dorm_number FROM users WHERE id = $1", id).Scan(&data.Name, &data.Room, &data.DormNumber)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (d *DB) CreateNewUser(ctx context.Context, name, telegram, password, room, dormNumber, login, role string) (int, error) {
	var id int
	h := sha256.New()
	h.Write([]byte(password))
	cryptedPassword := hex.EncodeToString(h.Sum(nil))

	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer tx.Rollback()
	err = tx.QueryRow("INSERT INTO users (name, crypted_password, telegram_id, room, dorm_number, login, role) values ($1, $2, $3, $4, $5, $6, $7) RETURNING id;", name, cryptedPassword, telegram, room, dormNumber, login, role).Scan(&id)
	if err != nil && strings.HasSuffix(err.Error(), "(SQLSTATE 23505)") {
		return -1, errors.New("login already occupied")
	}
	if err != nil {
		return -1, err
	}
	if err := tx.Commit(); err != nil {
		return -1, err
	}
	return id, nil
}

func (d *DB) SetRefreshKey(id int, refresh string) error {
	s := sha256.New()
	s.Write([]byte(refresh))
	cryptedRefreshBytes, err := bcrypt.GenerateFromPassword(s.Sum(nil), 12)
	if err != nil {
		return err
	}
	cryptedRefresh := string(cryptedRefreshBytes)
	_, err = d.conn.Exec("UPDATE users SET crypted_refresh = $1 WHERE id = $2", cryptedRefresh, id)
	if err != nil {
		return err
	}
	return nil
}

func (d *DB) DeleteUser(ctx context.Context, id int) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	_, err = tx.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *DB) CreateMessage(ctx context.Context, from, to int, content string) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	_, err = tx.Exec("INSERT INTO messages (author, recipient, content) VALUES ($1, $2, $3)", from, to, content)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *DB) CreateCard(ctx context.Context, id int, name string, price float64, photos []string, description string) error {
	var cardID int
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	err = tx.QueryRow("INSERT INTO cards (created_by, name, card_status, price, description) VALUES ($1, $2, $3, $4, $5) RETURNING id;", id, name, "on moderation", price, description).Scan(&cardID)
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	tx2, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx2.Rollback()
	for _, photo := range photos {
		_, err = tx2.Exec("INSERT INTO card_photos (card_id, photo_link) VALUES ($1, $2)", cardID, photo)
		if err != nil {
			return err
		}
	}
	if err = tx2.Commit(); err != nil {
		return err
	}
	return nil
}

func (d *DB) GetOnModeration() (*model.Card, error) {
	card := model.Card{}
	row := d.conn.QueryRow("SELECT c.card_status, c.id, c.price, c.name, u.name, u.dorm_number, u.room, u.telegram_id, c.description FROM cards c INNER JOIN users u ON u.id = c.created_by WHERE c.id = (SELECT MIN(id) from cards where card_status = 'on moderation') and card_status='on moderation'")
	err := row.Scan(&card.Status, &card.Id, &card.Price, &card.Name, &card.CreatorName, &card.DormNumber, &card.Room, &card.Telegram, &card.Description)
	if err == sql.ErrNoRows {
		return &card, nil
	}
	if err != nil {
		return nil, err
	}
	photos, err := d.conn.Query("SELECT photo_link from card_photos WHERE card_id = $1", card.Id)
	if err != nil {
		return nil, err
	}
	for photos.Next() {
		var cardPhoto string
		photos.Scan(&cardPhoto)
		card.Photos = append(card.Photos, cardPhoto)
	}
	return &card, nil
}

func (d *DB) GetImage(name string) (string, error) {
	file, err := os.Open("internal/images/" + name)
	if err != nil {
		return "", err
	}
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil {
		return "", err
	}
	enc := base64.StdEncoding.EncodeToString(bs)
	return enc, nil
}

func (d *DB) GetModerated() ([]*model.Card, error) {
	cards := []*model.Card{}
	rows, err := d.conn.Query("SELECT c.id, c.price, c.name, u.name, u.dorm_number, u.room, u.telegram_id, c.description FROM cards c INNER JOIN users u ON u.id = c.created_by WHERE card_status = 'accepted'")
	if err != nil {
		return cards, err
	}
	for rows.Next() {
		card := model.Card{}
		rows.Scan(&card.Id, &card.Price, &card.Name, &card.CreatorName, &card.DormNumber, &card.Room, &card.Telegram, &card.Description)
		photos, err := d.conn.Query("SELECT photo_link from card_photos WHERE card_id = $1", card.Id)
		if err != nil {
			return cards, err
		}
		for photos.Next() {
			var cardPhoto string
			photos.Scan(&cardPhoto)
			card.Photos = append(card.Photos, cardPhoto)
		}
		cards = append(cards, &card)
	}
	return cards, nil
}

func (d *DB) VerifyPassword(ctx context.Context, login, password string) (int, string, error) {
	var id int
	var role string
	h := sha256.New()
	h.Write([]byte(password))
	cryptedPassword := hex.EncodeToString(h.Sum(nil))
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return -1, "", err
	}
	defer tx.Rollback()
	err = tx.QueryRow("SELECT id, role from users WHERE login = $1 AND crypted_password =  $2;", login, cryptedPassword).Scan(&id, &role)
	if err != nil {
		return -1, "", err
	}
	if id == 0 {
		return -1, "", errors.New("user not found")
	}
	if err := tx.Commit(); err != nil {
		return -1, "", err
	}
	return id, role, nil
}

func (d *DB) ChangeCardStatus(ctx context.Context, newStatus string, id int) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.Exec("UPDATE cards SET card_status = $1 WHERE id = $2", newStatus, id)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
