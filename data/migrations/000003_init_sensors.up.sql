CREATE TABLE IF NOT EXISTS sensors (
    id serial primary key,
    district_id int references airquality_db.public.districts (id),
    api_id int,
    address varchar(600),
    lat float4,
    lon float4
);