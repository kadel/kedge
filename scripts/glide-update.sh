#!/bin/sh 

# for more information why we are doing this see comments glide.yaml


OPENSHIFT_VENDOR="vendor/github.com/openshift/origin/vendor"
PROJECT_VENDOR="./vendor"
TMPDIR=`mktemp -d`


# get all dependencies
glide update

# copy OpenShift vendor to temporary directory
cp -r $OPENSHIFT_VENDOR/* $TMPDIR


#strip vendor (we have to do this so we don't end up with nasted vendoring)
# glide is going to use cache, so it shouldn't require downloading pkgs again
glide update -v


# how deep is the package in directory strucutre?
# example: package hosted in github.com has three level strucure (github.com/<namespace>/<pkgname>)
#          packages in k8s.io has only two level strucutre (k8s.io/<pkgname>)
# we need to know where it is so we can move whole package (nothing more nothing less)
#
# If we were to move whole github.com than whole github.com gets replaced in target directory,
# and if target directory had some extra library from gihtub.com that is not in openshift vendor it would get removed.

TWO_LEVEL="cloud.google.com go.pedge.io go4.org google.golang.org gopkg.in k8s.io vbom.ml"
THREE_LEVEL="bitbucket.org github.com golang.org"


# check if we cover everthing from openshift vendor
# every domain in OpenShift vendor has to be covered in lists above (defined if its 2 or 3 level structure)
for path in `find $TMPDIR -maxdepth 1 -mindepth 1 -type d | sort`; do
    domain=`basename $path`

    found=false
    for t in $TWO_LEVEL $THREE_LEVEL; do
        if [ "$t" == "$domain" ]; then
            found=true
        fi
    done

    if [ $found == false ]; then
        echo "ERROR: structure for $domain is not defined"
        exit 1
    fi

done



# move packages from tmp dir to project vendor directory

for domain in $TWO_LEVEL; do
    pkgs=`find -L "${TMPDIR}/$domain"  -maxdepth 1 -mindepth 1 -type d -printf "$domain/%P\n"`
    for pkg in $pkgs; do
        mkdir -p $PROJECT_VENDOR/$pkg
        target=`dirname $PROJECT_VENDOR/$pkg`
        mv -f $TMPDIR/$pkg $target
    done
done

for domain in $THREE_LEVEL; do
    pkgs=`find -L "${TMPDIR}/$domain"  -maxdepth 2 -mindepth 2 -type d -printf "$domain/%P\n"`
    for pkg in $pkgs; do
        mkdir -p $PROJECT_VENDOR/$pkg
        target=`dirname $PROJECT_VENDOR/$pkg`
        mv -f $TMPDIR/$pkg $target
    done
done




rm -r $TMPDIR