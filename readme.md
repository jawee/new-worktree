# new-worktree

Tool for easy creation of bug/feature/fix worktrees without having to remember the exact syntax.

Executes
```bash
git worktree add {short-type}-branch-name -b {type}/branch-name
```

Where {short-type} is
* feat
* bug
* fix

and {type} is
* feature
* bug
* fix

