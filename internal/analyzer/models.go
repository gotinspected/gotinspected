package analyzer

// AuthorStat holds the activity metrics for a single author.
type AuthorStat struct {
	Name       string
	Commits    int
	Insertions int
	Deletions  int
}

// Result holds the aggregated output of a repository analysis.
type Result struct {
	Authors      []*AuthorStat
	TotalCommits int
	TotalIns     int
	TotalDel     int
	TotalChanges int
}
