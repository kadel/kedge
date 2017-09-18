#!/usr/bin/groovy
@Library('github.com/fabric8io/fabric8-pipeline-library@master')
def dummy
goNode{
  dockerNode{
    ws{
      if (env.BRANCH_NAME.startsWith('PR-')) {
        dir(buildPath) {
          checkout scm
          
          container(name: 'go') {
            stage ('run test') {
              version = "SNAPSHOT-${env.BRANCH_NAME}-${env.BUILD_NUMBER}"
              sh "make clean test"    
            }
          }
        }
      }
    } 
  }
}


