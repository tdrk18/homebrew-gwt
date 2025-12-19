class Gwt < Formula
  desc "Interactive TUI helper for git worktree management"
  homepage "https://github.com/tdrk18/homebrew-gwt"
  url "https://github.com/tdrk18/homebrew-gwt/releases/download/0.1.0/gwt-0.1.0.tar.gz"
  sha256 "5c79eb63e438584e4fbf4672451c48b51ed49d2d1d78af10c56c3e705e25810d"

  def install
    bin.install "gwt-bin"
    pkgshare.install "shell/gwt.zsh"
    pkgshare.install "shell/gwt.bash"
  end

  def caveats
    <<~EOS
      gwt requires shell integration to work correctly.

      Add one of the following lines to your shell config:

        # zsh
        source #{opt_pkgshare}/gwt.zsh

        # bash
        source #{opt_pkgshare}/gwt.bash

      Restart your shell after updating the config.
    EOS
  end
end
