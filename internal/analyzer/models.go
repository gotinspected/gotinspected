package analyzer

type AuthorStat struct {
	Name          string
	Commits       int
	Insertions    int
	Deletions     int
	RawInsertions int
	RawDeletions  int
	EmptyLines    int
	Comments      int
	CurrentLines  int
}

type Result struct {
	Authors         []*AuthorStat
	TotalCommits    int
	TotalIns        int
	TotalDel        int
	TotalChanges    int
	TotalRawIns     int
	TotalRawDel     int
	TotalRawChanges int
	TotalEmpty      int
	TotalComments   int
	TotalCurrent    int
}
