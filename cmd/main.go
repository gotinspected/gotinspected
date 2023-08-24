package main

import (
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gotinspected/gotinspected/util"
	"sync"
)

type FileAggregate struct {
	FileName               string
	AuthorToLinesChanged   map[string]int
	LineToLineLength       map[int]int
	LineToLineCommentShare map[int]int
}

func AggregateFileByAuthor(fileBlame *git.BlameResult, resultChan chan *FileAggregate, wg *sync.WaitGroup) error {
	defer wg.Done()

	fileAggregate := &FileAggregate{
		FileName:               fileBlame.Path,
		AuthorToLinesChanged:   make(map[string]int),
		LineToLineLength:       make(map[int]int),
		LineToLineCommentShare: make(map[int]int),
	}

	for lineNumber, line := range fileBlame.Lines {
		author := line.AuthorName
		authorLinesChanged, ok := fileAggregate.AuthorToLinesChanged[author]
		if !ok {
			authorLinesChanged = 1
		} else {
			authorLinesChanged += 1
		}
		lineLength := len(line.Text)
		fileAggregate.LineToLineLength[lineNumber] = lineLength
		fileAggregate.AuthorToLinesChanged[author] = authorLinesChanged
	}

	resultChan <- fileAggregate
	return nil
}

type AuthorAggregate struct {
	AuthorName   string
	LinesChanged int
	// TODO
}

func (authorAgg *AuthorAggregate) String() string {
	return fmt.Sprintf("%s: %d", authorAgg.AuthorName, authorAgg.LinesChanged)
}

type RepoAggregate struct {
	RepoLocation       string
	AggregatesByAuthor map[string]*AuthorAggregate
	// TODO
}

func (repoAgg *RepoAggregate) String() string {
	var res string
	res += fmt.Sprintf("GotInspected Statistics for the Repo '%s'\n", repoAgg.RepoLocation)
	res += "\tAuthor: Lines Changed\n"
	var foundStats = false
	for _, authorAgg := range repoAgg.AggregatesByAuthor {
		res += fmt.Sprintf("\t%s\n", authorAgg.String())
		foundStats = true
	}
	if !foundStats {
		res += "\tNo Stats for the repo where found\n"
	}
	return res
}

// Basic example of how to blame a repository.
func main() {

	// argument parsing
	var repoPath string
	flag.StringVar(&repoPath, "repo", ".", "path to the repository to analyze")
	flag.Parse()
	fmt.Println(fmt.Sprintf("Analyze repo given by path '%s'", repoPath))

	repo, err := git.PlainOpen(repoPath)
	util.ExitOnError(err)
	// Retrieve the branch's HEAD, to then get the HEAD commit.
	headRef, err := repo.Head()
	util.ExitOnError(err)
	// get the commit obj for the current head commit
	commitObj, err := repo.CommitObject(headRef.Hash())
	util.ExitOnError(err)

	var wg sync.WaitGroup
	var channel []chan *FileAggregate
	starterFunc := func(file *object.File) error {
		fileBlame, err := git.Blame(commitObj, file.Name)
		if err != nil {
			return err
		}

		aggChan := make(chan *FileAggregate)
		channel = append(channel, aggChan)
		// TODO: handel error gracefully -> use context
		wg.Add(1)
		go AggregateFileByAuthor(fileBlame, aggChan, &wg)
		return nil
	}
	t, _ := commitObj.Tree()
	err = t.Files().ForEach(starterFunc)
	util.ExitOnError(err) // TODO: gracefull -> use context

	wg.Add(1)
	repoAgg := RepoAggregate{
		RepoLocation:       repoPath,
		AggregatesByAuthor: make(map[string]*AuthorAggregate),
	}
	go func() {
		defer wg.Done()
		for _, c := range channel {
			aggregate := <-c
			for author, linesChanged := range aggregate.AuthorToLinesChanged {
				authorAgg, ok := repoAgg.AggregatesByAuthor[author]
				if !ok {
					authorAgg = &AuthorAggregate{
						AuthorName:   author,
						LinesChanged: linesChanged,
					}
				} else {
					authorAgg.LinesChanged += linesChanged
				}
				repoAgg.AggregatesByAuthor[author] = authorAgg
			}
		}
		fmt.Println(fmt.Sprintf("%s", repoAgg.String()))
	}()

	wg.Wait()
}
