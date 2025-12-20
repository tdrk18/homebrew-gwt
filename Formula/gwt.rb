class Gwt < Formula
  desc "Interactive TUI helper for git worktree management"
  homepage "https://github.com/tdrk18/homebrew-gwt"
  url "https://github.com/tdrk18/homebrew-gwt/releases/download/v0.1.0/gwt-v0.1.0.tar.gz"
  sha256 "b46733f295b80ce2ca5ec35f093fa0f371bbec54f8945bcf28fe07b846d661a4"

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
