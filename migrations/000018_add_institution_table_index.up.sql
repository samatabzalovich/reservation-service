CREATE INDEX institution_search_idx ON institution USING gin(to_tsvector('simple', name || ' ' || description));
