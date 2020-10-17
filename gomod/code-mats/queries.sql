/* Number of sha's */
select count(distinct sha) from log;

/* Files touched the most (exchange filename for any other field) */
select count(*), filename from log group by filename order by count(*) desc;

select count(*) from log a, log b where a.sha=b.sha and a.filename = 'migrations/v0.0.8-v1.0.0/manually.js' and b.filename = 'migrations/v0.0.8-v1.0.0/package.json';

select count(*), a.filename, b.filename from log a, log b where a.sha=b.sha and a.filename = (select filename from log) and b.filename = (select filename from log) and a.filename = b.filename;

select a.filename, b.filename from log a, log b where a.sha=b.sha and a.filename != b.filename;

with filenames as (
  select a.filename as aname, b.filename as bname from log a, log b where a.sha=b.sha and a.filename != b.filename
)
select count(*) from log a, log b, filenames where a.sha=b.sha and a.filename = filenames.aname and b.filename = filenames.bname;
