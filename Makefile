all: install move

install:
	go install

bashrc:
	$(shell scripts/add_line_bashrc.sh)

move:
	$(shell scripts/move_db.sh)


