CREATE TABLE notifications (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type             VARCHAR(50) NOT NULL,
    title            TEXT NOT NULL,
    body             TEXT NOT NULL,
    is_read          BOOLEAN NOT NULL DEFAULT FALSE,
    related_log_id   BIGINT REFERENCES sync_logs(id) ON DELETE SET NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX notifications_user_unread ON notifications(user_id, is_read) WHERE is_read = FALSE;

COMMENT ON TABLE  notifications                  IS 'インアプリ通知';
COMMENT ON COLUMN notifications.type             IS 'SYNC_COMPLETED | SYNC_FAILED';
COMMENT ON COLUMN notifications.is_read          IS '既読フラグ';
COMMENT ON COLUMN notifications.related_log_id   IS '関連する同期ログID';
