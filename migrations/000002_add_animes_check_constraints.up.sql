alter table animes add constraint animes_runtime_check check ( runtime >= 0 );

alter table animes add constraint animes_year_check check ( year between 1888 and date_part('year', now()) );

alter table animes add constraint genre_length_check check ( array_length(genres, 1) between 1 and 5);