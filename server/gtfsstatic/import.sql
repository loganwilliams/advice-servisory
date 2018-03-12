DROP TABLE IF EXISTS stops CASCADE;
DROP TABLE IF EXISTS routes CASCADE;

CREATE SCHEMA IF NOT EXISTS import;
DROP TABLE IF EXISTS import.stops;
CREATE TABLE import.stops(
	stop_id text,
	stop_code text default '0',
	stop_name text,
	stop_desc text,
	stop_lat text,
	stop_lon text,
	zone_id text,
	stop_url text,
	location_type text,
	parent_station text
);
COPY import.stops FROM '/Users/loganw/go/src/github.com/loganwilliams/adviceservisory/server/gtfs-static/stops.txt' WITH DELIMITER ',' HEADER CSV;

CREATE TABLE import.routes(
	route_id text,
	agency_id text,
	route_short_name text,
	route_long_name text,
	route_desc text,
	route_type text,
	route_url text,
	route_color text,
	route_text_color text
);
COPY import.routes FROM '/Users/loganw/go/src/github.com/loganwilliams/adviceservisory/server/gtfs-static/routes.txt' WITH DELIMITER ',' HEADER CSV;

CREATE TABLE stops(
	id serial primary key,
	code text,
	name varchar(200),
	latitude real,
	longitude real,
	location_type int,
	parent_station text
);

CREATE TABLE routes(
	id serial primary key,
	code varchar(5),
	short_name varchar(5),
	name varchar(200),
	description text,
	type int,
	url varchar(200),
	color varchar(6)
);

INSERT INTO stops(
	code,
	name,
	latitude,
	longitude,
	location_type,
	parent_station
) SELECT
	import.stops.stop_id AS code,
	import.stops.stop_name AS name,
	import.stops.stop_lat::real AS latitude,
	import.stops.stop_lon::real AS longitude,
	import.stops.location_type::int AS location_type,
	import.stops.parent_station AS parent_station
FROM import.stops;

INSERT INTO routes(
	code,
	short_name,
	name,
	description,
	type,
	url,
	color
) SELECT
	import.routes.route_id AS code,
	import.routes.route_short_name as short_name,
	import.routes.route_long_name AS name,
	import.routes.route_desc AS description,
	import.routes.route_type::int AS type,
	import.routes.route_url AS url,
	import.routes.route_color AS color
FROM import.routes;
