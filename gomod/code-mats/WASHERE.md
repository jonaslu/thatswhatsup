# WASHERE
I made covariance work by using a set.

Now I need somewhere to store the things and count the instances.
I'm thinking a hashkey, but I can also use an array?

Or, I want some way to simply count stuff.

# Considerations
If there are multiple merges, I will use the last commit as the real commit,
the rest seems like only merge stuff.

Example:

Also, binary files does not have added and deleted, they have - -
so I'm translating this into "0" "0" to have something.

# Next
Massage data from the database in psql

# Useful queries
queries.sql

# Get all files...
For all files, log each commit they appear in

Then, for each file:
For each commit the file appears in, give one point to all files in that commit

File - { other files, points }
Scan through to find the highest number

What do I want?
I want the set of all files that are in more than one commit together

So, for all files:
Generate the set of all combinations
Then for each entry in the set - get the # of commits that they appear in

select commit, count(*) from log a, log b where a.sha=b.sha and a.filename = 'XX' and b.filename = 'YY' group by commit having count(*) > 1;
