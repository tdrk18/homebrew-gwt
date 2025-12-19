# wt

自分用の git worktree TUI ラッパー。

`git worktree` の煩雑な操作を、  
**一覧 → 選択 → cd** だけで完結させる。

---

## 使い方

```sh
gwt
````

### キー操作

* `j / k` : 移動（vim key）
* `Enter` : 選択した worktree に移動
* `n` : 新しい worktree を作成
* `d` : worktree + branch を削除
* `Esc` / `Ctrl+C` : 終了

---

## 仕組み

* TUI は **stderr** に描画
* 選択結果（path）は **stdout** に出力
* shell function 側で `cd` を行う

```sh
gwt() {
  local out
  out="$(gwt-bin)" || return
  [ -n "$out" ] && cd "$out"
}
```

---

## 表示ルール

* `>` : カーソル
* `*` : current worktree
* `!` : dirty
* `@` : detached HEAD
* path は repo root からの相対パス

---

## 設計方針

* 自分のみを想定
* 再描画しない（操作後は終了）
* 安全第一（壊れないことを優先）
* 余計な機能は足さない

---

## 注意

* git repository 内で実行すること
* bash / zsh を想定
