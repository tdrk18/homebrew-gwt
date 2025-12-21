package main

import "testing"

func TestDetectRepoRoot(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	root, err := detectRepoRoot()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if root != "/repo" {
		t.Fatalf("root = %q, want /repo", root)
	}
}

func TestListWorktrees(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	wts, err := listWorktrees()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(wts) != 3 {
		t.Fatalf("len = %d, want 3", len(wts))
	}

	if wts[1].Branch != "feature" {
		t.Fatalf("branch = %q, want feature", wts[1].Branch)
	}

	if wts[1].IsDetached != false {
		t.Fatalf("IsDetached = %v, want false", wts[1].IsDetached)
	}

	if wts[2].Branch != "" {
		t.Fatalf("branch = %q, want \"\"", wts[2].Branch)
	}

	if wts[2].IsDetached != true {
		t.Fatalf("IsDetached = %v, want true", wts[2].IsDetached)
	}
}

func TestIsWorktreeDirty_Clean(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	if isWorktreeDirty("/clean") {
		t.Fatalf("expected clean worktree")
	}
}

func TestIsWorktreeDirty_Dirty(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	if !isWorktreeDirty("/dirty") {
		t.Fatalf("expected dirty worktree")
	}
}

func TestIsWorktreeDirty_Error(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	if isWorktreeDirty("/unknown") {
		t.Fatalf("expected false on error")
	}
}

func TestAddWorktree_Success(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	err := addWorktree("success")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAddWorktree_FormatsStderrError(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	err := addWorktree("failure")
	if err == nil {
		t.Fatalf("expected error")
	}

	want := "fatal: a branch named 'foo' already exists"
	if err.Error() != want {
		t.Fatalf("err = %q, want %q", err.Error(), want)
	}
}

func TestRemoveWorktree_Success(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	ctx := Context{
		Path:       "/ok",
		Branch:     "feature",
		IsDetached: false,
	}

	if err := removeWorktree(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveWorktree_DetachedSkipsBranchDelete(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	ctx := Context{
		Path:       "/ok",
		Branch:     "",
		IsDetached: true,
	}

	if err := removeWorktree(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveWorktree_RemoveFails(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	ctx := Context{
		Path:       "/fail-remove",
		Branch:     "feature",
		IsDetached: false,
	}

	err := removeWorktree(ctx)
	if err == nil {
		t.Fatalf("expected error")
	}

	want := "fatal: cannot remove worktree"
	if err.Error() != want {
		t.Fatalf("err = %q, want %q", err.Error(), want)
	}
}

func TestRemoveWorktree_BranchDeleteFails(t *testing.T) {
	withFakeExec(t, fakeExecCommand)

	ctx := Context{
		Path:       "/ok",
		Branch:     "fail-branch",
		IsDetached: false,
	}

	err := removeWorktree(ctx)
	if err == nil {
		t.Fatalf("expected error")
	}

	want := "error: failed to delete branch"
	if err.Error() != want {
		t.Fatalf("err = %q, want %q", err.Error(), want)
	}
}
