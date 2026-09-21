-- Deploy activity-reports:1.0.0/1_create_activity-report_table to pg

BEGIN;

CREATE TABLE activity_reports (
    id                  UUID            PRIMARY KEY DEFAULT gen_random_uuid(),
    month               INTEGER         NOT NULL,
    year                INTEGER         NOT NULL,
    prospection_days    DECIMAL         NOT NULL DEFAULT 0,
    vacation_days       DECIMAL         NOT NULL DEFAULT 0,
    days_away           DECIMAL         NOT NULL DEFAULT 0,
    formation_days      DECIMAL         NOT NULL DEFAULT 0,
    created_at          TIMESTAMP       NOT NULL DEFAULT NOW()
);

CREATE TABLE report_items (
    id              UUID    PRIMARY KEY DEFAULT gen_random_uuid(),
    mission_id      UUID    NOT NULL,
    report_id       UUID    NOT NULL,
    nb_days         DECIMAL NOT NULL,
    progress_report TEXT    NOT NULL,

    CONSTRAINT fk_activity_reports  FOREIGN KEY (report_id) REFERENCES activity_reports (id),
    CONSTRAINT unique_item          UNIQUE (mission_id, report_id)
);

COMMIT;
