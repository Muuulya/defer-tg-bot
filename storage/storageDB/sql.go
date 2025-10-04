package storageDB

const (
	// Init
	initSQL = `
		PRAGMA journal_mode=WAL;
		PRAGMA foreign_keys=ON;
		PRAGMA busy_timeout=5000;
	`
	createTablsSQL = `
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			user_name TEXT,
			current_state_name TEXT,
			current_dialog_message_id INTEGER,
			current_channel_page INTAGER,
			selected_channel_id INTEGER
		);

		CREATE TABLE IF NOT EXISTS channels (
			user_id INTEGER,
			channel_id INTEGER,
			channel_name TEXT,
			PRIMARY KEY (user_id, channel_id),
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS scheduled_posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			channel_id INTEGER NOT NULL,
			message_id INTEGER NOT NULL,
			scheduled_time INTEGER NOT NULL,
			is_posted INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
			FOREIGN KEY (user_id, channel_id) REFERENCES channels(user_id, channel_id) ON DELETE CASCADE
		);

		CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_post
		ON scheduled_posts(user_id, channel_id, message_id, scheduled_time);
	`

	// User
	addUserSQL = `INSERT OR IGNORE INTO users(user_id, user_name, current_state_name, current_dialog_message_id, current_channel_page, selected_channel_id) 
		VALUES(?, ?, ?, ?, ?, ?)`
	updateUserStateSQL                = `UPDATE users SET current_state_name = ? WHERE user_id = ?`
	updateUserCurrentDialogMessageSQL = `UPDATE users SET current_dialog_message_id = ? WHERE user_id = ?`
	updateUserSelectedChannelIDSQL    = `UPDATE users SET selected_channel_id = ? WHERE user_id = ?`
	updateUserCurrentChannelPageSQL   = `UPDATE users SET current_channel_page = ? WHERE user_id = ?`
	getUserSQL                        = `SELECT * FROM users WHERE user_id = ? LIMIT 1`
	getAllUsersSQL                    = `SELECT * FROM users`

	// Channel
	addChannelSQL = `INSERT OR IGNORE INTO channels(user_id, channel_id, channel_name)
		VALUES (?, ?, ?);`
	updateChannelNameSQL = `UPDATE channels SET channel_name = ? WHERE user_id = ? AND channel_id = ?`
	getChannelSQL        = `SELECT channel_id, channel_name FROM channels WHERE user_id = ? AND channel_id = ? LIMIT 1`
	getUserChannelsSQL   = `SELECT channel_id, channel_name FROM channels WHERE user_id = ?`
	removeChannelSQL     = `DELETE FROM channels WHERE user_id = ? AND channel_id = ?`

	// Message
	addMessageSQL = `INSERT OR IGNORE INTO scheduled_posts(user_id, channel_id, message_id, scheduled_time, is_posted)
		VALUES (?, ?, ?, ?, ?);`
	getAllMessagesSQL = `SELECT id, user_id, channel_id, message_id, scheduled_time, is_posted FROM scheduled_posts`

	getUserChannelMessagesAfterSQL = `
		SELECT id, user_id, channel_id, message_id, scheduled_time, is_posted
		FROM scheduled_posts
		WHERE user_id = ? AND channel_id = ? AND scheduled_time >= ?
	`
	getAllMissedMessagesBeforSQL = `
		SELECT id, user_id, channel_id, message_id, scheduled_time, is_posted
		FROM scheduled_posts
		WHERE scheduled_time < ? AND is_posted = ?
	`

	updateMessageStatusSQL = `
		UPDATE scheduled_posts SET is_posted = ? WHERE id = ?
	`

	removeMessageSQL = `
		DELETE FROM scheduled_posts
		WHERE id = ?
	`
)
