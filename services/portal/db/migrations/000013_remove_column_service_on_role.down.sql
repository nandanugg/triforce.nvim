begin;

alter table role add column service varchar(50);

comment on column role.service is 'deprecated';

commit;
