CREATE INDEX idx_project_meta_canonical ON projects USING GIN ((meta_info->'canonical'));
-- Using expression index
CREATE INDEX idx_project_id_canonical_expr ON projects (id, (meta_info->>'canonical'));

CREATE INDEX IF NOT EXISTS idx_project_canonical ON projects ((meta_info->>'canonical')) WHERE meta_info->>'canonical' IS NOT NULL;