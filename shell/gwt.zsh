gwt() {
  local out
  local exit_code

  out="$(gwt-bin)"
  exit_code=$?

  if [ $exit_code -ne 0 ] || [ -z "$out" ]; then
    return
  fi

  cd "$out" || return
}
