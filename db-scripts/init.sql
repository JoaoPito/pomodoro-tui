--- PROJECTS ---
CREATE TABLE IF NOT EXISTS public.projects
(
    id integer NOT NULL DEFAULT nextval('projects_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'::text),
    repository character varying(500) COLLATE pg_catalog."default",
    creator character varying(500) COLLATE pg_catalog."default",
    archived boolean NOT NULL DEFAULT false,
    updated_at timestamp with time zone NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'::text),
    CONSTRAINT projects_pkey PRIMARY KEY (id),
    CONSTRAINT projects_name_unique UNIQUE (name),
    CONSTRAINT projects_name_not_empty CHECK (length(TRIM(BOTH FROM name)) > 0)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.projects
    OWNER to postgres;

COMMENT ON TABLE public.projects
    IS 'Stores project information with git repository references';
-- Index: idx_projects_archived

-- DROP INDEX IF EXISTS public.idx_projects_archived;

CREATE INDEX IF NOT EXISTS idx_projects_archived
    ON public.projects USING btree
    (archived ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: idx_projects_created_at

-- DROP INDEX IF EXISTS public.idx_projects_created_at;

CREATE INDEX IF NOT EXISTS idx_projects_created_at
    ON public.projects USING btree
    (created_at DESC NULLS FIRST)
    TABLESPACE pg_default;
-- Index: idx_projects_repository

-- DROP INDEX IF EXISTS public.idx_projects_repository;

CREATE INDEX IF NOT EXISTS idx_projects_repository
    ON public.projects USING btree
    (repository COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE repository IS NOT NULL;

-- Trigger: update_projects_updated_at

-- DROP TRIGGER IF EXISTS update_projects_updated_at ON public.projects;

CREATE OR REPLACE TRIGGER update_projects_updated_at
    BEFORE UPDATE 
    ON public.projects
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at_column();

--- TASKS ---
CREATE TABLE IF NOT EXISTS public.tasks
(
    id integer NOT NULL DEFAULT nextval('tasks_id_seq'::regclass),
    name character varying(500) COLLATE pg_catalog."default" NOT NULL,
    project_id integer,
    created_at timestamp with time zone NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'::text),
    completed boolean NOT NULL DEFAULT false,
    estimated_duration_min integer,
    deadline timestamp with time zone,
    times_procrastinated integer NOT NULL DEFAULT 0,
    updated_at timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'::text),
    deleted_at timestamp with time zone,
    completed_at timestamp with time zone,
    description text COLLATE pg_catalog."default",
    priority integer DEFAULT 0,
    CONSTRAINT tasks_pkey PRIMARY KEY (id),
    CONSTRAINT fk_tasks_project FOREIGN KEY (project_id)
        REFERENCES public.projects (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
    CONSTRAINT tasks_completed_at_logic CHECK (completed = false AND completed_at IS NULL OR completed = true AND completed_at IS NOT NULL),
    CONSTRAINT tasks_estimated_duration_positive CHECK (estimated_duration_min IS NULL OR estimated_duration_min > 0),
    CONSTRAINT tasks_name_not_empty CHECK (length(TRIM(BOTH FROM name)) > 0),
    CONSTRAINT tasks_priority_valid CHECK (priority = ANY (ARRAY[NULL::integer, 0, 1, 2, 3])),
    CONSTRAINT tasks_times_procrastinated_non_negative CHECK (times_procrastinated >= 0)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.tasks
    OWNER to postgres;

COMMENT ON TABLE public.tasks
    IS 'Individual tasks associated with projects';

COMMENT ON COLUMN public.tasks.times_procrastinated
    IS 'Counter for how many times the task was postponed';
-- Index: idx_tasks_completed

-- DROP INDEX IF EXISTS public.idx_tasks_completed;

CREATE INDEX IF NOT EXISTS idx_tasks_completed
    ON public.tasks USING btree
    (completed ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: idx_tasks_created_at

-- DROP INDEX IF EXISTS public.idx_tasks_created_at;

CREATE INDEX IF NOT EXISTS idx_tasks_created_at
    ON public.tasks USING btree
    (created_at DESC NULLS FIRST)
    TABLESPACE pg_default;
-- Index: idx_tasks_deadline

-- DROP INDEX IF EXISTS public.idx_tasks_deadline;

CREATE INDEX IF NOT EXISTS idx_tasks_deadline
    ON public.tasks USING btree
    (deadline ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE deadline IS NOT NULL;
-- Index: idx_tasks_priority

-- DROP INDEX IF EXISTS public.idx_tasks_priority;

CREATE INDEX IF NOT EXISTS idx_tasks_priority
    ON public.tasks USING btree
    (priority ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: idx_tasks_project_completed

-- DROP INDEX IF EXISTS public.idx_tasks_project_completed;

CREATE INDEX IF NOT EXISTS idx_tasks_project_completed
    ON public.tasks USING btree
    (project_id ASC NULLS LAST, completed ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: idx_tasks_project_id

-- DROP INDEX IF EXISTS public.idx_tasks_project_id;

CREATE INDEX IF NOT EXISTS idx_tasks_project_id
    ON public.tasks USING btree
    (project_id ASC NULLS LAST)
    TABLESPACE pg_default;

-- Trigger: set_task_completed_at_trigger

-- DROP TRIGGER IF EXISTS set_task_completed_at_trigger ON public.tasks;

CREATE OR REPLACE TRIGGER set_task_completed_at_trigger
    BEFORE UPDATE 
    ON public.tasks
    FOR EACH ROW
    EXECUTE FUNCTION public.set_task_completed_at();

-- Trigger: update_tasks_updated_at

-- DROP TRIGGER IF EXISTS update_tasks_updated_at ON public.tasks;

CREATE OR REPLACE TRIGGER update_tasks_updated_at
    BEFORE UPDATE 
    ON public.tasks
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at_column();

--- FOCUS SESSIONS ---
CREATE TABLE IF NOT EXISTS public.focus
(
    id uuid NOT NULL DEFAULT gen_random_uuid(),
    start_time timestamp with time zone NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'::text),
    end_time timestamp with time zone,
    task_id integer,
    end_status character varying(50) COLLATE pg_catalog."default",
    device character varying(250) COLLATE pg_catalog."default",
    capture_mode character varying(50) COLLATE pg_catalog."default",
    session character varying(50) COLLATE pg_catalog."default",
    duration_minutes integer GENERATED ALWAYS AS (
CASE
    WHEN (end_time IS NOT NULL) THEN (EXTRACT(epoch FROM (end_time - start_time)) / (60)::numeric)
    ELSE NULL::numeric
END) STORED,
    CONSTRAINT focus_pkey PRIMARY KEY (id),
    CONSTRAINT fk_focus_task FOREIGN KEY (task_id)
        REFERENCES public.tasks (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL,
    CONSTRAINT focus_end_time_after_start CHECK (end_time IS NULL OR end_time > start_time)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.focus
    OWNER to postgres;

COMMENT ON TABLE public.focus
    IS 'Focus sessions (pomodoro) tracking actual time spent on tasks';

COMMENT ON COLUMN public.focus.capture_mode
    IS 'How the focus session was captured (manual, timer, etc)';
-- Index: idx_focus_active_sessions

-- DROP INDEX IF EXISTS public.idx_focus_active_sessions;

CREATE INDEX IF NOT EXISTS idx_focus_active_sessions
    ON public.focus USING btree
    (end_time ASC NULLS LAST)
    TABLESPACE pg_default
    WHERE end_time IS NULL;
-- Index: idx_focus_device

-- DROP INDEX IF EXISTS public.idx_focus_device;

CREATE INDEX IF NOT EXISTS idx_focus_device
    ON public.focus USING btree
    (device COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: idx_focus_end_status

-- DROP INDEX IF EXISTS public.idx_focus_end_status;

CREATE INDEX IF NOT EXISTS idx_focus_end_status
    ON public.focus USING btree
    (end_status COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: idx_focus_start_time

-- DROP INDEX IF EXISTS public.idx_focus_start_time;

CREATE INDEX IF NOT EXISTS idx_focus_start_time
    ON public.focus USING btree
    (start_time DESC NULLS FIRST)
    TABLESPACE pg_default;
-- Index: idx_focus_task_id

-- DROP INDEX IF EXISTS public.idx_focus_task_id;

CREATE INDEX IF NOT EXISTS idx_focus_task_id
    ON public.focus USING btree
    (task_id ASC NULLS LAST)
    TABLESPACE pg_default;
