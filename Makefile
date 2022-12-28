# Makefile CHEATSHEET: https://devhints.io/makefile
##############################################################################
include Makefile.settings
##############################################################################
# Meta

menu :
	$(INFO) 'Manage source code :'
	@echo '	push  : git push -u origin master'
	@echo '	tag   : git tag v${VER_APP}  (VER_APP)'
	@echo '	untag : git … : remove v${VER_APP}  (VER_APP)'

##############################################################################
# Source 

# git remote add origin git@github.com:$_USERNAME/$_REPONAME.git  # ssh mode
push :
	gc
	git push -u origin master
tag :
ifeq (v${VER_APP}, $(shell git tag |grep v${VER_APP}))
	@echo 'repo ALREADY tagged @ "v${VER_APP}" : VER_APP'
else 
	git tag v${VER_APP}
	git push origin v${VER_APP}
	git tag
endif
untag :
	git tag -d v${VER_APP}
	git push origin --delete v${VER_APP}
markup :
	bash make.md2html.sh
tarball :
	bash make.tarball.sh
perms :
	bash make.perms.sh
