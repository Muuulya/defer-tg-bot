package bot

import (
	"encoding/json"
	"log"
	"time"
)

func (b *Bot) StartScheduler() {
	for {
		rows, err := b.DB.Query("SELECT id, channel_id, user_id, message_ids, scheduled_time FROM scheduled_posts")
		if err != nil {
			log.Println("Scheduler DB error:", err)
			time.Sleep(time.Minute)
			continue
		}

		for rows.Next() {
			var id, channelID, userID int
			var messageIDsJSON string
			var scheduledTimeStr string

			if err := rows.Scan(&id, &channelID, &userID, &messageIDsJSON, &scheduledTimeStr); err != nil {
				log.Println("Row scan error:", err)
				continue
			}

			var messageIDs []int
			_ = json.Unmarshal([]byte(messageIDsJSON), &messageIDs)

			scheduledTime, _ := time.Parse(time.RFC3339, scheduledTimeStr)
			now := time.Now()

			if scheduledTime.Before(now) || scheduledTime.Equal(now) {
				go b.PublishMessages(id, int64(channelID), int64(userID), messageIDs)
			} else {
				delay := time.Until(scheduledTime)
				go func(postID int, chID, uID int64, msgIDs []int, d time.Duration) {
					time.Sleep(d)
					b.PublishMessages(postID, chID, uID, msgIDs)
				}(id, int64(channelID), int64(userID), messageIDs, delay)
			}
		}

		rows.Close()
		time.Sleep(time.Minute)
	}
}
