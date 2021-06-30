ALTER TABLE post
    ADD COLUMN searchable_text tsvector;

CREATE INDEX searchable_text_idx ON post USING GIN (to_tsvector('english', post.text));

CREATE TRIGGER post_searchable_text_update
    BEFORE INSERT OR UPDATE
    ON post
    FOR EACH ROW
EXECUTE PROCEDURE
    tsvector_update_trigger(searchable_text, 'pg_catalog.english', text);
