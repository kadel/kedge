#!/bin/bash

DATE=`date --iso-8601=date`
TIME=`date --iso-8601=seconds`


BINTRAY_SUBJECT="kedgeproject"
BINTRAY_REPO="kedge"
BINTRAY_PACKAGE="kedge"
BINTRAY_API="https://bintray.com/api/v1"
BINTRAY_AUTH="tkral:5ecf6c9b2c4e850e582f6c6845c5e9b0e677a104"

if [ "${TRAVIS_BRANCH}" == "master" ] && [ "${TRAVIS_PULL_REQUEST}" == "false" ]; then
# this is not a pull request, its build on master branch
    echo "master build"
fi


if [ "${TRAVIS_PULL_REQUEST}" == "false" ]; then
# this is not a pull request, its build on master branch
    echo "master build"
fi


upload_file() {
    PUT /content/:subject/:repo/:package/:version/:file_path[?publish=0/1][?override=0/1][?explode=0/1]

    curl -T <FILE.EXT> -utkral:<API_KEY> https://api.bintray.com/content/kedgeproject/kedge/kedge/latest/<FILE_TARGET_PATH>

}


create_version() {
    version="1.1.5"
    desc="description"
    vcs_tag="1.1.5"

    read -r -d '' data << EOM
        {
            "name": "${version}", 
            "released": "${TIME}",
            "desc": "${desc}",
            "vcs_tag": "${vcs_tag}"
        }
EOM

    curl  -v -X POST -H "Content-Type: application/json"  -u ${BINTRAY_AUTH} -d "${data}" ${BINTRAY_API}/packages/${BINTRAY_SUBJECT}/${BINTRAY_REPO}/${BINTRAY_PACKAGE}/versions
}
