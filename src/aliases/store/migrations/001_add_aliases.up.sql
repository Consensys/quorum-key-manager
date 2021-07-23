begin;

create type alias_kind as enum (
	'string',
	'array'
);

create table if not exists aliases (
	id serial primary key,
	key text,
	registry_name text,
	kind alias_kind,
	value json
);
create unique index on aliases (registry_name, key);
create index on aliases (registry_name);

commit;
