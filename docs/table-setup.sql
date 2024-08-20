create table the_grid.workout (
                                  wrk_id numeric(15) not null primary key,
                                  wrk_date date not null unique
);
create sequence the_grid.seq_wrk_id start with 1 increment by 1;
insert into the_grid.workout (wrk_id, wrk_date) values (nextval('the_grid.seq_wrk_id'), '2024-08-19');
select * from the_grid.workout;

create table the_grid.exercise (
                                   exr_id numeric(15) not null primary key,
                                   exr_name varchar(256) not null,
                                   exr_group varchar(256) not null,
                                   constraint chk_exr_group check (exr_group in (
                                                                                 'BICEP',
                                                                                 'BACK',
                                                                                 'TRICEP',
                                                                                 'CHEST',
                                                                                 'SHOULDER',
                                                                                 'LEGS',
                                                                                 'ABS',
                                                                                 'CARDIO',
                                                                                 'MISC'
                                       ))
);
create sequence the_grid.seq_exr_id start with 1 increment by 1;

insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Concentrated Curls', 'BICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Seated Curls', 'BICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Straight Bar Curls', 'BICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Head Curls', 'BICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Seated Cable Curls', 'BICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Back Extension', 'BACK');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Seated Back Flies', 'BACK');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Standing Lat Pulldown', 'BACK');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Seated Rows', 'BACK');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Trap Pulls', 'BACK');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'BenchPress', 'CHEST');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Incline Press Dumb Bell', 'CHEST');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Chestovers', 'CHEST');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Dumbbell reach', 'CHEST');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Dumbbell Flies', 'CHEST');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Dumbbell Skull Crusher', 'TRICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Overhead seated press', 'TRICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Push downs', 'TRICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Cable push down', 'TRICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Cable pull over', 'TRICEP');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Dumbbell Flies', 'SHOULDER');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Arnold Press', 'SHOULDER');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Barbell Press', 'SHOULDER');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Lateral dumbbell raise', 'SHOULDER');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Lateral Plate Raise', 'SHOULDER');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Dumbbell Press', 'SHOULDER');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Cable pull down', 'SHOULDER');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Corkscrew', 'ABS');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Situp', 'ABS');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Pullover Machine', 'ABS');
insert into the_grid.exercise (exr_id, exr_name, exr_group) values (nextval('the_grid.seq_exr_id'), 'Running', 'CARDIO');


create table the_grid.set (
                              set_id numeric(15) not null primary key,
                              set_wrk_id numeric(15) not null references the_grid.workout(wrk_id),
                              set_exr_id  numeric(15) not null references the_grid.exercise(exr_id),
                              set_weight numeric(15) not null,
                              set_reps numeric(15) not null
);
create sequence the_grid.seq_set_id start with 1 increment by 1;

insert into the_grid.set (set_id, set_wrk_id, set_exr_id, set_weight, set_reps) values (nextval('the_grid.seq_set_id'), 1, 2, 30, 20);
insert into the_grid.set (set_id, set_wrk_id, set_exr_id, set_weight, set_reps) values (nextval('the_grid.seq_set_id'), 1, 2, 20, 20);
insert into the_grid.set (set_id, set_wrk_id, set_exr_id, set_weight, set_reps) values (nextval('the_grid.seq_set_id'), 1, 16, 20, 10);
insert into the_grid.set (set_id, set_wrk_id, set_exr_id, set_weight, set_reps) values (nextval('the_grid.seq_set_id'), 1, 16, 10, 15);
