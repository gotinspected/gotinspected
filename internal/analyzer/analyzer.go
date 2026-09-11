package analyzer

import (
	"sort"
	"strings"
	"sync"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/yourusername/gitinspector-go/internal/plugins"
)

type GitAnalyzer struct {
	repoPath string
	registry *plugins.Registry
}

func New(repoPath string, registry *plugins.Registry) *GitAnalyzer {
	return &GitAnalyzer{repoPath: repoPath, registry: registry}
}

func (a *GitAnalyzer) Analyze() (*Result, error) {
	repo, err := git.PlainOpenWithOptions(a.repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}

	headRef, err := repo.Head()
	if err != nil {
		return nil, err
	}

	headCommit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		return nil, err
	}

	// 1. Initialize concurrency tools and private stat maps
	var wg sync.WaitGroup
	wg.Add(2)

	diffStats := make(map[string]*AuthorStat)
	blameStats := make(map[string]*AuthorStat)
	var diffErr, blameErr error

	// 2. Run Historical Diff Analysis concurrently
	go func() {
		defer wg.Done()
		diffErr = a.runDiffPhase(repo, headRef, diffStats)
	}()

	// 3. Run Blame Analysis concurrently
	go func() {
		defer wg.Done()
		blameErr = a.runBlamePhase(repo, headCommit, blameStats)
	}()

	// Wait for both phases to complete
	wg.Wait()

	if diffErr != nil {
		return nil, diffErr
	}
	if blameErr != nil {
		return nil, blameErr
	}

	// 4. Safely merge the results from both phases
	finalStats := make(map[string]*AuthorStat)
	for name, stat := range diffStats {
		finalStats[name] = stat
	}
	for name, bStat := range blameStats {
		if _, exists := finalStats[name]; !exists {
			finalStats[name] = &AuthorStat{Name: name}
		}
		finalStats[name].CurrentLines += bStat.CurrentLines
	}

	return a.aggregate(finalStats), nil
}

func (a *GitAnalyzer) runDiffPhase(repo *git.Repository, headRef *plumbing.Reference, statsMap map[string]*AuthorStat) error {
	commitIter, err := repo.Log(&git.LogOptions{From: headRef.Hash()})
	if err != nil {
		return err
	}

	return commitIter.ForEach(func(c *object.Commit) error {
		authorName := c.Author.Name
		if _, exists := statsMap[authorName]; !exists {
			statsMap[authorName] = &AuthorStat{Name: authorName}
		}
		statsMap[authorName].Commits++

		currentTree, _ := c.Tree()
		var parentTree *object.Tree
		if c.NumParents() > 0 {
			parent, err := c.Parent(0)
			if err == nil {
				parentTree, _ = parent.Tree()
			}
		}
		if parentTree == nil {
			parentTree = &object.Tree{}
		}

		changes, _ := object.DiffTree(parentTree, currentTree)
		patch, _ := changes.Patch()

		for _, fp := range patch.FilePatches() {
			from, to := fp.Files()
			filename := ""
			if to != nil {
				filename = to.Path()
			} else if from != nil {
				filename = from.Path()
			}

			plugin := a.registry.Get(filename)
			if plugin.ShouldExclude(filename) {
				continue
			}

			for _, chunk := range fp.Chunks() {
				if chunk.Type() == diff.Equal {
					continue
				}

				lines := strings.Split(chunk.Content(), "\n")
				fileAnalyzer := plugin.NewAnalyzer()

				for i, line := range lines {
					if i == len(lines)-1 && line == "" {
						continue
					}

					if chunk.Type() == diff.Add {
						statsMap[authorName].RawInsertions++
					}
					if chunk.Type() == diff.Delete {
						statsMap[authorName].RawDeletions++
					}

					switch fileAnalyzer.AnalyzeLine(line) {
					case plugins.TypeCode:
						if chunk.Type() == diff.Add {
							statsMap[authorName].Insertions++
						}
						if chunk.Type() == diff.Delete {
							statsMap[authorName].Deletions++
						}
					case plugins.TypeEmpty:
						statsMap[authorName].EmptyLines++
					case plugins.TypeComment:
						statsMap[authorName].Comments++
					}
				}
			}
		}
		return nil
	})
}

func (a *GitAnalyzer) runBlamePhase(repo *git.Repository, headCommit *object.Commit, statsMap map[string]*AuthorStat) error {
	headTree, err := headCommit.Tree()
	if err != nil {
		return err
	}

	commitCache := make(map[plumbing.Hash]*object.Commit)

	return headTree.Files().ForEach(func(f *object.File) error {
		plugin := a.registry.Get(f.Name)
		if plugin.ShouldExclude(f.Name) {
			return nil
		}

		blameResult, err := git.Blame(headCommit, f.Name)
		if err != nil {
			return nil
		}

		fileAnalyzer := plugin.NewAnalyzer()

		for _, line := range blameResult.Lines {
			c, cached := commitCache[line.Hash]
			if !cached {
				c, _ = repo.CommitObject(line.Hash)
				commitCache[line.Hash] = c
			}

			authorName := line.Author
			if c != nil {
				authorName = c.Author.Name
			}

			if _, exists := statsMap[authorName]; !exists {
				statsMap[authorName] = &AuthorStat{Name: authorName}
			}

			if fileAnalyzer.AnalyzeLine(line.Text) == plugins.TypeCode {
				statsMap[authorName].CurrentLines++
			}
		}
		return nil
	})
}

func (a *GitAnalyzer) aggregate(statsMap map[string]*AuthorStat) *Result {
	res := &Result{}
	for _, stat := range statsMap {
		res.Authors = append(res.Authors, stat)
		res.TotalCommits += stat.Commits
		res.TotalIns += stat.Insertions
		res.TotalDel += stat.Deletions
		res.TotalRawIns += stat.RawInsertions
		res.TotalRawDel += stat.RawDeletions
		res.TotalEmpty += stat.EmptyLines
		res.TotalComments += stat.Comments
		res.TotalCurrent += stat.CurrentLines
	}
	res.TotalChanges = res.TotalIns + res.TotalDel
	res.TotalRawChanges = res.TotalRawIns + res.TotalRawDel

	sort.Slice(res.Authors, func(i, j int) bool {
		impactI := res.Authors[i].Insertions + res.Authors[i].Deletions
		impactJ := res.Authors[j].Insertions + res.Authors[j].Deletions
		if impactI == impactJ {
			return res.Authors[i].Commits > res.Authors[j].Commits
		}
		return impactI > impactJ
	})
	return res
}
