update_submodule:
	git submodule init
	git submodule update --remote --merge
	git commit -am 'Updated submodules'
	git push -u origin master

.PHONY: update_submodule