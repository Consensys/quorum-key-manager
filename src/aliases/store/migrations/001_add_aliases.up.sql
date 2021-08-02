begin;

create table if not exists aliases (
	key text,
	registry_name text,
	value json,

	primary key (registry_name, key)
);

commit;
