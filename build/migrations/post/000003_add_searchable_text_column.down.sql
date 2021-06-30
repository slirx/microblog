drop trigger post_searchable_text_update on post;
drop index searchable_text_idx;
alter table post drop column searchable_text;
