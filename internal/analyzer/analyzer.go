package analyzer

import (
	"sort"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/yourusername/gitinspector-go/internal/plugins"
)

type GitAnalyzer struct {
	repoPath string
	registry *plugins.Registry
}

func New(repoPath string, registry *plugins.Registry) *GitAnalyzer {
	return &GitAnalyzer{
		repoPath: repoPath,
		registry: registry,
	}
}

func (a *GitAnalyzer) Analyze() (*Result, error) {
	// FIX: Use DetectDotGit so it automatically finds the repo root even if
	// the positional path points to a deep subdirectory.
	repo, err := git.PlainOpenWithOptions(a.repoPath, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	head, err := repo.Head()
	if err != nil {
		return nil, err
	}

	commitIter, err := repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		return nil, err
	}

	statsMap := make(map[string]*AuthorStat)

	err = commitIter.ForEach(func(c *object.Commit) error {
		authorName := c.Author.Name
		if _, exists := statsMap[authorName]; !exists {
			statsMap[authorName] = &AuthorStat{Name: authorName}
		}
		statsMap[authorName].Commits++

		currentTree, err := c.Tree()
		if err != nil {
			return err
		}

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

		changes, err := object.DiffTree(parentTree, currentTree)
		if err != nil {
			return err
		}

		patch, err := changes.Patch()
		if err != nil {
			return err
		}

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
				for _, line := range lines {
					if line == "" {
						continue
					}

					if plugin.IsSignificant(line) {
						if chunk.Type() == diff.Add {
							statsMap[authorName].Insertions++
						} else if chunk.Type() == diff.Delete {
							statsMap[authorName].Deletions++
						}
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return a.aggregate(statsMap), nil
}

func (a *GitAnalyzer) aggregate(statsMap map[string]*AuthorStat) *Result {
	res := &Result{}
	for _, stat := range statsMap {
		res.Authors = append(res.Authors, stat)
		res.TotalCommits += stat.Commits
		res.TotalIns += stat.Insertions
		res.TotalDel += stat.Deletions
	}
	res.TotalChanges = res.TotalIns + res.TotalDel

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
