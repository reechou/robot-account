# maybe more powerful
# for mac (sed for linux is different)
dir=`echo ${PWD##*/}`
grep "real-film" * -R | grep -v Godeps | awk -F: '{print $1}' | sort | uniq | xargs sed -i '' "s#real-film#$dir#g"
#mv http_demo.ini $dir.ini

